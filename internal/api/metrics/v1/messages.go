package v1

import (
	"context"

	"github.com/therenotomorrow/ex"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	"google.golang.org/grpc/codes"
)

type ErrorMarshaler struct {
	errorToCode      map[error]codes.Code
	errorToErrorCode map[error]typespb.ErrorCode
}

func NewErrorMarshaler() *ErrorMarshaler {
	return &ErrorMarshaler{
		errorToCode: map[error]codes.Code{
			context.Canceled: codes.Canceled,
			ex.ErrUnexpected: codes.Internal,
			ErrStream:        codes.Unavailable,
		},
		errorToErrorCode: map[error]typespb.ErrorCode{
			context.Canceled: typespb.ErrorCode_ERROR_CODE_INTERNAL,
			ex.ErrUnexpected: typespb.ErrorCode_ERROR_CODE_INTERNAL,
			ErrStream:        typespb.ErrorCode_ERROR_CODE_INTERNAL,
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
