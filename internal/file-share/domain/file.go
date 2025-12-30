package domain

import (
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type RoomFile struct {
	ID        uuid.UUID
	Path      string
	Name      string
	Size      int64
	CreatedAt time.Time
}

func NewRoomFile(path, name string, size int64, now time.Time) (*RoomFile, error) {
	if path == "" || name == "" || size <= 0 {
		return nil, ErrInvalidFile
	}

	safeName := filepath.Base(name)

	return &RoomFile{
		ID:        uuid.New(),
		Path:      path,
		Name:      safeName,
		Size:      size,
		CreatedAt: now,
	}, nil
}
