package email

import (
	"sync"

	"github.com/therenotomorrow/gotes/internal/domain"
)

const (
	ErrInvalidEmail domain.Error = "invalid email"
)

var (
	validator Validator
	once      = sync.Once{}
)

type Validator interface {
	Validate(raw string) (bool, error)
}

func SetValidator(v Validator) {
	once.Do(func() {
		validator = v
	})
}
