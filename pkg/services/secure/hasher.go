package secure

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"

	"github.com/therenotomorrow/ex"
	"golang.org/x/crypto/argon2"
)

const (
	hashTime    = 3
	hashMemory  = 64 * 1024
	hashThreads = 4
	hashKeyLen  = 32
	saltLen     = 16

	ErrEncodePassword ex.Error = "encode password error"
	ErrDecodePassword ex.Error = "decode password error"
)

type PasswordHasher struct {
	encoder *base64.Encoding
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{encoder: base64.RawStdEncoding}
}

func (p *PasswordHasher) Hash(plain string) (string, error) {
	salt := make([]byte, saltLen)

	_, err := rand.Read(salt)
	if err != nil {
		return "", ErrEncodePassword.Because(err)
	}

	hash := p.hash(plain, salt)

	return p.encoder.EncodeToString(append(salt, hash...)), nil
}

func (p *PasswordHasher) Verify(plain, encoded string) (bool, error) {
	data, err := p.encoder.DecodeString(encoded)
	if err != nil {
		return false, ErrDecodePassword.Because(err)
	}

	if len(data) != saltLen+hashKeyLen {
		return false, ErrDecodePassword.Reason("corrupted data")
	}

	hash := p.hash(plain, data[:saltLen])

	return subtle.ConstantTimeCompare(data[saltLen:], hash) == 1, nil
}

func (p *PasswordHasher) hash(plain string, salt []byte) []byte {
	return argon2.IDKey([]byte(plain), salt, hashTime, hashMemory, hashThreads, hashKeyLen)
}
