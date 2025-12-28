package v1

import (
	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/domain"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
	pb "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UnmarshalRetrieveNoteRequest(request *pb.RetrieveNoteRequest) (*RetrieveNoteInput, error) {
	ident, err := id.Conv(request.GetId().GetValue())
	if err != nil {
		return nil, err
	}

	return &RetrieveNoteInput{ID: ident}, nil
}

func UnmarshalDeleteNoteRequest(request *pb.DeleteNoteRequest) (*DeleteNoteInput, error) {
	ident, err := id.Conv(request.GetId().GetValue())
	if err != nil {
		return nil, err
	}

	return &DeleteNoteInput{ID: ident}, nil
}

func UnmarshalCreateNoteRequest(request *pb.CreateNoteRequest) *CreateNoteInput {
	return &CreateNoteInput{
		Title:   request.GetTitle(),
		Content: request.GetContent(),
	}
}

func MarshalNote(note *entities.Note) *pb.Note {
	return &pb.Note{
		Id:        &typespb.ID{Value: note.ID.Value()},
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: timestamppb.New(note.CreatedAt),
		UpdatedAt: timestamppb.New(note.UpdatedAt),
	}
}

func MarshalNotes(notes []*entities.Note) []*pb.Note {
	pbNotes := make([]*pb.Note, len(notes))
	for i, note := range notes {
		pbNotes[i] = MarshalNote(note)
	}

	return pbNotes
}

var errorToCodeMapping = map[error]codes.Code{
	entities.ErrEmptyTitle:   codes.InvalidArgument,
	entities.ErrEmptyName:    codes.InvalidArgument,
	entities.ErrEmptyContent: codes.InvalidArgument,
	id.ErrNegativeID:         codes.InvalidArgument,
	domain.ErrNoteNotFound:   codes.NotFound,
	ex.ErrUnexpected:         codes.Internal,
}

func MarshalError(err error) error {
	code, ok := errorToCodeMapping[err]
	if !ok {
		code = codes.Unknown
	}

	text := err.Error()
	if code == codes.Internal {
		text = "internal error"
	}

	if code == codes.Unknown {
		text = "unknown error"
	}

	return status.Error(code, text)
}
