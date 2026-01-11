package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/api/users/v1/usecases"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/email"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
	commands "github.com/therenotomorrow/gotes/internal/storages/postgres/commands/users"
	queries "github.com/therenotomorrow/gotes/internal/storages/postgres/queries/users"
)

type UsersRepository struct {
	commands commands.Querier
	queries  queries.Querier
}

func NewUsersRepository(dbtx postgres.DBTX) *UsersRepository {
	return &UsersRepository{commands: commands.New(dbtx), queries: queries.New(dbtx)}
}

func (u *UsersRepository) SaveUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	ident, err := u.commands.InsertUser(ctx, commands.NewInsertUserParams(user))
	if err != nil {
		return nil, ex.Unexpected(err)
	}

	user.ID = id.New(ident)

	return user, nil
}

func (u *UsersRepository) GetUserByEmail(ctx context.Context, mail email.Email) (*entities.User, error) {
	user, err := u.queries.SelectUserByEmail(ctx, mail.Value())

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, usecases.ErrUserNotFound
	case err != nil:
		return nil, ex.Unexpected(err)
	}

	return user.ToEntity(), nil
}

func (u *UsersRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	err := u.commands.UpdateUser(ctx, commands.NewUpdateUserParams(user))

	return ex.Unexpected(err)
}
