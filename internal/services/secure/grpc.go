package secure

import (
	"context"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (*entities.User, error)
}

func UnaryServerInterceptor(auth Authenticator, allowlist ...string) grpc.UnaryServerInterceptor {
	interceptor := authenticateInterceptor(auth, allowlist...)

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx, err := interceptor(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(auth Authenticator, allowlist ...string) grpc.StreamServerInterceptor {
	interceptor := authenticateInterceptor(auth, allowlist...)

	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, err := interceptor(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		wrapped := middleware.WrapServerStream(stream)
		wrapped.WrappedContext = ctx

		return handler(srv, wrapped)
	}
}

func authenticateInterceptor(
	auth Authenticator,
	allowlist ...string,
) func(ctx context.Context, method string) (context.Context, error) {
	allowed := make(map[string]struct{}, len(allowlist))
	for _, a := range allowlist {
		allowed[a] = struct{}{}
	}

	return func(ctx context.Context, method string) (context.Context, error) {
		if _, ok := allowed[method]; ok {
			return ctx, nil
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return ctx, status.Error(codes.Unauthenticated, "missing metadata")
		}

		vals := md.Get(authKey)
		if len(vals) == 0 {
			return ctx, status.Error(codes.Unauthenticated, "missing token")
		}

		user, err := auth.Authenticate(ctx, vals[0])
		if err != nil {
			return ctx, status.Error(codes.Unauthenticated, err.Error())
		}

		return context.WithValue(ctx, secureKey, user), nil
	}
}
