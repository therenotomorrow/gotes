package secure

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

func User(ctx context.Context) (*entities.User, error) {
	user, ok := ctx.Value(secureKey).(*entities.User)
	if !ok {
		return nil, ErrUnauthorized.Because(errMissingUserInContext)
	}

	return user, nil
}

func NewUserContext(ctx context.Context, user *entities.User) context.Context {
	return context.WithValue(ctx, secureKey, user)
}
