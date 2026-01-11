package entities

import (
	"time"

	"github.com/therenotomorrow/gotes/internal/domain"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
)

const (
	ErrEmptyTitle   domain.Error = "empty title"
	ErrEmptyContent domain.Error = "empty content"
)

type Note struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Owner     *User
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
		Owner:     nil,
	}, nil
}

func (n *Note) IsOwner(u *User) bool {
	return n.Owner.ID == u.ID
}

func (n *Note) SetOwner(u *User) {
	n.Owner = u
}
