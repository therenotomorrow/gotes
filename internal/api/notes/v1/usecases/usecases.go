package usecases

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/api/notes/v1/ports"
	"github.com/therenotomorrow/gotes/internal/domain"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
)

const (
	ErrNoteNotFound     domain.Error = "note not found"
	ErrZeroEvents       domain.Error = "zero events"
	ErrPermissionDenied domain.Error = "permission denied"
)

type UseCases struct {
	uow   ports.UnitOfWork
	store ports.Store
}

func NewCases(uow ports.UnitOfWork, store ports.Store) *UseCases {
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
		note, err = store.Notes.SaveNote(ctx, note)
		if err != nil {
			return err
		}

		event := entities.NewEvent(entities.EventTypeCreated, note)

		return store.Events.SaveEvent(ctx, event)
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
		note, err := store.Notes.GetNote(ctx, ident)
		if err != nil {
			return err
		}

		err = use.permit(user, note)
		if err != nil {
			return err
		}

		err = store.Notes.DeleteNote(ctx, note)
		if err != nil {
			return err
		}

		event := entities.NewEvent(entities.EventTypeDeleted, note)

		return store.Events.SaveEvent(ctx, event)
	})
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

	note, err := use.store.Notes.GetNote(ctx, ident)
	if err != nil {
		return nil, err
	}

	err = use.permit(user, note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (use *UseCases) ListNotes(ctx context.Context, user *entities.User) ([]*entities.Note, error) {
	return use.store.Notes.GetNotesByUser(ctx, user)
}

func (use *UseCases) UnreadEvents(ctx context.Context, user *entities.User) (int32, error) {
	return use.store.Events.CountEvents(ctx, user)
}

func (use *UseCases) GetNextEvent(ctx context.Context, user *entities.User) (*entities.Event, error) {
	return use.store.Events.GetEvent(ctx, user)
}

func (use *UseCases) permit(user *entities.User, note *entities.Note) error {
	if !note.IsOwner(user) {
		return ErrPermissionDenied
	}

	return nil
}
