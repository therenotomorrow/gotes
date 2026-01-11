package ports_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/adapters/mocks"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/adapters/postgres"
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/ports"
)

func TestNotesRepository(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*ports.NotesRepository)(nil), new(postgres.NotesRepository))
	assert.Implements(t, (*ports.NotesRepository)(nil), new(mocks.MockNotesRepository))
}

func TestStoreProvider(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*ports.StoreProvider)(nil), new(postgres.StoreProvider))
	assert.Implements(t, (*ports.StoreProvider)(nil), new(mocks.MockStoreProvider))
}

func TestUnitOfWork(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*ports.UnitOfWork)(nil), new(postgres.UnitOfWork))
	assert.Implements(t, (*ports.UnitOfWork)(nil), new(mocks.MockUnitOfWork))
}
