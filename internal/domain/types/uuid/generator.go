package uuid

import (
	"sync"

	"github.com/therenotomorrow/gotes/internal/domain"
)

const (
	ErrInvalidUUID domain.Error = "invalid uuid"
)

var (
	generator Generator
	once      = sync.Once{}
)

type Generator interface {
	Generate() (string, error)
	Validate(raw string) (bool, error)
}

func SetGenerator(g Generator) {
	once.Do(func() {
		generator = g
	})
}
