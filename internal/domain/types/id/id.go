package id

import (
	"encoding/json"

	"github.com/therenotomorrow/gotes/internal/domain"
)

const (
	ErrInvalidID domain.Error = "invalid id"
	ErrMarshal   domain.Error = "cannot marshal id"
	ErrUnmarshal domain.Error = "cannot unmarshal id"
)

var zeroID ID

type ID struct {
	value int64
}

func (id *ID) MarshalJSON() ([]byte, error) {
	val, err := json.Marshal(id.value)
	if err != nil {
		return nil, ErrMarshal
	}

	return val, nil
}

func (id *ID) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &id.value)
	if err != nil {
		return ErrUnmarshal
	}

	return nil
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

func (id *ID) Value() int64 {
	return id.value
}

func (id *ID) ValuePtr() *int64 {
	return &id.value
}
