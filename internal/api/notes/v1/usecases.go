package v1

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/api/notes/v1/ports"
	"github.com/therenotomorrow/gotes/internal/domain"
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

func (use *UseCases) CreateNote(
	ctx context.Context,
	user *entities.User,
	input *CreateNoteInput,
) (*entities.Note, error) {
	note, err := entities.NewNote(input.Title, input.Content)
	if err != nil {
		return nil, err
	}

	note.SetOwner(user)

	err = use.uow.Do(ctx, func(store ports.Store) error {
		note, err = store.Notes().SaveNote(ctx, note)

		return err
	})

	return note, err
}

type DeleteNoteInput struct {
	ID int64
}

func (use *UseCases) DeleteNote(ctx context.Context, user *entities.User, input *DeleteNoteInput) error {
	ident, err := id.Conv(input.ID)
	if err != nil {
		return err
	}

	return use.uow.Do(ctx, func(store ports.Store) error {
		repo := store.Notes()

		note, err := repo.GetNote(ctx, ident)
		if err != nil {
			return err
		}

		err = use.permit(user, note)
		if err != nil {
			return err
		}

		return repo.DeleteNote(ctx, note)
	})
}

func (use *UseCases) ListNotes(ctx context.Context, user *entities.User) ([]*entities.Note, error) {
	return use.store.Notes().GetNotesByUser(ctx, user)
}

type RetrieveNoteInput struct {
	ID int64
}

func (use *UseCases) RetrieveNote(
	ctx context.Context,
	user *entities.User,
	input *RetrieveNoteInput,
) (*entities.Note, error) {
	ident, err := id.Conv(input.ID)
	if err != nil {
		return nil, err
	}

	note, err := use.store.Notes().GetNote(ctx, ident)
	if err != nil {
		return nil, err
	}

	return note, use.permit(user, note)
}

func (use *UseCases) permit(user *entities.User, note *entities.Note) error {
	if !note.IsOwner(user) {
		return domain.ErrPermissionDenied
	}

	return nil
}
