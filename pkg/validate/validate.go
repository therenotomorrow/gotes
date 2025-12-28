package validate

import (
	"errors"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/therenotomorrow/ex"
)

const (
	ErrRegisterValidator ex.Error = "register validator"
)

func RegisterPostgresDSN(v *validator.Validate) error {
	err := v.RegisterValidation("postgres_dsn", func(fl validator.FieldLevel) bool {
		dsn := fl.Field().String()
		_, err := pgx.ParseConfig(dsn)

		return err == nil
	})
	if err != nil {
		err = ErrRegisterValidator.Because(err)
	}

	return err
}

var singleton = sync.OnceValue(validate)

var (
	Var    = singleton().Var
	Struct = singleton().Struct
)

func validate() *validator.Validate {
	val := validator.New()
	err := errors.Join(RegisterPostgresDSN(val))

	ex.Panic(err)

	return val
}
