package server

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/therenotomorrow/ex"
	openapinotesv1 "github.com/therenotomorrow/gotes/docs/api/notes/v1"
	openapiusersv1 "github.com/therenotomorrow/gotes/docs/api/users/v1"
	"github.com/therenotomorrow/gotes/internal/config"
	pbnotesv1 "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	pbusersv1 "github.com/therenotomorrow/gotes/pkg/api/users/v1"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
	"github.com/therenotomorrow/gotes/tools/swagger"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGateway(cfg *config.Config, logger *slog.Logger) *http.Server {
	ctx := context.Background()
	tracer := trace.Service("gateway", logger)

	creds := insecure.NewCredentials()
	if cfg.Server.Secure {
		creds = credentials.NewTLS(new(tls.Config))
	}

	options := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	handler := http.NewServeMux()
	handler.Handle("GET /docs/", http.StripPrefix("/docs", http.FileServer(http.FS(swagger.Content))))

	// ---- NotesService
	notesGateway := runtime.NewServeMux()
	notesMiddlewares := []func(next http.Handler) http.Handler{
		tracer.Middleware,
		LoggingMiddleware(tracer),
		CORSMiddleware(cfg.Server.Gateway.CORS),
		TrimSlashMiddleware,
		WebSocketMiddleware,
	}
	err := pbnotesv1.RegisterNotesServiceHandlerFromEndpoint(ctx, notesGateway, cfg.Server.Address, options)
	ex.Panic(err)
	handler.Handle("/api/v1/", ApplyMiddlewares(notesGateway, notesMiddlewares...))
	handler.Handle(
		"GET /docs/notes/",
		http.StripPrefix("/docs/notes", http.FileServer(http.FS(openapinotesv1.Content))),
	)

	// ---- UsersService
	usersGateway := runtime.NewServeMux()
	usersMiddlewares := []func(next http.Handler) http.Handler{
		tracer.Middleware,
		LoggingMiddleware(tracer),
		CORSMiddleware(cfg.Server.Gateway.CORS),
	}
	err = pbusersv1.RegisterUsersServiceHandlerFromEndpoint(ctx, usersGateway, cfg.Server.Address, options)
	ex.Panic(err)
	handler.Handle("/api/v1/users/", ApplyMiddlewares(usersGateway, usersMiddlewares...))
	handler.Handle(
		"GET /docs/users/",
		http.StripPrefix("/docs/users", http.FileServer(http.FS(openapiusersv1.Content))),
	)

	gateway := new(http.Server)
	gateway.Addr = cfg.Server.Gateway.Address
	gateway.Handler = wsproxy.WebsocketProxy(handler)

	return gateway
}
