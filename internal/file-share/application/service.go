package application

import (
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
)

type Service struct {
	rooms  ports.RoomRepository
	files  ports.FileStore
	hasher ports.PasswordHasher
	now    func() time.Time
}

func NewService(rooms ports.RoomRepository, files ports.FileStore, hasher ports.PasswordHasher) *Service {
	return &Service{
		rooms:  rooms,
		files:  files,
		hasher: hasher,
		now:    time.Now,
	}
}
