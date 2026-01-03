package uuid

import (
	"github.com/google/uuid"
	"github.com/therenotomorrow/ex"
)

const (
	ErrInvalidUUID ex.Error = "invalid uuid"
)

var zeroUUID UUID

type UUID struct {
	value string
}

func New() UUID {
	val := ex.Critical(uuid.NewRandom())

	return UUID{value: val.String()}
}

func Parse(raw string) (UUID, error) {
	val, err := uuid.Parse(raw)
	if err != nil {
		return zeroUUID, ErrInvalidUUID.Because(err)
	}

	return UUID{value: val.String()}, nil
}

func (u UUID) Value() string {
	return u.value
}

func (u UUID) String() string {
	return u.value
}
