package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/storages/postgres/commands"
	"github.com/therenotomorrow/gotes/internal/storages/postgres/queries"
	"github.com/therenotomorrow/gotes/pkg/tracer"
	"github.com/therenotomorrow/gotes/pkg/validate"
)

type txCtx string

const (
	txKey txCtx = "tx"

	ErrTxError       ex.Error = "transaction error"
	ErrInvalidConfig ex.Error = "invalid config"
)

type Config struct {
	Logger *slog.Logger
	DSN    string `json:"dsn" validate:"required,postgres_dsn"`
}

type Database struct {
	pool   *pgxpool.Pool
	tracer *tracer.Tracer
	config Config
}

func New(cfg Config) (*Database, error) {
	err := validate.Struct(cfg)
	if err != nil {
		return nil, ErrInvalidConfig.Because(err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		return nil, ErrInvalidConfig.Because(err)
	}

	return &Database{
		config: cfg,
		pool:   pool,
		tracer: tracer.Service("postgres", cfg.Logger),
	}, nil
}

func MustNew(cfg Config) *Database {
	db, err := New(cfg)

	return ex.Critical(db, err)
}

func (db *Database) Tx(ctx context.Context, call func(ctx context.Context) error) (err error) {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if ok {
		return call(ctx)
	}

	tx, err = db.pool.Begin(ctx)
	if err != nil {
		return ErrTxError.Because(err)
	}

	defer func() {
		var txErr error

		if err != nil {
			txErr = tx.Rollback(ctx)
		} else {
			txErr = tx.Commit(ctx)
		}

		if txErr != nil {
			db.tracer.Error(ctx, "transaction failed", txErr)

			err = errors.Join(err, ErrTxError.Because(txErr))
		}
	}()

	return call(context.WithValue(ctx, txKey, tx))
}

func (db *Database) Close() {
	db.pool.Close()
}

type CQRS struct {
	Commands commands.Querier
	Queries  queries.Querier
}

func (db *Database) CQRS(ctx context.Context) *CQRS {
	cmds := commands.New(db.pool)
	query := queries.New(db.pool)

	if tx, ok := ctx.Value(txKey).(pgx.Tx); ok {
		cmds = cmds.WithTx(tx)
		query = query.WithTx(tx)
	}

	return &CQRS{Commands: cmds, Queries: query}
}
