package v1

import (
	"context"
	"log/slog"

	"github.com/therenotomorrow/gotes/internal/api"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/adapters"
	"github.com/therenotomorrow/gotes/internal/services/auth"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
	pb "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	"github.com/therenotomorrow/gotes/pkg/services/trace"
)

type NotesService struct {
	pb.UnimplementedNotesServiceServer

	secure auth.Securable
	tracer *trace.Tracer
	cases  *UseCases
}

func Service(secure auth.Securable, database postgres.Database, logger *slog.Logger) *NotesService {
	var (
		uow    = adapters.NewUnitOfWork(database)
		store  = adapters.NewStore(database.Conn(context.Background()))
		tracer = trace.Service("notes.v1", logger)
	)

	return &NotesService{
		UnimplementedNotesServiceServer: pb.UnimplementedNotesServiceServer{},
		secure:                          secure,
		tracer:                          tracer,
		cases:                           New(uow, store),
	}
}

func (svc *NotesService) ListNotes(ctx context.Context, _ *pb.ListNotesRequest) (*pb.ListNotesResponse, error) {
	notes, err := svc.cases.ListNotes(ctx, svc.secure.User(ctx))
	if err != nil {
		return nil, api.Error(err)
	}

	return &pb.ListNotesResponse{Notes: MarshalNotes(notes)}, nil
}

func (svc *NotesService) RetrieveNote(
	ctx context.Context,
	request *pb.RetrieveNoteRequest,
) (*pb.RetrieveNoteResponse, error) {
	note, err := svc.cases.RetrieveNote(ctx, svc.secure.User(ctx), &RetrieveNoteInput{
		ID: request.GetId().GetValue(),
	})
	if err != nil {
		return nil, api.Error(err)
	}

	return &pb.RetrieveNoteResponse{Note: MarshalNote(note)}, nil
}

func (svc *NotesService) CreateNote(
	ctx context.Context,
	request *pb.CreateNoteRequest,
) (*pb.CreateNoteResponse, error) {
	note, err := svc.cases.CreateNote(ctx, svc.secure.User(ctx), &CreateNoteInput{
		Title:   request.GetTitle(),
		Content: request.GetContent(),
	})
	if err != nil {
		return nil, api.Error(err)
	}

	return &pb.CreateNoteResponse{Note: MarshalNote(note)}, nil
}

func (svc *NotesService) DeleteNote(
	ctx context.Context,
	request *pb.DeleteNoteRequest,
) (*pb.DeleteNoteResponse, error) {
	err := svc.cases.DeleteNote(ctx, svc.secure.User(ctx), &DeleteNoteInput{
		ID: request.GetId().GetValue(),
	})
	if err != nil {
		return nil, api.Error(err)
	}

	return &pb.DeleteNoteResponse{}, nil
}
