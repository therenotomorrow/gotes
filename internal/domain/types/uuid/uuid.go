package uuid

type UUID struct {
	value string
}

func New() UUID {
	uuid, err := Create()
	if err != nil {
		panic(err)
	}

	return uuid
}

func Create() (UUID, error) {
	val, err := generator.Generate()
	if err != nil {
		return UUID{value: ""}, err
	}

	return UUID{value: val}, nil
}

func Conv(uuid string) UUID {
	val, err := Parse(uuid)
	if err != nil {
		panic(err)
	}

	return val
}

func Parse(raw string) (UUID, error) {
	valid, err := generator.Validate(raw)

	switch {
	case err != nil:
		return UUID{value: ""}, err
	case !valid:
		return UUID{value: ""}, ErrInvalidUUID
	}

	return UUID{value: raw}, nil
}

func (u UUID) Value() string {
	return u.value
}
