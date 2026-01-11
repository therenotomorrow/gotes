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

type Store struct {
	Notes NotesRepository
}

type StoreProvider interface {
	Provide(ctx context.Context) Store
}

type UnitOfWork interface {
	Do(ctx context.Context, work func(store Store) error) error
}
