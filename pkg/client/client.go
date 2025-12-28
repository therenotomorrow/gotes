package client

import (
	"crypto/tls"

	"github.com/therenotomorrow/ex"
	pb "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	"github.com/therenotomorrow/gotes/pkg/validate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ErrInvalidConfig ex.Error = "invalid config"
)

type Config struct {
	Address string `json:"address" validate:"required"`
	Secure  bool   `json:"secure"`
}

type Client struct {
	pb.NotesServiceClient

	conn *grpc.ClientConn

	config Config
}

func New(cfg Config, options ...grpc.DialOption) (*Client, error) {
	err := validate.Struct(cfg)
	if err != nil {
		return nil, ErrInvalidConfig.Because(err)
	}

	creds := insecure.NewCredentials()
	if cfg.Secure {
		creds = credentials.NewTLS(new(tls.Config))
	}

	options = append(options, grpc.WithTransportCredentials(creds))

	conn, err := grpc.NewClient(cfg.Address, options...)
	if err != nil {
		return nil, ex.Unexpected(err)
	}

	client := pb.NewNotesServiceClient(conn)

	return &Client{NotesServiceClient: client, conn: conn, config: cfg}, nil
}

func MustNew(cfg Config, options ...grpc.DialOption) *Client {
	client, err := New(cfg, options...)
	ex.Panic(err)

	return client
}

func (c *Client) Close() {
	_ = c.conn.Close()
}
