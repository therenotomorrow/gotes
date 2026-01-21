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

func MarshalUnread(unread int32) *pb.SubscribeToEventsResponse_Unread {
	return &pb.SubscribeToEventsResponse_Unread{
		Unread: &pb.Unread{Events: unread},
	}
}

func MarshalEvent(event *entities.Event) *pb.SubscribeToEventsResponse_Event {
	var eventType pb.EventType

	switch event.EventType {
	case entities.EventTypeCreated:
		eventType = pb.EventType_EVENT_TYPE_CREATED
	case entities.EventTypeDeleted:
		eventType = pb.EventType_EVENT_TYPE_DELETED
	default:
		eventType = pb.EventType_EVENT_TYPE_UNKNOWN
	}

	return &pb.SubscribeToEventsResponse_Event{
		Event: &pb.Event{
			Id:        event.ID.Value(),
			Type:      eventType,
			NoteId:    &typespb.ID{Value: event.Note.ID.Value()},
			EventTime: timestamppb.New(event.EventTime),
		},
	}
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
			ErrSend:                      codes.Unavailable,
		},
		errorToErrorCode: map[error]typespb.ErrorCode{
			usecases.ErrNoteNotFound:     typespb.ErrorCode_ERROR_CODE_ENTITY_NOT_FOUND,
			usecases.ErrPermissionDenied: typespb.ErrorCode_ERROR_CODE_PERMISSION_DENIED,
			context.Canceled:             typespb.ErrorCode_ERROR_CODE_INTERNAL,
			ex.ErrUnexpected:             typespb.ErrorCode_ERROR_CODE_INTERNAL,
			secure.ErrUnauthorized:       typespb.ErrorCode_ERROR_CODE_PERMISSION_DENIED,
			ErrSend:                      typespb.ErrorCode_ERROR_CODE_INTERNAL,
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
