package trace

import (
	"context"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
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

func (t *Tracer) StreamServerInterceptor(
	srv any,
	stream grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	wrapped := middleware.WrapServerStream(stream)
	wrapped.WrappedContext = t.Context(stream.Context())

	return handler(srv, wrapped)
}
