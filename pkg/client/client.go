package client

import (
	"context"
	"crypto/tls"

	"buf.build/go/protovalidate"
	"github.com/therenotomorrow/ex"
	notesv1 "github.com/therenotomorrow/gotes/pkg/api/notes/v1"
	usersv1 "github.com/therenotomorrow/gotes/pkg/api/users/v1"
	"github.com/therenotomorrow/gotes/pkg/services/validate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	authKey = "authorization"

	ErrInvalidConfig ex.Error = "invalid config"
)

type Config struct {
	Address string `json:"address" validate:"required"`
	Secure  bool   `json:"secure"`
}

type Client struct {
	notesv1.NotesServiceClient
	usersv1.UsersServiceClient

	conn   *grpc.ClientConn
	config Config
}

func New(cfg Config, options ...grpc.DialOption) (*Client, error) {
	err := validate.Struct(cfg)
	if err != nil {
		return nil, ErrInvalidConfig.Because(err)
	}

	validator, err := protovalidate.New()
	if err != nil {
		return nil, ErrInvalidConfig.Because(err)
	}

	creds := insecure.NewCredentials()
	if cfg.Secure {
		creds = credentials.NewTLS(new(tls.Config))
	}

	options = append(options,
		grpc.WithTransportCredentials(creds),
		grpc.WithChainUnaryInterceptor(
			validate.UnaryClientInterceptor(validator),
		),
	)

	conn, err := grpc.NewClient(cfg.Address, options...)
	if err != nil {
		return nil, ErrInvalidConfig.Because(err)
	}

	return &Client{
		NotesServiceClient: notesv1.NewNotesServiceClient(conn),
		UsersServiceClient: usersv1.NewUsersServiceClient(conn),
		conn:               conn,
		config:             cfg,
	}, nil
}

func MustNew(cfg Config, options ...grpc.DialOption) *Client {
	cli, err := New(cfg, options...)

	return ex.Critical(cli, err)
}

func (c *Client) Authenticate(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, authKey, token)
}

func (c *Client) Close() {
	err := c.conn.Close()

	ex.Skip(err)
}
