package usecases

import (
	"context"
	"errors"

	"github.com/therenotomorrow/gotes/internal/api/users/v1/ports"
	"github.com/therenotomorrow/gotes/internal/domain"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/email"
	"github.com/therenotomorrow/gotes/internal/domain/types/password"
)

const (
	ErrUserNotFound     domain.Error = "user not found"
	ErrPermissionDenied domain.Error = "permission denied"
)

type UseCases struct {
	uow ports.UnitOfWork
}

func NewCases(uow ports.UnitOfWork) *UseCases {
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
		saved, err := store.Users.GetUserByEmail(ctx, user.Email)

		switch {
		case errors.Is(err, ErrUserNotFound):
			user, err = store.Users.SaveUser(ctx, user)
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
		user, err = store.Users.GetUserByEmail(ctx, mail)

		switch {
		case errors.Is(err, ErrUserNotFound):
			return ErrPermissionDenied
		case err != nil:
			return err
		}

		err = use.checkPassword(user, input.Password)
		if err != nil {
			return err
		}

		user.RefreshToken()

		return store.Users.UpdateUser(ctx, user)
	})

	return user, err
}

func (use *UseCases) checkPassword(user *entities.User, pass string) error {
	err := user.Password.Check(pass)

	if errors.Is(err, password.ErrInvalidPassword) {
		return ErrPermissionDenied
	}

	return err
}
