package ports_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/therenotomorrow/gotes/internal/api/users/v1/adapters/mocks"
	"github.com/therenotomorrow/gotes/internal/api/users/v1/adapters/postgres"
	"github.com/therenotomorrow/gotes/internal/api/users/v1/ports"
)

func TestNotesRepository(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*ports.UsersRepository)(nil), new(postgres.UsersRepository))
	assert.Implements(t, (*ports.UsersRepository)(nil), new(mocks.MockUsersRepository))
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
