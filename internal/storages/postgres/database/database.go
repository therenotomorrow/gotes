package database

import (
	"context"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
	"github.com/therenotomorrow/gotes/pkg/services/validate"
)

const (
	ErrInvalidConfig ex.Error = "invalid config"
)

type Config struct {
	DSN string `json:"dsn" validate:"required,postgres_dsn"`
}

type Postgres struct {
	trm    *manager.Manager
	ctx    *pgxv5.CtxGetter
	pool   *pgxpool.Pool
	config Config
}

func New(cfg Config, logger *slog.Logger) (*Postgres, error) {
	err := validate.Struct(cfg)
	if err != nil {
		return nil, ErrInvalidConfig.Because(err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		return nil, ErrInvalidConfig.Because(err)
	}

	tracer := trace.Service("postgres", logger)
	manOpts := []manager.Opt{
		manager.WithLog(log(func(ctx context.Context, msg string) { tracer.Warning(ctx, msg) })),
	}

	trm, err := manager.New(pgxv5.NewDefaultFactory(pool), manOpts...)
	if err != nil {
		return nil, ErrInvalidConfig.Because(err)
	}

	ctx := pgxv5.DefaultCtxGetter

	return &Postgres{trm: trm, ctx: ctx, config: cfg, pool: pool}, nil
}

func MustNew(cfg Config, logger *slog.Logger) *Postgres {
	db, err := New(cfg, logger)

	return ex.Critical(db, err)
}

func (db *Postgres) Tx(ctx context.Context, call func(ctx context.Context) error) error {
	return db.trm.Do(ctx, call)
}

func (db *Postgres) Conn(ctx context.Context) postgres.DBTX {
	return db.ctx.DefaultTrOrDB(ctx, db.pool)
}

func (db *Postgres) Default() postgres.DBTX {
	return db.pool
}

func (db *Postgres) Close() {
	db.pool.Close()
}

type log func(ctx context.Context, msg string)

func (l log) Warning(ctx context.Context, msg string) {
	l(ctx, msg)
}
