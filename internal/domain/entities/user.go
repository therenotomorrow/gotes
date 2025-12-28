package entities

import (
	"github.com/therenotomorrow/gotes/internal/domain/types/email"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
	"github.com/therenotomorrow/gotes/internal/domain/types/uuid"
)

const (
	ErrEmptyName Error = "empty name"
)

type User struct {
	Name  string
	Email email.Email
	Token uuid.UUID
	ID    id.ID
}

func NewUser(name string, mail email.Email) (*User, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	return &User{
		Name:  name,
		Email: mail,
		ID:    id.ID{},
		Token: uuid.New(),
	}, nil
}
