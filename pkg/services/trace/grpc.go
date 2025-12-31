package trace

import (
	"context"

	"google.golang.org/grpc"
)

func (t *Tracer) UnaryServerInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	ctx = t.Context(ctx)

	return handler(ctx, req)
}
