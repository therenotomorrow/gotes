package entities

import (
	"time"

	"github.com/therenotomorrow/gotes/internal/domain"
	"github.com/therenotomorrow/gotes/internal/domain/types/email"
	"github.com/therenotomorrow/gotes/internal/domain/types/id"
	"github.com/therenotomorrow/gotes/internal/domain/types/password"
	"github.com/therenotomorrow/gotes/internal/domain/types/uuid"
)

const (
	ErrEmptyEmail    domain.Error = "empty email"
	ErrEmptyPassword domain.Error = "empty password"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     email.Email
	Password  password.Password
	Token     uuid.UUID
	ID        id.ID
}

func NewUser(name, addr, pass string) (*User, error) {
	if addr == "" {
		return nil, ErrEmptyEmail
	}

	if pass == "" {
		return nil, ErrEmptyPassword
	}

	mail, err := email.Parse(addr)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &User{
		Name:      name,
		Email:     mail,
		Password:  password.New(pass),
		ID:        id.ID{},
		Token:     uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) Same(o *User) bool {
	return u.Email.Equals(o.Email)
}

func (u *User) RefreshToken() {
	u.Token = uuid.New()
	u.UpdatedAt = time.Now()
}
