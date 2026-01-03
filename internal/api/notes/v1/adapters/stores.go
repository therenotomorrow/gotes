package adapters

import (
	"github.com/therenotomorrow/gotes/internal/api/notes/v1/ports"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
)

type Store struct {
	notes *NotesRepository
}

func NewStore(cqrs *postgres.CQRS) *Store {
	return &Store{notes: NewNotesRepository(cqrs)}
}

func (s *Store) Notes() ports.NotesRepository {
	return s.notes
}
