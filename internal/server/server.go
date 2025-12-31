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
	notesv1 "github.com/therenotomorrow/gotes/internal/api/notes/v1"
	usersv1 "github.com/therenotomorrow/gotes/internal/api/users/v1"
	"github.com/therenotomorrow/gotes/internal/config"
	"github.com/therenotomorrow/gotes/internal/domain/types/email"
	"github.com/therenotomorrow/gotes/internal/domain/types/password"
	"github.com/therenotomorrow/gotes/internal/domain/types/uuid"
	"github.com/therenotomorrow/gotes/internal/services/auth"
	"github.com/therenotomorrow/gotes/internal/storages/postgres/database"
	pbnotesv1 "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	pbusersv1 "github.com/therenotomorrow/gotes/pkg/api/users/v1"
	"github.com/therenotomorrow/gotes/pkg/services/generate"
	"github.com/therenotomorrow/gotes/pkg/services/secure"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
	"github.com/therenotomorrow/gotes/pkg/services/validate"
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
			auth.UnaryServerInterceptor(deps.Authenticator, []string{
				"/api.users.v1.UsersService/RegisterUser",
				"/api.users.v1.UsersService/RefreshToken",
			}...),
		),
	)

	email.SetValidator(deps.EmailValidator)
	uuid.SetGenerator(deps.UUIDGenerator)
	password.SetHasher(deps.PasswordHasher)

	pbnotesv1.RegisterNotesServiceServer(server, notesv1.Service(deps.Secure, deps.Database, logger))
	pbusersv1.RegisterUsersServiceServer(server, usersv1.Service(deps.Database, logger))

	if cfg.Debug {
		reflection.Register(server)
	}

	return &Server{logger: logger, grpc: server, config: cfg, deps: deps, once: sync.Once{}}
}

func Default(cfg *config.Config) *Server {
	var (
		logger        = trace.Logger(trace.JSON, cfg.Debug)
		postgres      = database.MustNew(database.Config{DSN: cfg.Postgres.DSN}, logger)
		securable     = auth.Secure{}
		authenticator = auth.NewTokenAuthenticator(postgres)
		hasher        = secure.NewPasswordHasher()
		generator     = generate.NewUUID()
		validator     = validate.NewEmail()
	)

	return New(cfg, &Dependencies{
		Database:       postgres,
		Secure:         securable,
		Authenticator:  authenticator,
		PasswordHasher: hasher,
		UUIDGenerator:  generator,
		EmailValidator: validator,
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
