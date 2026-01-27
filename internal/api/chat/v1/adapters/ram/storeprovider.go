package ram

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/api/chat/v1/ports"
)

type StoreProvider struct {
	repo ports.MessagesRepository
}

func NewStoreProvider() *StoreProvider {
	return &StoreProvider{repo: NewMessagesRepository()}
}

func (p *StoreProvider) Provide(_ context.Context) ports.Store {
	return ports.Store{Messages: p.repo}
}
