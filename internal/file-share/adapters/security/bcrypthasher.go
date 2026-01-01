package security

import (
	"context"

	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	Cost int
}

var _ ports.PasswordHasher = (*BcryptHasher)(nil)

func (h BcryptHasher) Hash(ctx context.Context, plain string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	cost := h.Cost
	if cost == 0 {
		cost = 12
	}
	b, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	return string(b), err
}

func (h BcryptHasher) Verify(ctx context.Context, plain, hash string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil, nil
}
