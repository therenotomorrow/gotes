package v1_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/therenotomorrow/ex"
	v1 "github.com/therenotomorrow/gotes/internal/api/notes/v1"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/adapters/mocks"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/ports"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/usecases"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
	"github.com/therenotomorrow/gotes/internal/services/secure"
	pb "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
	"github.com/therenotomorrow/gotes/pkg/testkit"
)

type unitOfWork struct {
	provider ports.StoreProvider
}

func (uow unitOfWork) Do(ctx context.Context, work func(store ports.Store) error) error {
	return work(uow.provider.Provide(ctx))
}

var (
	log   = trace.Logger(trace.TEXT, true)
	input = &pb.CreateNoteRequest{
		Title:   "title",
		Content: "content",
	}
)

func TestNotesService_CreateNote(t *testing.T) {
	t.Parallel()

	t.Run("secure", func(t *testing.T) {
		t.Parallel()

		var (
			ctx      = t.Context()
			notes    = mocks.NewMockNotesRepository(t)
			store    = ports.Store{Notes: notes}
			provider = mocks.NewMockStoreProvider(t)
		)

		provider.On("Provide", context.Background()).Return(store)

		svc := v1.NewServiceWithProvider(unitOfWork{provider: provider}, provider, log)

		resp, err := svc.CreateNote(ctx, input)
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = unauthorized")
		assert.Nil(t, resp)

		testkit.AssertErrorDetails(t, err, &typespb.Error{
			Code:   typespb.ErrorCode_ERROR_CODE_PERMISSION_DENIED,
			Reason: "missing user in context",
		})
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()

		var (
			ctx      = t.Context()
			user     = new(entities.User)
			notes    = mocks.NewMockNotesRepository(t)
			store    = ports.Store{Notes: notes}
			provider = mocks.NewMockStoreProvider(t)
		)

		ctx = secure.NewUserContext(ctx, user)

		provider.On("Provide", context.Background()).Return(store)
		provider.On("Provide", ctx).Return(store)
		notes.On("SaveNote", ctx, mock.AnythingOfType("*entities.Note")).
			Return(nil, ex.Conv(usecases.ErrNoteNotFound).Reason("test error"))

		svc := v1.NewServiceWithProvider(unitOfWork{provider: provider}, provider, log)

		resp, err := svc.CreateNote(ctx, input)
		require.EqualError(t, err, "rpc error: code = NotFound desc = note not found")
		assert.Nil(t, resp)

		testkit.AssertErrorDetails(t, err, &typespb.Error{
			Code:   typespb.ErrorCode_ERROR_CODE_ENTITY_NOT_FOUND,
			Reason: "test error",
		})
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		var (
			ctx      = t.Context()
			user     = new(entities.User)
			notes    = mocks.NewMockNotesRepository(t)
			store    = ports.Store{Notes: notes}
			provider = mocks.NewMockStoreProvider(t)
		)

		ident := id.New(42)
		ctx = secure.NewUserContext(ctx, user)

		provider.On("Provide", context.Background()).Return(store)
		provider.On("Provide", ctx).Return(store)

		notes.On("SaveNote", ctx, mock.AnythingOfType("*entities.Note")).
			Return(func(ctx context.Context, note *entities.Note) *entities.Note {
				note.ID = ident

				return note
			}, nil)

		svc := v1.NewServiceWithProvider(unitOfWork{provider: provider}, provider, log)

		resp, err := svc.CreateNote(ctx, input)
		require.NoError(t, err)

		resp.Note.CreatedAt = testkit.TruncateTimestamp(resp.GetNote().GetCreatedAt())
		resp.Note.UpdatedAt = testkit.TruncateTimestamp(resp.GetNote().GetUpdatedAt())

		now := testkit.NowByMinute()
		want := &pb.CreateNoteResponse{Note: &pb.Note{
			Id:        &typespb.ID{Value: ident.Value()},
			Title:     "title",
			Content:   "content",
			CreatedAt: testkit.TimeAsTimestamp(now),
			UpdatedAt: testkit.TimeAsTimestamp(now),
		}}

		assert.Equal(t, want, resp)
	})
}
