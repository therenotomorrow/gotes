package server

import (
	"github.com/therenotomorrow/gotes/internal/domain/types/email"
	"github.com/therenotomorrow/gotes/internal/domain/types/password"
	"github.com/therenotomorrow/gotes/internal/domain/types/uuid"
	"github.com/therenotomorrow/gotes/internal/services/secure"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
)

type Dependencies struct {
	Database       postgres.Database
	Authenticator  secure.Authenticator
	PasswordHasher password.Hasher
	UUIDGenerator  uuid.Generator
	EmailValidator email.Validator
}

func (d *Dependencies) Close() {
	d.Database.Close()
}
