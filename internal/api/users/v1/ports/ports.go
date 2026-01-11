package ports

import (
	"context"

	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/email"
)

type UsersRepository interface {
	SaveUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetUserByEmail(ctx context.Context, mail email.Email) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
}

type Store struct {
	Users UsersRepository
}

type StoreProvider interface {
	Provide(ctx context.Context) Store
}

type UnitOfWork interface {
	Do(ctx context.Context, unit func(store Store) error) error
}
