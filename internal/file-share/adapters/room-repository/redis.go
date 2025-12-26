package roomrepository

import (
	"context"
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"github.com/google/uuid"
)

type RedisRepo struct {
	rooms map[uuid.UUID]*domain.FileRoom
}

var _ ports.RoomRepository = (*RedisRepo)(nil)

func NewRedisRepo() *RedisRepo {
	return &RedisRepo{
		rooms: make(map[uuid.UUID]*domain.FileRoom),
	}
}

func (r *RedisRepo) Get(ctx context.Context, roomID uuid.UUID) (*domain.FileRoom, bool, error) {
	panic("TODO")
}

func (r *RedisRepo) ListSnapshots(ctx context.Context) ([]domain.RoomSnapshot, error) {
	panic("TODO")
}

func (r *RedisRepo) Create(ctx context.Context, room *domain.FileRoom) error {
	panic("TODO")
}

func (r *RedisRepo) Update(ctx context.Context, room *domain.FileRoom) error {
	panic("TODO")
}

func (r *RedisRepo) Delete(ctx context.Context, roomID uuid.UUID) error {
	panic("TODO")
}

func (r *RedisRepo) DeleteExpired(ctx context.Context, now time.Time) ([]uuid.UUID, error) {
	panic("TODO")
}
