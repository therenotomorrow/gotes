package api

import (
	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/domain"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errorToCodeMapping = map[error]codes.Code{
	domain.ErrNoteNotFound:     codes.NotFound,
	domain.ErrPermissionDenied: codes.PermissionDenied,
	ex.ErrUnexpected:           codes.Internal,
}

var codeToErrorCodeMapping = map[codes.Code]typespb.ErrorCode{
	codes.NotFound:         typespb.ErrorCode_ERROR_CODE_ENTITY_NOT_FOUND,
	codes.PermissionDenied: typespb.ErrorCode_ERROR_CODE_PERMISSION_DENIED,
	codes.Internal:         typespb.ErrorCode_ERROR_CODE_INTERNAL,
	codes.Unknown:          typespb.ErrorCode_ERROR_CODE_UNKNOWN,
}

func Error(err error) error {
	text := ""
	details := new(typespb.Error)

	err, cause := ex.Expose(err)
	if err != nil {
		text = err.Error()
		details.Reason = text
	}

	if cause != nil {
		details.Reason = cause.Error()
	}

	code, ok := errorToCodeMapping[err]

	switch {
	case !ok:
		code = codes.Unknown
		text = "unknown error"
	case code == codes.Internal:
		text = "internal error"
	}

	details.Code = codeToErrorCodeMapping[code]

	st := status.New(code, text)
	st, err = st.WithDetails(details)

	ex.Skip(err)

	return st.Err()
}
