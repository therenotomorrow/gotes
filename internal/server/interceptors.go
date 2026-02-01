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

type wrappedServerStream struct {
	grpc.ServerStream

	logger log
}

func (wss *wrappedServerStream) SendMsg(m any) error {
	wss.logger.Log(wss.Context(), slog.LevelInfo, "sending stream message", "message", m)

	return wss.ServerStream.SendMsg(m)
}

func (wss *wrappedServerStream) RecvMsg(m any) error {
	args := []any{"message", m}
	level := slog.LevelInfo

	err := wss.ServerStream.RecvMsg(m)
	if err != nil {
		level = slog.LevelError

		args = append(args, "error", err)
	}

	wss.logger.Log(wss.Context(), level, "receiving stream message", args...)

	return err
}

func LoggingStreamServerInterceptor(logger log) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		args := []any{"method", info.FullMethod}
		level := slog.LevelInfo

		wrapped := &wrappedServerStream{ServerStream: stream, logger: logger}

		logger.Log(wrapped.Context(), level, "start streaming", args...)

		err := handler(srv, wrapped)
		if err != nil {
			level = slog.LevelError

			args = append(args, "error", err)
		}

		logger.Log(wrapped.Context(), level, "end streaming", args...)

		return err
	}
}
