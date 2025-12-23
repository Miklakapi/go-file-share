package ports

import (
	"context"
	"time"
)

type TokenService interface {
	Issue(ctx context.Context, ttl time.Duration) (token string, expiresAt time.Time, err error)
	Validate(ctx context.Context, token string) error
}
