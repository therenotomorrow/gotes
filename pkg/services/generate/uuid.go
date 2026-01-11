package generate

import (
	"github.com/google/uuid"
	"github.com/therenotomorrow/ex"
)

type UUID struct{}

func NewUUID() *UUID {
	return &UUID{}
}

func (g *UUID) Generate() (string, error) {
	val, err := uuid.NewRandom()
	if err != nil {
		return "", ex.Unexpected(err)
	}

	return val.String(), nil
}

func (g *UUID) Validate(raw string) (bool, error) {
	_, err := uuid.Parse(raw)

	return err == nil, ex.Unexpected(err)
}
