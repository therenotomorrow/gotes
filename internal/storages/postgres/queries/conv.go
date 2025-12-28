package queries

import (
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
)

func (n *Note) Entity() *entities.Note {
	return &entities.Note{
		ID:        id.New(n.ID),
		Title:     n.Title,
		Content:   n.Content,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

type Notes []*Note

func (n Notes) Entities() []*entities.Note {
	notes := make([]*entities.Note, len(n))
	for i, note := range n {
		notes[i] = note.Entity()
	}

	return notes
}
