package ports

import (
	"context"
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/google/uuid"
)

type RoomRepository interface {
	Get(ctx context.Context, roomID uuid.UUID) (*domain.FileRoom, bool, error)
	ListSnapshots(ctx context.Context) ([]domain.RoomSnapshot, error)
	Create(ctx context.Context, room *domain.FileRoom) error
	Update(ctx context.Context, room *domain.FileRoom) error
	Delete(ctx context.Context, roomID uuid.UUID) error
	DeleteExpired(ctx context.Context, now time.Time) ([]uuid.UUID, error)
}
