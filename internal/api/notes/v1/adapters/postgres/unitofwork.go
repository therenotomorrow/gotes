package postgres

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/api/notes/v1/ports"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
)

type UnitOfWork struct {
	db       postgres.Database
	provider ports.StoreProvider
}

func NewUnitOfWork(db postgres.Database, provider ports.StoreProvider) *UnitOfWork {
	return &UnitOfWork{db: db, provider: provider}
}

func (uow *UnitOfWork) Do(ctx context.Context, work func(store ports.Store) error) error {
	return uow.db.Tx(ctx, func(ctx context.Context) error {
		store := uow.provider.Provide(ctx)

		return work(store)
	})
}
