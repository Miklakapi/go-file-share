package sqliterepository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"github.com/google/uuid"
)

type SqliteRepo struct {
	db *sql.DB
}

var _ ports.RoomRepository = (*SqliteRepo)(nil)

func New(db *sql.DB) *SqliteRepo {
	return &SqliteRepo{
		db: db,
	}
}

func (r *SqliteRepo) Get(ctx context.Context, roomID uuid.UUID) (*domain.Room, bool, error) {
	panic("TODO")
}

func (r *SqliteRepo) List(ctx context.Context) ([]*domain.Room, error) {
	panic("TODO")
}

func (r *SqliteRepo) Create(ctx context.Context, room *domain.Room) error {
	panic("TODO")
}

func (r *SqliteRepo) Update(ctx context.Context, room *domain.Room) error {
	panic("TODO")
}

func (r *SqliteRepo) Delete(ctx context.Context, roomID uuid.UUID) ([]string, error) {
	panic("TODO")
}

func (r *SqliteRepo) DeleteExpired(ctx context.Context, now time.Time) ([]domain.ExpiredCleanup, error) {
	panic("TODO")
}

func (r *SqliteRepo) RemoveToken(ctx context.Context, roomID uuid.UUID, token string) (bool, error) {
	panic("TODO")
}

func (r *SqliteRepo) AddToken(ctx context.Context, roomID uuid.UUID, token string) error {
	panic("TODO")
}

func (r *SqliteRepo) AddFileByToken(ctx context.Context, roomID uuid.UUID, token string, file *domain.RoomFile) (bool, error) {
	panic("TODO")
}

func (r *SqliteRepo) DeleteFileByToken(ctx context.Context, roomID, fileID uuid.UUID, token string) (string, bool, error) {
	panic("TODO")
}
