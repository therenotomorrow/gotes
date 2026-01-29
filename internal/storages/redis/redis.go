package redis

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
	"github.com/therenotomorrow/gotes/pkg/services/validate"
)

const (
	ErrInvalidConfig ex.Error = "invalid config"
)

type Config struct {
	Address  string `json:"address"  validate:"required"`
	Password string `json:"password" validate:"required"`
}

func New(cfg Config, logger *slog.Logger) (*redis.Client, error) {
	err := validate.Struct(cfg)
	if err != nil {
		return nil, ErrInvalidConfig.Because(err)
	}

	tracer := trace.Service("redis", logger)

	redis.SetLogger(log(func(ctx context.Context, format string, args ...any) {
		tracer.Info(ctx, format, args...)
	}))

	options := new(redis.Options)

	options.Addr = cfg.Address
	options.Password = cfg.Password

	client := redis.NewClient(options)

	return client, nil
}

func MustNew(cfg Config, logger *slog.Logger) *redis.Client {
	cli, err := New(cfg, logger)

	ex.Panic(err)

	return cli
}

type log func(ctx context.Context, format string, args ...any)

func (l log) Printf(ctx context.Context, format string, args ...any) {
	l(ctx, format, args...)
}
