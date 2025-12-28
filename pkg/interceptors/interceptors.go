package interceptors

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/therenotomorrow/gotes/pkg/tracer"
	"google.golang.org/grpc"
)

func Logging(logger *slog.Logger) grpc.UnaryServerInterceptor {
	trace := tracer.New(logger)
	logFn := func(ctx context.Context, level logging.Level, message string, fields ...any) {
		fields = append(fields, "trace", trace.Trace(ctx))

		logger.Log(ctx, slog.Level(level), message, fields...)
	}

	return logging.UnaryServerInterceptor(
		logging.LoggerFunc(logFn),
		logging.WithDisableLoggingFields(
			"grpc.method_type", "protocol", "grpc", "grpc.component",
		),
	)
}

func Trace(logger *slog.Logger) grpc.UnaryServerInterceptor {
	trace := tracer.New(logger)

	return func(ctx context.Context, request any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = trace.New(ctx)

		return handler(ctx, request)
	}
}
