package server

import (
	"github.com/therenotomorrow/gotes/internal/domain/types/email"
	"github.com/therenotomorrow/gotes/internal/domain/types/password"
	"github.com/therenotomorrow/gotes/internal/domain/types/uuid"
	"github.com/therenotomorrow/gotes/internal/services/auth"
	"github.com/therenotomorrow/gotes/internal/storages/postgres"
)

type Dependencies struct {
	Database       postgres.Database
	Secure         auth.Secure
	Authenticator  auth.Authenticator
	PasswordHasher password.Hasher
	UUIDGenerator  uuid.Generator
	EmailValidator email.Validator
}

func (d *Dependencies) Close() {
	d.Database.Close()
}
