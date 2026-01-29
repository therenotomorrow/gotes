package ports

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/api/chat/v1/entities"
)

type MessagesRepository interface {
	SaveMessage(ctx context.Context, message *entities.Message) error
	Outbox(ctx context.Context) ([]*entities.Message, error)
	DeleteMessage(ctx context.Context, message *entities.Message) error
}

type Store struct {
	Messages MessagesRepository
}

type StoreProvider interface {
	Provide(ctx context.Context) Store
}
