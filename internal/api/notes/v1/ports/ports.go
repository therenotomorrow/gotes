package ports

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
)

type NotesRepository interface {
	SaveNote(ctx context.Context, note *entities.Note) (*entities.Note, error)
	GetNote(ctx context.Context, id id.ID) (*entities.Note, error)
	DeleteNote(ctx context.Context, note *entities.Note) error
	GetNotesByUser(ctx context.Context, user *entities.User) ([]*entities.Note, error)
}

type Store interface {
	Notes() NotesRepository
}

type UnitOfWork interface {
	Do(ctx context.Context, unit func(store Store) error) error
}
