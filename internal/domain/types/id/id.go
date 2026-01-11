package id

import (
	"github.com/therenotomorrow/gotes/internal/domain"
)

const (
	ErrInvalidID domain.Error = "invalid id"
)

var zeroID ID

type ID struct {
	value int64
}

func New(val int64) ID {
	id, err := Conv(val)
	if err != nil {
		panic(err)
	}

	return id
}

func Conv(val int64) (ID, error) {
	if val < 1 {
		return zeroID, ErrInvalidID
	}

	return ID{value: val}, nil
}

func (id ID) Value() int64 {
	return id.value
}

func (id ID) ValuePtr() *int64 {
	return &id.value
}
