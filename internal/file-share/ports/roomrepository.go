package ports

import (
	"context"
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/google/uuid"
)

type RoomRepository interface {
	Get(ctx context.Context, roomID uuid.UUID) (*domain.Room, bool, error)
	List(ctx context.Context) ([]*domain.Room, error)
	Create(ctx context.Context, room *domain.Room) error
	Delete(ctx context.Context, roomID uuid.UUID) ([]string, error)
	DeleteExpired(ctx context.Context, now time.Time) ([]domain.ExpiredCleanup, error)
	RemoveToken(ctx context.Context, roomID uuid.UUID, token string) (bool, error)
	AddToken(ctx context.Context, roomID uuid.UUID, token string) error
	AddFileByToken(ctx context.Context, roomID uuid.UUID, token string, file *domain.RoomFile) (bool, error)
	DeleteFileByToken(ctx context.Context, roomID, fileID uuid.UUID, token string) (string, bool, error)
}
