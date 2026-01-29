package usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/adapters/mocks"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/ports"
	v1 "github.com/therenotomorrow/gotes/internal/api/notes/v1/usecases"
	"github.com/therenotomorrow/gotes/internal/api/users/v1/usecases"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
	"github.com/therenotomorrow/gotes/internal/domain/types/uuid"
	"github.com/therenotomorrow/gotes/pkg/services/generate"
	"github.com/therenotomorrow/gotes/pkg/testkit"
)

type unitOfWork ports.Store

func (uow unitOfWork) Do(_ context.Context, work func(store ports.Store) error) error {
	return work(ports.Store(uow))
}

func TestUseCasesCreateNote(t *testing.T) {
	t.Parallel()

	uuid.SetGenerator(generate.NewUUID())

	t.Run("new note errors", func(t *testing.T) {
		t.Parallel()

		var (
			ctx   = t.Context()
			user  = new(entities.User)
			store = ports.Store{}
			use   = v1.NewCases(nil, store)
		)

		got, err := use.CreateNote(ctx, user, &v1.CreateNoteInput{
			Title:   "",
			Content: "content",
		})
		require.ErrorIs(t, err, entities.ErrEmptyTitle)
		assert.Nil(t, got)

		got, err = use.CreateNote(ctx, user, &v1.CreateNoteInput{
			Title:   "title",
			Content: "",
		})
		require.ErrorIs(t, err, entities.ErrEmptyContent)
		assert.Nil(t, got)
	})

	t.Run("store error", func(t *testing.T) {
		t.Parallel()

		var (
			ctx   = t.Context()
			user  = new(entities.User)
			input = &v1.CreateNoteInput{Title: "title", Content: "content"}
			notes = mocks.NewMockNotesRepository(t)
			store = ports.Store{Notes: notes}
			use   = v1.NewCases(unitOfWork(store), store)
		)

		notes.On("SaveNote", ctx, mock.AnythingOfType("*entities.Note")).
			Return(nil, ex.ErrUnknown)

		got, err := use.CreateNote(ctx, user, input)
		require.ErrorIs(t, err, ex.ErrUnknown)
		assert.Nil(t, got)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var (
			ctx    = t.Context()
			user   = new(entities.User)
			input  = &v1.CreateNoteInput{Title: "title", Content: "content"}
			notes  = mocks.NewMockNotesRepository(t)
			events = mocks.NewMockEventsRepository(t)
			store  = ports.Store{Notes: notes, Events: events}
			use    = v1.NewCases(unitOfWork(store), store)
		)

		ident := id.New(42)

		notes.On("SaveNote", ctx, mock.AnythingOfType("*entities.Note")).
			Return(func(ctx context.Context, note *entities.Note) *entities.Note {
				note.ID = ident

				return note
			}, nil)
		events.On("SaveEvent", ctx, mock.AnythingOfType("*entities.Event")).
			Return(nil)

		got, err := use.CreateNote(ctx, user, input)
		require.NoError(t, err)

		got.CreatedAt = testkit.TimeByMinute(got.CreatedAt)
		got.UpdatedAt = testkit.TimeByMinute(got.UpdatedAt)

		now := testkit.NowByMinute()
		want := &entities.Note{
			CreatedAt: now,
			UpdatedAt: now,
			Owner:     user,
			Title:     "title",
			Content:   "content",
			ID:        ident,
		}

		assert.Equal(t, want, got)
	})
}

func TestUseCasesDeleteNote(t *testing.T) {
	t.Parallel()

	uuid.SetGenerator(generate.NewUUID())

	t.Run("validation errors", func(t *testing.T) {
		t.Parallel()

		var (
			ctx   = t.Context()
			user  = new(entities.User)
			store = ports.Store{}
			use   = v1.NewCases(nil, store)
		)

		err := use.DeleteNote(ctx, user, &v1.DeleteNoteInput{
			ID: -42,
		})
		require.ErrorIs(t, err, id.ErrInvalidID)
	})

	t.Run("store get note error", func(t *testing.T) {
		t.Parallel()

		var (
			ctx   = t.Context()
			user  = new(entities.User)
			ident = id.New(42)
			input = &v1.DeleteNoteInput{ID: ident.Value()}
			notes = mocks.NewMockNotesRepository(t)
			store = ports.Store{Notes: notes}
			use   = v1.NewCases(unitOfWork(store), store)
		)

		notes.On("GetNote", ctx, ident).
			Return(nil, ex.ErrUnknown)

		err := use.DeleteNote(ctx, user, input)
		require.ErrorIs(t, err, ex.ErrUnknown)
	})

	t.Run("not permit", func(t *testing.T) {
		t.Parallel()

		var (
			ctx   = t.Context()
			ident = id.New(42)
			input = &v1.DeleteNoteInput{ID: ident.Value()}
			notes = mocks.NewMockNotesRepository(t)
			store = ports.Store{Notes: notes}
			use   = v1.NewCases(unitOfWork(store), store)
		)

		user := &entities.User{ID: id.New(20)}
		note := &entities.Note{Owner: &entities.User{ID: id.New(10)}}

		notes.On("GetNote", ctx, ident).
			Return(note, nil)

		err := use.DeleteNote(ctx, user, input)
		require.ErrorIs(t, err, usecases.ErrPermissionDenied)
	})

	t.Run("store delete note error", func(t *testing.T) {
		t.Parallel()

		var (
			ctx   = t.Context()
			owner = &entities.User{ID: id.New(30)}
			note  = &entities.Note{Owner: owner}
			ident = id.New(42)
			input = &v1.DeleteNoteInput{ID: ident.Value()}
			notes = mocks.NewMockNotesRepository(t)
			store = ports.Store{Notes: notes}
			use   = v1.NewCases(unitOfWork(store), store)
		)

		notes.On("GetNote", ctx, ident).
			Return(note, nil)
		notes.On("DeleteNote", ctx, note).
			Return(ex.ErrUnknown)

		err := use.DeleteNote(ctx, owner, input)
		require.ErrorIs(t, err, ex.ErrUnknown)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var (
			ctx    = t.Context()
			owner  = &entities.User{ID: id.New(40)}
			note   = &entities.Note{Owner: owner}
			ident  = id.New(42)
			input  = &v1.DeleteNoteInput{ID: ident.Value()}
			notes  = mocks.NewMockNotesRepository(t)
			events = mocks.NewMockEventsRepository(t)
			store  = ports.Store{Notes: notes, Events: events}
			use    = v1.NewCases(unitOfWork(store), store)
		)

		notes.On("GetNote", ctx, ident).
			Return(note, nil)
		notes.On("DeleteNote", ctx, note).
			Return(nil)
		events.On("SaveEvent", ctx, mock.AnythingOfType("*entities.Event")).
			Return(nil)

		err := use.DeleteNote(ctx, owner, input)
		require.NoError(t, err)
	})
}

