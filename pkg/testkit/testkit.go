package testkit

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therenotomorrow/ex"
	typespb "github.com/therenotomorrow/gotes/pkg/api/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	bufSize = 1024 * 1024
)

type TestServer struct {
	listen *bufconn.Listener
	Server *grpc.Server
}

func NewTestServer(t *testing.T, options ...grpc.ServerOption) *TestServer {
	t.Helper()

	return &TestServer{
		listen: bufconn.Listen(bufSize),
		Server: grpc.NewServer(options...),
	}
}

func (s *TestServer) Client() *grpc.ClientConn {
	conn, err := grpc.NewClient(
		":0",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return s.listen.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	ex.Panic(err)

	return conn
}

func (s *TestServer) Serve() {
	err := s.Server.Serve(s.listen)

	ex.Panic(err)
}

func (s *TestServer) Stop() {
	s.Server.Stop()
}

func TimeAsTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t.Truncate(time.Minute))
}

func TruncateTimestamp(ts *timestamppb.Timestamp) *timestamppb.Timestamp {
	return TimeAsTimestamp(ts.AsTime())
}

func TimeByMinute(t time.Time) time.Time {
	return t.Truncate(time.Minute)
}

func NowByMinute() time.Time {
	return TimeByMinute(time.Now())
}

func AssertErrorDetails(t *testing.T, gotErr error, wantErr *typespb.Error) {
	t.Helper()

	st, ok := status.FromError(gotErr)
	assert.True(t, ok)
	assert.Len(t, st.Details(), 1)

	got, ok := st.Details()[0].(*typespb.Error)
	assert.True(t, ok)

	require.Equal(t, wantErr.GetCode(), got.GetCode())
	require.Equal(t, wantErr.GetReason(), got.GetReason())
}
