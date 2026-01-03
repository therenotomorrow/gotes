package tracer

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/therenotomorrow/ex"
)

type tracerCtx string

const tracerKey tracerCtx = "tracer"

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
	tracer.logger = logger.With("service", service)

	return tracer
}

func (t *Tracer) New(ctx context.Context) context.Context {
	if _, ok := ctx.Value(tracerKey).(uuid.UUID); ok {
		return ctx
	}

	return context.WithValue(ctx, tracerKey, uuid.New())
}

func (t *Tracer) Log(ctx context.Context, where string, err error) {
	if err == nil {
		t.Info(ctx, where)

		return
	}

	err, cause := ex.Expose(err)
	t.Error(ctx, where, err, "cause", cause)
}

func (t *Tracer) Trace(ctx context.Context) uuid.UUID {
	if tid, ok := ctx.Value(tracerKey).(uuid.UUID); ok {
		return tid
	}

	return uuid.Nil
}

func (t *Tracer) Info(ctx context.Context, msg string, kvs ...any) {
	kvs = append(kvs, "trace", t.Trace(ctx))

	t.logger.InfoContext(ctx, msg, kvs...)
}

func (t *Tracer) Error(ctx context.Context, msg string, err error, kvs ...any) {
	kvs = append(kvs, "error", err, "trace", t.Trace(ctx))

	t.logger.ErrorContext(ctx, msg, kvs...)
}
