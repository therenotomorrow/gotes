package postgres

import (
	"context"

	"github.com/redis/go-redis/v9"
	adapters "github.com/therenotomorrow/gotes/internal/api/notes/v1/adapters/redis"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/ports"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
)

type StoreProvider struct {
	db  postgres.Database
	rdb redis.UniversalClient
}

func NewStoreProvider(db postgres.Database, rdb redis.UniversalClient) *StoreProvider {
	return &StoreProvider{db: db, rdb: rdb}
}

func (p *StoreProvider) Provide(ctx context.Context) ports.Store {
	conn := p.db.Conn(ctx)

	return ports.Store{
		Notes:  NewNotesRepository(conn),
		Events: adapters.NewEventsRepository(p.rdb),
	}
}
