package entities

import (
	"time"

	"github.com/therenotomorrow/gotes/internal/domain/types/uuid"
)

type EventType int

const (
	EventTypeCreated EventType = iota
	EventTypeDeleted
)

type Event struct {
	EventTime time.Time
	Note      *Note
	ID        uuid.UUID
	EventType EventType
}

func NewEvent(t EventType, n *Note) *Event {
	return &Event{
		ID:        uuid.New(),
		EventType: t,
		Note:      n,
		EventTime: time.Now(),
	}
}
