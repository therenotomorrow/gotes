package postgres

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/api/notes/v1/ports"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
)

type StoreProvider struct {
	db postgres.Database
}

func NewStoreProvider(db postgres.Database) *StoreProvider {
	return &StoreProvider{db: db}
}

func (p *StoreProvider) Provide(ctx context.Context) ports.Store {
	conn := p.db.Conn(ctx)

	return ports.Store{Notes: NewNotesRepository(conn)}
}
