package entities

import (
	"github.com/therenotomorrow/gotes/internal/domain"
	"github.com/therenotomorrow/gotes/internal/domain/types/uuid"
)

const (
	ErrEmptyText domain.Error = "empty text"
)

type Message struct {
	Header Header
	Text   string
}

func NewMessage(text, correlationID string) (*Message, error) {
	if text == "" {
		return nil, ErrEmptyText
	}

	if correlationID == "" {
		correlationID = uuid.New().Value()
	}

	return &Message{
		Header: Header{CorrelationID: correlationID},
		Text:   text,
	}, nil
}
