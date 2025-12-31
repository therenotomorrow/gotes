package queries

import (
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
)

func (n *Note) ToEntity() *entities.Note {
	return &entities.Note{
		ID:        id.New(n.ID),
		Title:     n.Title,
		Content:   n.Content,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Owner:     setOwner(n.UserID),
	}
}

func setOwner(userID *int64) *entities.User {
	if userID == nil {
		return nil
	}

	owner := new(entities.User)
	owner.ID = id.New(*userID)

	return owner
}

type Notes []*Note

func (n Notes) ToEntities() []*entities.Note {
	notes := make([]*entities.Note, len(n))
	for i, note := range n {
		notes[i] = note.ToEntity()
	}

	return notes
}
