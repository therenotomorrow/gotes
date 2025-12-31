package validate

import (
	"github.com/therenotomorrow/ex"
)

type Email struct{}

func NewEmail() *Email {
	return &Email{}
}

func (e *Email) Validate(raw string) (bool, error) {
	err := Var(raw, "email")

	return err == nil, ex.Unexpected(err)
}
