package v1

import (
	"context"

	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/api/chat/v1/entities"
	pb "github.com/therenotomorrow/gotes/pkg/api/chat/v1"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func MarshalMessage(message *entities.Message) *pb.DispatchResponse_Message {
	return &pb.DispatchResponse_Message{
		Message: &pb.Message{
			Header: &pb.Header{CorrelationId: message.Header.CorrelationID},
			Text:   message.Text,
		},
	}
}

func MarshalDetails(st *status.Status) []*anypb.Any {
	details := make([]*anypb.Any, 0)

	for _, detail := range st.Details() {
		if msg, ok := detail.(proto.Message); ok {
			anyDetail, err := anypb.New(msg)
			if err == nil {
				details = append(details, anyDetail)
			}
		}
	}

	return details
}

func MarshalStatus(st *status.Status) *pb.DispatchResponse_Status {
	return &pb.DispatchResponse_Status{
		Status: &statuspb.Status{
			Code:    int32(st.Code()), //nolint:gosec // allowed conversation
			Message: st.Message(),
			Details: MarshalDetails(st),
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
			context.Canceled: codes.Canceled,
			ex.ErrUnexpected: codes.Internal,
			ErrChat:          codes.Unavailable,
		},
		errorToErrorCode: map[error]typespb.ErrorCode{
			context.Canceled: typespb.ErrorCode_ERROR_CODE_INTERNAL,
			ex.ErrUnexpected: typespb.ErrorCode_ERROR_CODE_INTERNAL,
			ErrChat:          typespb.ErrorCode_ERROR_CODE_INTERNAL,
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
