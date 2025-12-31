package auth

import (
	"context"

	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/internal/domain/entities"
)

type secureCtx string

const (
	authKey             = "authorization"
	secureKey secureCtx = "secure"

	errMissingUserInContext ex.Error = "missing user in context"
)

type Securable interface {
	User(ctx context.Context) *entities.User
}

type Secure struct{}

func (s Secure) User(ctx context.Context) *entities.User {
	user, ok := ctx.Value(secureKey).(*entities.User)
	if !ok {
		ex.Panic(errMissingUserInContext)
	}

	return user
}
