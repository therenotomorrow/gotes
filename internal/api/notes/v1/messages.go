package v1

import (
	"context"

	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/usecases"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/services/secure"
	pb "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

type ErrorMarshaler struct {
	errorToCode      map[error]codes.Code
	errorToErrorCode map[error]typespb.ErrorCode
}

func NewErrorMarshaler() *ErrorMarshaler {
	return &ErrorMarshaler{
		errorToCode: map[error]codes.Code{
			usecases.ErrNoteNotFound:     codes.NotFound,
			usecases.ErrPermissionDenied: codes.PermissionDenied,
			context.Canceled:             codes.Canceled,
			ex.ErrUnexpected:             codes.Internal,
			secure.ErrUnauthorized:       codes.Unauthenticated,
		},
		errorToErrorCode: map[error]typespb.ErrorCode{
			usecases.ErrNoteNotFound:     typespb.ErrorCode_ERROR_CODE_ENTITY_NOT_FOUND,
			usecases.ErrPermissionDenied: typespb.ErrorCode_ERROR_CODE_PERMISSION_DENIED,
			context.Canceled:             typespb.ErrorCode_ERROR_CODE_INTERNAL,
			ex.ErrUnexpected:             typespb.ErrorCode_ERROR_CODE_INTERNAL,
			secure.ErrUnauthorized:       typespb.ErrorCode_ERROR_CODE_PERMISSION_DENIED,
		},
	}
}

func (e *ErrorMarshaler) Code(err error) (codes.Code, bool) {
	code, ok := e.errorToCode[err]

	return code, ok
}

func (e *ErrorMarshaler) ErrorCode(err error) (typespb.ErrorCode, bool) {
	errCode, ok := e.errorToErrorCode[err]

	return errCode, ok
}
