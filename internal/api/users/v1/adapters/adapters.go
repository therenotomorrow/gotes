package adapters

import (
	"context"
	"database/sql"
	"errors"

	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/api/users/v1/ports"
	"github.com/therenotomorrow/gotes/internal/domain"
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
		return nil, domain.ErrUserNotFound
	case err != nil:
		return nil, ex.Unexpected(err)
	}

	return user.ToEntity(), nil
}

func (u *UsersRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	err := u.commands.UpdateUser(ctx, commands.NewUpdateUserParams(user))

	return ex.Unexpected(err)
}

type Store struct {
	users *UsersRepository
}

func NewStore(dbtx postgres.DBTX) *Store {
	return &Store{users: NewUsersRepository(dbtx)}
}

func (s *Store) Users() ports.UsersRepository {
	return s.users
}

type UnitOfWork struct {
	database postgres.Database
}

func NewUnitOfWork(database postgres.Database) *UnitOfWork {
	return &UnitOfWork{database: database}
}

func (u *UnitOfWork) Do(ctx context.Context, work func(store ports.Store) error) error {
	return u.database.Tx(ctx, func(ctx context.Context) error {
		cqrs := u.database.Conn(ctx)
		store := NewStore(cqrs)

		return work(store)
	})
}
