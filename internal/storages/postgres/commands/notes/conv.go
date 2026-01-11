package commands

import "github.com/therenotomorrow/gotes/internal/domain/entities"

func NewInsertNoteParams(note *entities.Note) *InsertNoteParams {
	return &InsertNoteParams{
		Title:     note.Title,
		Content:   note.Content,
		UserID:    note.Owner.ID.ValuePtr(),
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}
