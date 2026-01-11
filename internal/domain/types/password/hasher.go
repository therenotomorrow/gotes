package password

import (
	"sync"

	"github.com/therenotomorrow/gotes/internal/domain"
)

const (
	ErrInvalidPassword domain.Error = "invalid password"
)

var (
	hasher Hasher
	once   = sync.Once{}
)

type Hasher interface {
	Hash(plain string) (string, error)
	Verify(plain, encoded string) (bool, error)
}

func SetHasher(h Hasher) {
	once.Do(func() {
		hasher = h
	})
}
