package api

import (
	"github.com/therenotomorrow/ex"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorMarshaler interface {
	Code(err error) (codes.Code, bool)
	ErrorCode(err error) (typespb.ErrorCode, bool)
}

type ErrorHandlerFunc func(err error) error

func ErrorHandler(marshaler ErrorMarshaler) ErrorHandlerFunc {
	return func(err error) error {
		err, cause := ex.Expose(err)

		var text, reason string
		if err != nil {
			text = err.Error()
			reason = text
		}

		if cause != nil {
			reason = cause.Error()
		}

		code, exist := marshaler.Code(err)
		switch {
		case !exist:
			code = codes.Unknown
			text = "unknown error"
		case code == codes.Internal:
			// hide any details that are only internal
			text = "internal error"
		}

		errorCode, exist := marshaler.ErrorCode(err)
		if !exist {
			errorCode = typespb.ErrorCode_ERROR_CODE_UNKNOWN
		}

		st, err := status.New(code, text).WithDetails(&typespb.Error{
			Code:   errorCode,
			Reason: reason,
		})

		ex.Skip(err)

		return st.Err()
	}
}
