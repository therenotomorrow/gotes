package validate

import (
	"context"
	"errors"
	"strings"

	"buf.build/go/protovalidate"
	"github.com/therenotomorrow/ex"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
)

const (
	numParts = 2
)

func UnaryServerInterceptor(validator protovalidate.Validator) grpc.UnaryServerInterceptor {
	interceptor := validateInterceptor(validator)

	return func(ctx context.Context, request any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if msg, ok := request.(proto.Message); ok {
			err := interceptor(msg)
			if err != nil {
				return nil, err
			}
		}

		return handler(ctx, request)
	}
}

func UnaryClientInterceptor(validator protovalidate.Validator) grpc.UnaryClientInterceptor {
	interceptor := validateInterceptor(validator)

	return func(
		ctx context.Context,
		method string,
		request, reply any,
		conn *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if msg, ok := request.(proto.Message); ok {
			err := interceptor(msg)
			if err != nil {
				return err
			}
		}

		return invoker(ctx, method, request, reply, conn, opts...)
	}
}

func validateInterceptor(validator protovalidate.Validator) func(message proto.Message) error {
	return func(message proto.Message) error {
		err := validator.Validate(message)
		if err == nil {
			return nil
		}

		details := make([]protoadapt.MessageV1, 0)

		verr := new(protovalidate.ValidationError)
		if errors.As(err, &verr) {
			for _, violation := range verr.Violations {
				parts := strings.SplitN(violation.String(), ": ", numParts)

				name := parts[0]
				code := typespb.ErrorCode_ERROR_CODE_UNKNOWN

				switch name {
				case "id.value":
					code = typespb.ErrorCode_ERROR_CODE_INVALID_ID
				case "email":
					code = typespb.ErrorCode_ERROR_CODE_INVALID_EMAIL
				case "password":
					code = typespb.ErrorCode_ERROR_CODE_INVALID_PASSWORD
				case "title":
					code = typespb.ErrorCode_ERROR_CODE_INVALID_TITLE
				case "content":
					code = typespb.ErrorCode_ERROR_CODE_INVALID_CONTENT
				}

				details = append(details, &typespb.Error{Code: code, Reason: parts[1]})
			}
		}

		st := status.New(codes.InvalidArgument, err.Error())
		st, err = st.WithDetails(details...)

		ex.Skip(err)

		return st.Err()
	}
}
