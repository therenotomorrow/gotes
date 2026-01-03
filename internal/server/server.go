package server

import (
	"context"
	"log/slog"
	"net"
	"os/signal"
	"sync"
	"syscall"

	"github.com/therenotomorrow/ex"
	v1 "github.com/therenotomorrow/gotes/internal/api/notes/v1"
	"github.com/therenotomorrow/gotes/internal/config"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
	pb "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	"github.com/therenotomorrow/gotes/pkg/interceptors"
	"github.com/therenotomorrow/gotes/pkg/tracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Dependencies struct {
	Database *postgres.Database
	Logger   *slog.Logger
}

func (d *Dependencies) Close() {
	d.Database.Close()
}

type Server struct {
	deps   Dependencies
	tracer *tracer.Tracer
	grpc   *grpc.Server
	config config.Config
	once   sync.Once
}

func New(cfg config.Config, deps Dependencies) *Server {
	logger := deps.Logger
	server := setupServer(deps, logger)

	return &Server{
		config: cfg,
		deps:   deps,
		tracer: tracer.New(logger),
		grpc:   server,
		once:   sync.Once{},
	}
}

func setupServer(deps Dependencies, logger *slog.Logger) *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.Trace(logger),
			interceptors.Logging(logger),
		),
	)

	pb.RegisterNotesServiceServer(srv, v1.Service(deps.Database))
	reflection.Register(srv)

	return srv
}

func (s *Server) Serve() {
	ctx := context.Background()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var lis net.Listener

	s.once.Do(func() {
		var lc net.ListenConfig

		lis = ex.Critical(lc.Listen(ctx, "tcp", s.config.Server.Address))
	})

	go func() {
		s.tracer.Info(ctx, "listen...", "address", s.config.Server.Address)

		err := s.grpc.Serve(lis)
		if err != nil {
			s.tracer.Error(ctx, "listen failure", ex.Unexpected(err))

			stop()
		}
	}()

	<-ctx.Done()

	s.tracer.Info(ctx, "shutdown...")

	s.Stop()
}

func (s *Server) Stop() {
	s.grpc.Stop()
	s.deps.Close()
}
