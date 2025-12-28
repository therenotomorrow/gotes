package v1

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/api/notes/v1/adapters"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
	pb "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
)

type NotesService struct {
	pb.UnimplementedNotesServiceServer

	cases *UseCases
}

func Service(database *postgres.Database) *NotesService {
	uow := adapters.NewUnitOfWork(database)
	store := adapters.NewStore(database.CQRS(context.Background()))
	cases := New(uow, store)

	return &NotesService{
		UnimplementedNotesServiceServer: pb.UnimplementedNotesServiceServer{},
		cases:                           cases,
	}
}

func (svc *NotesService) ListNotes(ctx context.Context, _ *pb.ListNotesRequest) (*pb.ListNotesResponse, error) {
	notes, err := svc.cases.ListNotes(ctx)
	if err != nil {
		return nil, MarshalError(err)
	}

	return &pb.ListNotesResponse{Notes: MarshalNotes(notes)}, nil
}

func (svc *NotesService) RetrieveNote(
	ctx context.Context,
	request *pb.RetrieveNoteRequest,
) (*pb.RetrieveNoteResponse, error) {
	input, err := UnmarshalRetrieveNoteRequest(request)
	if err != nil {
		return nil, MarshalError(err)
	}

	note, err := svc.cases.RetrieveNote(ctx, input)
	if err != nil {
		return nil, MarshalError(err)
	}

	return &pb.RetrieveNoteResponse{Note: MarshalNote(note)}, nil
}

func (svc *NotesService) CreateNote(
	ctx context.Context,
	request *pb.CreateNoteRequest,
) (*pb.CreateNoteResponse, error) {
	note, err := svc.cases.CreateNote(ctx, UnmarshalCreateNoteRequest(request))
	if err != nil {
		return nil, MarshalError(err)
	}

	return &pb.CreateNoteResponse{Note: MarshalNote(note)}, nil
}

func (svc *NotesService) DeleteNote(
	ctx context.Context,
	request *pb.DeleteNoteRequest,
) (*pb.DeleteNoteResponse, error) {
	input, err := UnmarshalDeleteNoteRequest(request)
	if err != nil {
		return nil, MarshalError(err)
	}

	err = svc.cases.DeleteNote(ctx, input)
	if err != nil {
		return nil, MarshalError(err)
	}

	return &pb.DeleteNoteResponse{}, nil
}
