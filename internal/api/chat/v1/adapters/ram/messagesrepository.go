package ram

import (
	"context"
	"maps"
	"slices"
	"sync"

	"github.com/therenotomorrow/gotes/internal/api/chat/v1/entities"
)

type MessagesRepository struct {
	data  map[string]*entities.Message
	mutex sync.RWMutex
}

func NewMessagesRepository() *MessagesRepository {
	return &MessagesRepository{
		mutex: sync.RWMutex{},
		data:  make(map[string]*entities.Message),
	}
}

func (m *MessagesRepository) SaveMessage(_ context.Context, message *entities.Message) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data[message.Header.CorrelationID] = message

	return nil
}

func (m *MessagesRepository) Outbox(_ context.Context) ([]*entities.Message, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return slices.Collect(maps.Values(m.data)), nil
}

func (m *MessagesRepository) DeleteMessage(_ context.Context, message *entities.Message) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.data, message.Header.CorrelationID)

	return nil
}
