package trace

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type tracerCtx string

const tracerKey tracerCtx = "trace"

type Tracer struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Tracer {
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	return &Tracer{logger: logger}
}

func Service(service string, logger *slog.Logger) *Tracer {
	tracer := New(logger)
	tracer.logger = tracer.logger.With("service", service)

	return tracer
}

func (t *Tracer) Context(ctx context.Context) context.Context {
	if _, ok := ctx.Value(tracerKey).(uuid.UUID); ok {
		return ctx
	}

	return context.WithValue(ctx, tracerKey, uuid.New())
}

func (t *Tracer) Log(ctx context.Context, level slog.Level, msg string, kvs ...any) {
	kvs = t.trace(ctx, kvs...)

	t.logger.Log(ctx, level, msg, kvs...)
}

func (t *Tracer) Info(ctx context.Context, msg string, kvs ...any) {
	kvs = t.trace(ctx, kvs...)

	t.logger.InfoContext(ctx, msg, kvs...)
}

func (t *Tracer) Warning(ctx context.Context, msg string, kvs ...any) {
	kvs = t.trace(ctx, kvs...)

	t.logger.WarnContext(ctx, msg, kvs...)
}

func (t *Tracer) Error(ctx context.Context, msg string, err error, kvs ...any) {
	kvs = t.trace(ctx, append(kvs, "error", err)...)

	t.logger.ErrorContext(ctx, msg, kvs...)
}

func (t *Tracer) trace(ctx context.Context, kvs ...any) []any {
	trace := uuid.Nil
	if tid, ok := ctx.Value(tracerKey).(uuid.UUID); ok {
		trace = tid
	}

	return append(kvs, "trace", trace)
}
