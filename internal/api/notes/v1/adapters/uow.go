package adapters

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/api/notes/v1/ports"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
)

type UnitOfWork struct {
	database *postgres.Database
}

func NewUnitOfWork(database *postgres.Database) *UnitOfWork {
	return &UnitOfWork{database: database}
}

func (u *UnitOfWork) Do(ctx context.Context, work func(store ports.Store) error) error {
	return u.database.Tx(ctx, func(ctx context.Context) error {
		cqrs := u.database.CQRS(ctx)
		store := NewStore(cqrs)

		return work(store)
	})
}
