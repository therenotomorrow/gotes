package v1

import (
	"context"
	"errors"

	"github.com/therenotomorrow/gotes/internal/api/users/v1/ports"
	"github.com/therenotomorrow/gotes/internal/domain"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/email"
	"github.com/therenotomorrow/gotes/internal/domain/types/password"
)

type UseCases struct {
	uow ports.UnitOfWork
}

func New(uow ports.UnitOfWork) *UseCases {
	return &UseCases{uow: uow}
}

type RegisterUserInput struct {
	Name     string
	Email    string
	Password string
}

func (use *UseCases) RegisterUser(ctx context.Context, input *RegisterUserInput) (*entities.User, error) {
	user, err := entities.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	err = use.uow.Do(ctx, func(store ports.Store) error {
		repo := store.Users()
		saved, err := repo.GetUserByEmail(ctx, user.Email)

		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			user, err = repo.SaveUser(ctx, user)
			if err != nil {
				return err
			}
		case err != nil:
			return err
		case user.Same(saved):
			user = saved
		}

		return use.checkPassword(user, input.Password)
	})

	return user, err
}

type RefreshTokenInput struct {
	Email    string
	Password string
}

func (use *UseCases) RefreshToken(ctx context.Context, input *RefreshTokenInput) (*entities.User, error) {
	mail, err := email.Parse(input.Email)
	if err != nil {
		return nil, err
	}

	var user *entities.User

	err = use.uow.Do(ctx, func(store ports.Store) error {
		repo := store.Users()
		user, err = repo.GetUserByEmail(ctx, mail)

		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return domain.ErrPermissionDenied
		case err != nil:
			return err
		}

		err = use.checkPassword(user, input.Password)
		if err != nil {
			return err
		}

		user.RefreshToken()

		return repo.UpdateUser(ctx, user)
	})

	return user, err
}

func (use *UseCases) checkPassword(user *entities.User, pass string) error {
	err := user.Password.Check(pass)

	if errors.Is(err, password.ErrInvalidPassword) {
		return domain.ErrPermissionDenied
	}

	return err
}
