package email

type Email struct {
	value string
}

func New(val string) Email {
	email, err := Parse(val)
	if err != nil {
		panic(err)
	}

	return email
}

func Parse(raw string) (Email, error) {
	valid, err := validator.Validate(raw)

	switch {
	case err != nil:
		return Email{value: ""}, err
	case !valid:
		return Email{value: ""}, ErrInvalidEmail
	}

	return Email{value: raw}, nil
}

func (e Email) Equals(o Email) bool {
	return e.value == o.value
}

func (e Email) Value() string {
	return e.value
}
