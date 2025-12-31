package adapters

import (
	"context"
	"database/sql"
	"errors"

	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/ports"
	"github.com/therenotomorrow/gotes/internal/domain"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
	commands "github.com/therenotomorrow/gotes/internal/storages/postgres/commands/notes"
	queries "github.com/therenotomorrow/gotes/internal/storages/postgres/queries/notes"
)

type NotesRepository struct {
	commands commands.Querier
	queries  queries.Querier
}

func NewNotesRepository(dbtx postgres.DBTX) *NotesRepository {
	return &NotesRepository{commands: commands.New(dbtx), queries: queries.New(dbtx)}
}

func (r *NotesRepository) SaveNote(ctx context.Context, note *entities.Note) (*entities.Note, error) {
	ident, err := r.commands.InsertNote(ctx, commands.NewInsertNoteParams(note))
	if err != nil {
		return nil, ex.Unexpected(err)
	}

	note.ID = id.New(ident)

	return note, nil
}

func (r *NotesRepository) GetNote(ctx context.Context, ident id.ID) (*entities.Note, error) {
	note, err := r.queries.SelectNote(ctx, ident.Value())

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrNoteNotFound
	case err != nil:
		return nil, ex.Unexpected(err)
	}

	return note.ToEntity(), nil
}

func (r *NotesRepository) GetNotesByUser(ctx context.Context, user *entities.User) ([]*entities.Note, error) {
	notes, err := r.queries.SelectNotesByUser(ctx, user.ID.ValuePtr())
	if err != nil {
		return nil, ex.Unexpected(err)
	}

	return queries.Notes(notes).ToEntities(), nil
}

func (r *NotesRepository) DeleteNote(ctx context.Context, note *entities.Note) error {
	err := r.commands.DeleteNote(ctx, note.ID.Value())

	return ex.Unexpected(err)
}

type Store struct {
	notes *NotesRepository
}

func NewStore(dbtx postgres.DBTX) *Store {
	return &Store{notes: NewNotesRepository(dbtx)}
}

func (s *Store) Notes() ports.NotesRepository {
	return s.notes
}

type UnitOfWork struct {
	database postgres.Database
}

func NewUnitOfWork(database postgres.Database) *UnitOfWork {
	return &UnitOfWork{database: database}
}

func (u *UnitOfWork) Do(ctx context.Context, work func(store ports.Store) error) error {
	return u.database.Tx(ctx, func(ctx context.Context) error {
		dbtx := u.database.Conn(ctx)
		store := NewStore(dbtx)

		return work(store)
	})
}
