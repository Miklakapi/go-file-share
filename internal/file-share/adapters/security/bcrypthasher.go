package security

import (
	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	Cost int
}

var _ ports.PasswordHasher = (*BcryptHasher)(nil)

func (h BcryptHasher) Hash(plain string) (string, error) {
	cost := h.Cost
	if cost == 0 {
		cost = 14
	}
	b, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	return string(b), err
}

func (h BcryptHasher) Verify(plain, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
