package v1

import (
	"context"

	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/api/users/v1/usecases"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	pb "github.com/therenotomorrow/gotes/pkg/api/users/v1"
	"google.golang.org/grpc/codes"
)

func MarshalUser(user *entities.User) *pb.User {
	return &pb.User{
		Id:    &typespb.ID{Value: user.ID.Value()},
		Name:  user.Name,
		Email: user.Email.Value(),
		Token: user.Token.Value(),
	}
}

type ErrorMarshaler struct {
	errorToCode      map[error]codes.Code
	errorToErrorCode map[error]typespb.ErrorCode
}

func NewErrorMarshaler() *ErrorMarshaler {
	return &ErrorMarshaler{
		errorToCode: map[error]codes.Code{
			usecases.ErrPermissionDenied: codes.PermissionDenied,
			context.Canceled:             codes.Canceled,
			ex.ErrUnexpected:             codes.Internal,
		},
		errorToErrorCode: map[error]typespb.ErrorCode{
			usecases.ErrPermissionDenied: typespb.ErrorCode_ERROR_CODE_PERMISSION_DENIED,
			context.Canceled:             typespb.ErrorCode_ERROR_CODE_INTERNAL,
			ex.ErrUnexpected:             typespb.ErrorCode_ERROR_CODE_INTERNAL,
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
