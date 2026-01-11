package server

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type log interface {
	Log(ctx context.Context, level slog.Level, msg string, kvs ...any)
}

func LoggingUnaryServerInterceptor(logger log) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		args := []any{"method", info.FullMethod}
		level := slog.LevelInfo
		start := time.Now()

		logger.Log(ctx, level, "request", args...)

		resp, err := handler(ctx, req)

		args = append(args, "duration", time.Since(start).String(), "status", status.Code(err).String())

		if err != nil {
			level = slog.LevelError

			args = append(args, "error", err)
		}

		logger.Log(ctx, level, "response", args...)

		return resp, err
	}
}
