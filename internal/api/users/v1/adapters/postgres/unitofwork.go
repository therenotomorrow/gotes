package postgres

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/api/users/v1/ports"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
)

type UnitOfWork struct {
	db       postgres.Database
	provider ports.StoreProvider
}

func NewUnitOfWork(db postgres.Database, provider ports.StoreProvider) *UnitOfWork {
	return &UnitOfWork{db: db, provider: provider}
}

func (u *UnitOfWork) Do(ctx context.Context, work func(store ports.Store) error) error {
	return u.db.Tx(ctx, func(ctx context.Context) error {
		store := u.provider.Provide(ctx)

		return work(store)
	})
}
