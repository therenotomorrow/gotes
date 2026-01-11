package queries

import (
	"github.com/therenotomorrow/gotes/internal/domain/entities"
	"github.com/therenotomorrow/gotes/internal/domain/types/email"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
	"github.com/therenotomorrow/gotes/internal/domain/types/password"
	"github.com/therenotomorrow/gotes/internal/domain/types/uuid"
)

func (u *User) ToEntity() *entities.User {
	return &entities.User{
		Name:      u.Name,
		Email:     email.New(u.Email),
		Password:  password.Conv(u.Password),
		Token:     uuid.Conv(u.Token.String()),
		ID:        id.New(u.ID),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
