package email

import (
	"github.com/therenotomorrow/ex"
	"github.com/therenotomorrow/gotes/pkg/validate"
)

const (
	ErrInvalidEmail ex.Error = "invalid email"
)

var zeroEmail Email

type Email struct {
	value string
}

func New(val string) Email {
	return ex.Critical(Parse(val))
}

func Parse(raw string) (Email, error) {
	err := validate.Var(raw, "email")
	if err != nil {
		return zeroEmail, ErrInvalidEmail.Because(err)
	}

	return Email{value: raw}, nil
}

func (e Email) Value() string {
	return e.value
}

func (e Email) String() string {
	return e.value
}
