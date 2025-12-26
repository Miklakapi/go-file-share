package ports

import (
	"context"
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/google/uuid"
)

type RoomRepository interface {
	Get(ctx context.Context, roomID uuid.UUID) (*domain.FileRoom, bool, error)
	GetByToken(ctx context.Context, roomID uuid.UUID, token string) (*domain.FileRoom, bool, error)
	ListSnapshots(ctx context.Context) ([]domain.RoomSnapshot, error)
	Create(ctx context.Context, room *domain.FileRoom) error
	Update(ctx context.Context, room *domain.FileRoom) error
	Delete(ctx context.Context, roomID uuid.UUID) ([]string, error)
	DeleteExpired(ctx context.Context, now time.Time) ([]domain.ExpiredCleanup, error)
	RemoveToken(ctx context.Context, roomID uuid.UUID, token string) (bool, error)
	GetPasswordHash(ctx context.Context, roomID uuid.UUID) (hash string, ok bool, err error)
	AddToken(ctx context.Context, roomID uuid.UUID, token string) error
}
