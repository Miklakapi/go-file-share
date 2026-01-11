package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TokenService interface {
	Issue(ctx context.Context, roomID uuid.UUID, ttl time.Duration) (token string, expiresAt time.Time, err error)
	Validate(ctx context.Context, token string) error
	ValidateWithRoom(ctx context.Context, roomID uuid.UUID, token string) error
}
