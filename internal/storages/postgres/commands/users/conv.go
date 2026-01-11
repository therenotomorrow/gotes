package commands

import (
	"github.com/google/uuid"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
)

func NewInsertUserParams(user *entities.User) *InsertUserParams {
	return &InsertUserParams{
		Name:      user.Name,
		Email:     user.Email.Value(),
		Password:  user.Password.Value(),
		Token:     uuid.MustParse(user.Token.Value()),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewUpdateUserParams(user *entities.User) *UpdateUserParams {
	return &UpdateUserParams{
		Name:      user.Name,
		Email:     user.Email.Value(),
		Password:  user.Password.Value(),
		Token:     uuid.MustParse(user.Token.Value()),
		UpdatedAt: user.UpdatedAt,
		ID:        user.ID.Value(),
	}
}
