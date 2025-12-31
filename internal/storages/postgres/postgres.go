package postgres

import (
	"context"

	commandsnotes "github.com/therenotomorrow/gotes/internal/storages/postgres/commands/notes"
	commandsusers "github.com/therenotomorrow/gotes/internal/storages/postgres/commands/users"
	queriesnotes "github.com/therenotomorrow/gotes/internal/storages/postgres/queries/notes"
	queriesusers "github.com/therenotomorrow/gotes/internal/storages/postgres/queries/users"
)

type DBTX interface {
	commandsnotes.DBTX
	commandsusers.DBTX
	queriesusers.DBTX
	queriesnotes.DBTX
}

type Database interface {
	Conn(ctx context.Context) DBTX
	Tx(ctx context.Context, call func(ctx context.Context) error) error
	Close()
}