func TestUseCasesRetrieveNote(t *testing.T) {
	t.Parallel()

	t.Run("validation errors", func(t *testing.T) {
		t.Parallel()

		var (
			ctx   = t.Context()
			user  = new(entities.User)
			store = ports.Store{}
			use   = v1.NewCases(nil, store)
		)

		got, err := use.RetrieveNote(ctx, user, &v1.RetrieveNoteInput{
			ID: -42,
		})
		require.ErrorIs(t, err, id.ErrInvalidID)
		assert.Nil(t, got)
	})

	t.Run("store error", func(t *testing.T) {
		t.Parallel()

		var (
			ctx   = t.Context()
			user  = new(entities.User)
			ident = id.New(42)
			input = &v1.RetrieveNoteInput{ID: ident.Value()}
			notes = mocks.NewMockNotesRepository(t)
			store = ports.Store{Notes: notes}
			use   = v1.NewCases(unitOfWork(store), store)
		)

		notes.On("GetNote", ctx, ident).
			Return(nil, ex.ErrUnknown)

		got, err := use.RetrieveNote(ctx, user, input)
		require.ErrorIs(t, err, ex.ErrUnknown)
		assert.Nil(t, got)
	})

	t.Run("not permit", func(t *testing.T) {
		t.Parallel()

		var (
			ctx   = t.Context()
			ident = id.New(42)
			input = &v1.RetrieveNoteInput{ID: ident.Value()}
			notes = mocks.NewMockNotesRepository(t)
			store = ports.Store{Notes: notes}
			use   = v1.NewCases(unitOfWork(store), store)
		)

		user := &entities.User{ID: id.New(20)}
		note := &entities.Note{Owner: &entities.User{ID: id.New(10)}}

		notes.On("GetNote", ctx, ident).
			Return(note, nil)

		got, err := use.RetrieveNote(ctx, user, input)
		require.ErrorIs(t, err, usecases.ErrPermissionDenied)
		assert.Nil(t, got)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var (
			ctx   = t.Context()
			owner = &entities.User{ID: id.New(40)}
			note  = &entities.Note{Owner: owner}
			ident = id.New(42)
			input = &v1.RetrieveNoteInput{ID: ident.Value()}
			notes = mocks.NewMockNotesRepository(t)
			store = ports.Store{Notes: notes}
			use   = v1.NewCases(unitOfWork(store), store)
		)

		notes.On("GetNote", ctx, ident).
			Return(note, nil)

		got, err := use.RetrieveNote(ctx, owner, input)
		require.NoError(t, err)

		want := note

		assert.Equal(t, want, got)
	})
}

func TestUseCasesListNotes(t *testing.T) {
	t.Parallel()

	t.Run("store error", func(t *testing.T) {
		t.Parallel()

		var (
			ctx   = t.Context()
			user  = new(entities.User)
			notes = mocks.NewMockNotesRepository(t)
			store = ports.Store{Notes: notes}
			use   = v1.NewCases(unitOfWork(store), store)
		)

		notes.On("GetNotesByUser", ctx, user).
			Return(nil, ex.ErrUnknown)

		got, err := use.ListNotes(ctx, user)
		require.ErrorIs(t, err, ex.ErrUnknown)
		assert.Nil(t, got)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var (
			ctx   = t.Context()
			user  = new(entities.User)
			notes = mocks.NewMockNotesRepository(t)
			store = ports.Store{Notes: notes}
			use   = v1.NewCases(unitOfWork(store), store)
		)

		notes.On("GetNotesByUser", ctx, user).
			Return([]*entities.Note{new(entities.Note), new(entities.Note)}, nil)

		got, err := use.ListNotes(ctx, user)
		require.NoError(t, err)

		want := []*entities.Note{new(entities.Note), new(entities.Note)}

		assert.Equal(t, want, got)
	})
}
