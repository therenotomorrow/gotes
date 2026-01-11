package v1

import (
	"context"
	"log/slog"

	"github.com/therenotomorrow/gotes/internal/api"
	adapters "github.com/therenotomorrow/gotes/internal/api/notes/v1/adapters/postgres"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/ports"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/usecases"
	"github.com/therenotomorrow/gotes/internal/services/secure"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
	pb "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
)

type NotesService struct {
	pb.UnimplementedNotesServiceServer

	handle api.ErrorHandlerFunc
	tracer *trace.Tracer
	cases  *usecases.UseCases
}

func NewService(db postgres.Database, logger *slog.Logger) *NotesService {
	provider := adapters.NewStoreProvider(db)
	uow := adapters.NewUnitOfWork(db, provider)

	return NewServiceWithProvider(uow, provider, logger)
}

func NewServiceWithProvider(uow ports.UnitOfWork, provider ports.StoreProvider, logger *slog.Logger) *NotesService {
	store := provider.Provide(context.Background())

	return &NotesService{
		UnimplementedNotesServiceServer: pb.UnimplementedNotesServiceServer{},
		handle:                          api.ErrorHandler(NewErrorMarshaler()),
		tracer:                          trace.Service("notes.v1", logger),
		cases:                           usecases.NewCases(uow, store),
	}
}

func (svc *NotesService) CreateNote(
	ctx context.Context,
	request *pb.CreateNoteRequest,
) (*pb.CreateNoteResponse, error) {
	user, err := secure.User(ctx)
	if err != nil {
		return nil, svc.handle(err)
	}

	note, err := svc.cases.CreateNote(ctx, user, &usecases.CreateNoteInput{
		Title:   request.GetTitle(),
		Content: request.GetContent(),
	})
	if err != nil {
		svc.tracer.Error(ctx, "CreateNote", err, "user", user.ID)

		return nil, svc.handle(err)
	}

	return &pb.CreateNoteResponse{Note: MarshalNote(note)}, nil
}

func (svc *NotesService) DeleteNote(
	ctx context.Context,
	request *pb.DeleteNoteRequest,
) (*pb.DeleteNoteResponse, error) {
	user, err := secure.User(ctx)
	if err != nil {
		return nil, svc.handle(err)
	}

	err = svc.cases.DeleteNote(ctx, user, &usecases.DeleteNoteInput{
		ID: request.GetId().GetValue(),
	})
	if err != nil {
		return nil, svc.handle(err)
	}

	return &pb.DeleteNoteResponse{}, nil
}

func (svc *NotesService) RetrieveNote(
	ctx context.Context,
	request *pb.RetrieveNoteRequest,
) (*pb.RetrieveNoteResponse, error) {
	user, err := secure.User(ctx)
	if err != nil {
		return nil, svc.handle(err)
	}

	note, err := svc.cases.RetrieveNote(ctx, user, &usecases.RetrieveNoteInput{
		ID: request.GetId().GetValue(),
	})
	if err != nil {
		return nil, svc.handle(err)
	}

	return &pb.RetrieveNoteResponse{Note: MarshalNote(note)}, nil
}

func (svc *NotesService) ListNotes(ctx context.Context, _ *pb.ListNotesRequest) (*pb.ListNotesResponse, error) {
	user, err := secure.User(ctx)
	if err != nil {
		return nil, svc.handle(err)
	}

	notes, err := svc.cases.ListNotes(ctx, user)
	if err != nil {
		return nil, svc.handle(err)
	}

	return &pb.ListNotesResponse{Notes: MarshalNotes(notes)}, nil
}
