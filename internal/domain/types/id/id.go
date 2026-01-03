package id

import (
	"strconv"

	"github.com/therenotomorrow/ex"
)

const (
	ErrInvalidID  ex.Error = "invalid id"
	ErrNegativeID ex.Error = "negative id"
)

var zeroID ID

type ID struct {
	value int64
}

func New(val int64) ID {
	return ex.Critical(Conv(val))
}

func Conv(val int64) (ID, error) {
	if val < 1 {
		return zeroID, ErrNegativeID
	}

	return ID{value: val}, nil
}

func Parse(raw string) (ID, error) {
	val, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return zeroID, ErrInvalidID.Because(err)
	}

	if val < 1 {
		return zeroID, ErrNegativeID
	}

	return ID{value: val}, nil
}

func (id ID) Value() int64 {
	return id.value
}

func (id ID) String() string {
	return strconv.FormatInt(id.value, 10)
}
