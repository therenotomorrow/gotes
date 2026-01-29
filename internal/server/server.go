package server

import (
	"context"
	"log/slog"
	"net"
	"os/signal"
	"sync"
	"syscall"

	"buf.build/go/protovalidate"
	"github.com/therenotomorrow/ex"
	chatv1 "github.com/therenotomorrow/gotes/internal/api/chat/v1"
	metricsv1 "github.com/therenotomorrow/gotes/internal/api/metrics/v1"
	notesv1 "github.com/therenotomorrow/gotes/internal/api/notes/v1"
	usersv1 "github.com/therenotomorrow/gotes/internal/api/users/v1"
	"github.com/therenotomorrow/gotes/internal/config"
	"github.com/therenotomorrow/gotes/internal/domain/types/email"
	"github.com/therenotomorrow/gotes/internal/domain/types/password"
	"github.com/therenotomorrow/gotes/internal/domain/types/uuid"
	"github.com/therenotomorrow/gotes/internal/services/secure"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
	"github.com/therenotomorrow/gotes/internal/storages/redis"
	pbchatv1 "github.com/therenotomorrow/gotes/pkg/api/chat/v1"
	pbmetricsv1 "github.com/therenotomorrow/gotes/pkg/api/metrics/v1"
	pbnotesv1 "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	pbusersv1 "github.com/therenotomorrow/gotes/pkg/api/users/v1"
	"github.com/therenotomorrow/gotes/pkg/services/generate"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
	"github.com/therenotomorrow/gotes/pkg/services/validate"
	"github.com/therenotomorrow/gotes/pkg/services/vault"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	deps   *Dependencies
	logger *slog.Logger
	grpc   *grpc.Server
	config *config.Config
	once   sync.Once
}

func New(cfg *config.Config, deps *Dependencies, logger *slog.Logger) *Server {
	tracer := trace.New(logger)
	validator := ex.Critical(protovalidate.New())

	server := grpc.NewServer(
		grpc.MaxConcurrentStreams(cfg.Server.MaxConcurrentStreams),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionIdle:     cfg.Server.KeepAlive.MaxConnection.Idle,
				MaxConnectionAge:      cfg.Server.KeepAlive.MaxConnection.Age,
				MaxConnectionAgeGrace: cfg.Server.KeepAlive.MaxConnection.AgeGrace,
				Time:                  cfg.Server.KeepAlive.Time,
				Timeout:               cfg.Server.KeepAlive.Timeout,
			},
		),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             cfg.Server.KeepAlive.EnforcementPolicy.MinTime,
			PermitWithoutStream: cfg.Server.KeepAlive.EnforcementPolicy.PermitWithoutStream,
		}),
		grpc.ChainUnaryInterceptor(
			tracer.UnaryServerInterceptor,
			LoggingUnaryServerInterceptor(tracer),
			validate.UnaryServerInterceptor(validator),
			secure.UnaryServerInterceptor(deps.Authenticator, []string{
				"/api.users.v1.UsersService/RegisterUser",
				"/api.users.v1.UsersService/RefreshToken",
			}...),
		),
		grpc.ChainStreamInterceptor(
			tracer.StreamServerInterceptor,
			LoggingStreamServerInterceptor(tracer),
			secure.StreamServerInterceptor(deps.Authenticator, []string{
				"/api.chat.v1.ChatService/Dispatch",
				"/api.metrics.v1.MetricsService/UploadMetrics",
			}...),
		),
	)

	email.SetValidator(deps.EmailValidator)
	uuid.SetGenerator(deps.UUIDGenerator)
	password.SetHasher(deps.PasswordHasher)

	pbmetricsv1.RegisterMetricsServiceServer(server, metricsv1.NewService(logger))
	pbnotesv1.RegisterNotesServiceServer(server, notesv1.NewService(deps.Database, deps.Redis, logger))
	pbusersv1.RegisterUsersServiceServer(server, usersv1.NewService(deps.Database, logger))
	pbchatv1.RegisterChatServiceServer(server, chatv1.NewService(validator, logger))

	if cfg.Debug {
		reflection.Register(server)
	}

	return &Server{logger: logger, grpc: server, config: cfg, deps: deps, once: sync.Once{}}
}

func Default(cfg *config.Config) *Server {
	logger := trace.Logger(trace.JSON, cfg.Debug)
	database := postgres.MustNew(postgres.Config{DSN: cfg.Postgres.DSN}, logger)
	rdb := redis.MustNew(redis.Config{Address: cfg.Redis.Address, Password: cfg.Redis.Password}, logger)

	return New(cfg, &Dependencies{
		Database:       database,
		Redis:          rdb,
		Authenticator:  secure.NewTokenAuthenticator(database),
		PasswordHasher: vault.NewPasswordHasher(),
		UUIDGenerator:  generate.NewUUID(),
		EmailValidator: validate.NewEmail(),
	}, logger)
}

func (s *Server) Serve(ctx context.Context) {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var lis net.Listener

	s.once.Do(func() {
		var lc net.ListenConfig

		lis = ex.Critical(lc.Listen(ctx, "tcp", s.config.Server.Address))
	})

	s.logger.InfoContext(ctx, "listen...", "address", s.config.Server.Address)

	defer s.Stop(ctx)

	go func() {
		err := s.grpc.Serve(lis)

		ex.Panic(err)
	}()

	<-ctx.Done()
}

func (s *Server) Stop(ctx context.Context) {
	s.logger.InfoContext(ctx, "shutdown...")

	s.grpc.Stop()
	s.deps.Close()
}
