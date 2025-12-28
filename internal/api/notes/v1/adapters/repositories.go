package adapters

import (
	"context"
	"database/sql"
	"errors"

	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/domain"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
	"github.com/therenotomorrow/gotes/internal/storages/postgres/commands"
	"github.com/therenotomorrow/gotes/internal/storages/postgres/queries"
)

type NotesRepository struct {
	cqrs *postgres.CQRS
}

func NewNotesRepository(cqrs *postgres.CQRS) *NotesRepository {
	return &NotesRepository{cqrs: cqrs}
}

func (r *NotesRepository) SaveNote(ctx context.Context, note *entities.Note) (*entities.Note, error) {
	ident, err := r.cqrs.Commands.InsertNote(ctx, commands.NewInsertNoteParams(note))
	if err != nil {
		return nil, ex.Unexpected(err)
	}

	note.ID = id.New(ident)

	return note, nil
}

func (r *NotesRepository) GetNote(ctx context.Context, ident id.ID) (*entities.Note, error) {
	note, err := r.cqrs.Queries.SelectNote(ctx, ident.Value())

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrNoteNotFound
	case err != nil:
		return nil, ex.Unexpected(err)
	}

	return note.Entity(), nil
}

func (r *NotesRepository) GetNotes(ctx context.Context) ([]*entities.Note, error) {
	notes, err := r.cqrs.Queries.SelectNotes(ctx)
	if err != nil {
		return nil, ex.Unexpected(err)
	}

	return queries.Notes(notes).Entities(), nil
}

func (r *NotesRepository) DeleteNote(ctx context.Context, note *entities.Note) error {
	err := r.cqrs.Commands.DeleteNote(ctx, note.ID.Value())
	if err != nil {
		return ex.Unexpected(err)
	}

	return nil
}
