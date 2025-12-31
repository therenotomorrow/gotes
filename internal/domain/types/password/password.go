package password

type Password struct {
	value string
}

func New(val string) Password {
	password, err := Create(val)
	if err != nil {
		panic(err)
	}

	return password
}

func Conv(hash string) Password {
	return Password{value: hash}
}

func Create(plain string) (Password, error) {
	val, err := hasher.Hash(plain)
	if err != nil {
		return Password{value: ""}, err
	}

	return Password{value: val}, nil
}

func (p Password) Check(plain string) error {
	valid, err := hasher.Verify(plain, p.value)

	switch {
	case err != nil:
		return err
	case !valid:
		return ErrInvalidPassword
	}

	return nil
}

func (p Password) Value() string {
	return p.value
}
