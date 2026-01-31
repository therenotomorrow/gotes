package server

import (
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/therenotomorrow/gotes/internal/config"
	"github.com/therenotomorrow/gotes/internal/services/secure"
)

func CORSMiddleware(cfg config.CORS) func(next http.Handler) http.Handler {
	allowedMethods := strings.Join(cfg.AllowedMethods, ", ")
	allowedOrigins := cfg.AllowedOrigins
	allowedHeaders := cfg.AllowedHeaders

	return func(next http.Handler) http.Handler {
		handlerFunc := func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Access-Control-Allow-Methods", allowedMethods)
			writer.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
			writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders)

			if request.Method == http.MethodOptions {
				writer.WriteHeader(http.StatusNoContent)

				return
			}

			next.ServeHTTP(writer, request)
		}

		return http.HandlerFunc(handlerFunc)
	}
}

func LoggingMiddleware(logger log) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handlerFunc := func(writer http.ResponseWriter, request *http.Request) {
			args := []any{"method", request.Method, "path", request.URL.Path}
			level := slog.LevelInfo

			logger.Log(request.Context(), level, "request", args...)

			next.ServeHTTP(writer, request)

			logger.Log(request.Context(), level, "response", args...)
		}

		return http.HandlerFunc(handlerFunc)
	}
}

func TrimSlashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			request.URL.Path = strings.TrimSuffix(request.URL.Path, "/")
		}

		next.ServeHTTP(writer, request)
	})
}

func WebSocketMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.HasSuffix(request.URL.Path, "/events") {
			request.Body = http.NoBody // fix server-side streaming deadlock for an empty message

			if token := request.URL.Query().Get("token"); token != "" {
				request.Header.Set(secure.AuthKey, token) // add query-based authorization
			}
		}

		next.ServeHTTP(writer, request)
	})
}

func ApplyMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	slices.Reverse(middlewares)

	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}
