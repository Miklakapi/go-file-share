package redisrepository

import (
	"context"
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	db *redis.Client
}

func New(db *redis.Client) *RedisRepo {
	return &RedisRepo{
		db: db,
	}
}

func (r *RedisRepo) Get(ctx context.Context, roomID uuid.UUID) (*domain.Room, bool, error) {
	panic("TODO")
}

func (r *RedisRepo) List(ctx context.Context) ([]*domain.Room, error) {
	panic("TODO")
}

func (r *RedisRepo) Create(ctx context.Context, room *domain.Room) error {
	panic("TODO")
}

func (r *RedisRepo) Delete(ctx context.Context, roomID uuid.UUID) ([]string, error) {
	panic("TODO")
}

func (r *RedisRepo) DeleteExpired(ctx context.Context, now time.Time) ([]domain.ExpiredCleanup, error) {
	panic("TODO")
}

func (r *RedisRepo) RemoveToken(ctx context.Context, roomID uuid.UUID, token string) (bool, error) {
	panic("TODO")
}

func (r *RedisRepo) AddToken(ctx context.Context, roomID uuid.UUID, token string) error {
	panic("TODO")
}

func (r *RedisRepo) AddFileByToken(ctx context.Context, roomID uuid.UUID, token string, file *domain.RoomFile) (bool, error) {
	panic("TODO")
}

func (r *RedisRepo) DeleteFileByToken(ctx context.Context, roomID, fileID uuid.UUID, token string) (string, bool, error) {
	panic("TODO")
}
