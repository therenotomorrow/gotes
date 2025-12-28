package v1

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/api/notes/v1/ports"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
)

type UseCases struct {
	uow   ports.UnitOfWork
	store ports.Store
}

func New(uow ports.UnitOfWork, store ports.Store) *UseCases {
	return &UseCases{uow: uow, store: store}
}

type CreateNoteInput struct {
	Title   string
	Content string
}

func (use *UseCases) CreateNote(ctx context.Context, input *CreateNoteInput) (*entities.Note, error) {
	note, err := entities.NewNote(input.Title, input.Content)
	if err != nil {
		return nil, err
	}

	err = use.uow.Do(ctx, func(store ports.Store) error {
		var err error

		note, err = store.Notes().SaveNote(ctx, note)

		return err
	})

	return note, err
}

type DeleteNoteInput struct {
	ID id.ID
}

func (use *UseCases) DeleteNote(ctx context.Context, input *DeleteNoteInput) error {
	return use.uow.Do(ctx, func(store ports.Store) error {
		repo := store.Notes()

		note, err := repo.GetNote(ctx, input.ID)
		if err != nil {
			return err
		}

		return repo.DeleteNote(ctx, note)
	})
}

func (use *UseCases) ListNotes(ctx context.Context) ([]*entities.Note, error) {
	return use.store.Notes().GetNotes(ctx)
}

type RetrieveNoteInput struct {
	ID id.ID
}

func (use *UseCases) RetrieveNote(ctx context.Context, input *RetrieveNoteInput) (*entities.Note, error) {
	return use.store.Notes().GetNote(ctx, input.ID)
}
