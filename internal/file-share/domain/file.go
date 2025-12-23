package domain

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FileRoomFile struct {
	ID        uuid.UUID
	Path      string
	Name      string
	Size      int64
	CreatedAt time.Time
}

func NewFileRoomFile(path, name string, size int64, now time.Time) (*FileRoomFile, error) {
	path = strings.TrimSpace(path)
	name = strings.TrimSpace(name)

	if path == "" || name == "" || size <= 0 {
		return nil, ErrInvalidFile
	}

	safeName := filepath.Base(name)

	return &FileRoomFile{
		ID:        uuid.New(),
		Path:      path,
		Name:      safeName,
		Size:      size,
		CreatedAt: now,
	}, nil
}
