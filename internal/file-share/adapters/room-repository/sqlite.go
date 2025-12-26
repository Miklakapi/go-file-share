package roomrepository

import (
	"context"
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"github.com/google/uuid"
)

type SqliteRepo struct {
	rooms map[uuid.UUID]*domain.FileRoom
}

var _ ports.RoomRepository = (*SqliteRepo)(nil)

func NewSqliteRepo() *SqliteRepo {
	return &SqliteRepo{
		rooms: make(map[uuid.UUID]*domain.FileRoom),
	}
}

func (r *SqliteRepo) Get(ctx context.Context, roomID uuid.UUID) (*domain.FileRoom, bool, error) {
	panic("TODO")
}

func (r *SqliteRepo) ListSnapshots(ctx context.Context) ([]domain.RoomSnapshot, error) {
	panic("TODO")
}

func (r *SqliteRepo) Create(ctx context.Context, room *domain.FileRoom) error {
	panic("TODO")
}

func (r *SqliteRepo) Update(ctx context.Context, room *domain.FileRoom) error {
	panic("TODO")
}

func (r *SqliteRepo) Delete(ctx context.Context, roomID uuid.UUID) error {
	panic("TODO")
}

func (r *SqliteRepo) DeleteExpired(ctx context.Context, now time.Time) ([]uuid.UUID, error) {
	panic("TODO")
}
