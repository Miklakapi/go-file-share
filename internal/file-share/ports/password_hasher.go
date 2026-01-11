package ports

import "context"

type PasswordHasher interface {
	Hash(ctx context.Context, plain string) (string, error)
	Verify(ctx context.Context, plain, hash string) (bool, error)
}
