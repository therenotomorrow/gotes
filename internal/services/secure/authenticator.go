package secure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	domuuid "github.com/therenotomorrow/gotes/internal/domain/types/uuid"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
	queries "github.com/therenotomorrow/gotes/internal/storages/postgres/queries/users"
)

const (
	ErrUnauthorized     ex.Error = "unauthorized"
	errHideUserNotFound ex.Error = "wrong token or user not exists"
)

type TokenAuthenticator struct {
	db postgres.Database
}

func NewTokenAuthenticator(db postgres.Database) *TokenAuthenticator {
	return &TokenAuthenticator{db: db}
}

func (a *TokenAuthenticator) Authenticate(ctx context.Context, token string) (*entities.User, error) {
	querier := queries.New(a.db.Conn(ctx))

	tkn, err := domuuid.Parse(token)
	if err != nil {
		return nil, ErrUnauthorized.Because(err)
	}

	user, err := querier.SelectUserByToken(ctx, uuid.MustParse(tkn.Value()))

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrUnauthorized.Because(errHideUserNotFound)
	case err != nil:
		return nil, ErrUnauthorized.Because(err)
	}

	return user.ToEntity(), nil
}
