package auth

import (
	"context"

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
	allowed := make(map[string]struct{}, len(allowlist))
	for _, a := range allowlist {
		allowed[a] = struct{}{}
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if _, ok := allowed[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		vals := md.Get(authKey)
		if len(vals) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		user, err := auth.Authenticate(ctx, vals[0])
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return handler(context.WithValue(ctx, secureKey, user), req)
	}
}
