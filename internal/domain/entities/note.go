package entities

import (
	"time"

	"github.com/therenotomorrow/gotes/internal/domain/types/id"
)

const (
	ErrEmptyTitle   Error = "empty title"
	ErrEmptyContent Error = "empty content"
)

type Note struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
	Content   string
	ID        id.ID
}

func NewNote(title, content string) (*Note, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}

	if content == "" {
		return nil, ErrEmptyContent
	}

	now := time.Now()

	return &Note{
		ID:        id.ID{},
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
