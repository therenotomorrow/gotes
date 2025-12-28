package commands

import (
	"github.com/therenotomorrow/gotes/internal/domain/entities"
)

func NewInsertNoteParams(n *entities.Note) *InsertNoteParams {
	return &InsertNoteParams{
		Title:     n.Title,
		Content:   n.Content,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}
