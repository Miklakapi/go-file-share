package dto

import (
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/google/uuid"
)

type CreateRoomRequest struct {
	Password string `json:"password" form:"password" binding:"required"`
	Lifespan int    `json:"lifespan" form:"lifespan"`
}

type AuthRoomRequest struct {
	CreateRoomRequest
}

type Room struct {
	ID        uuid.UUID `json:"id"`
	ExpiresAt time.Time `json:"expiresAt"`
	Files     int       `json:"files"`
	Tokens    int       `json:"tokens"`
}

func NewRoom(s domain.RoomSnapshot) Room {
	return Room{
		ID:        s.ID,
		ExpiresAt: s.ExpiresAt,
		Files:     s.Files,
		Tokens:    s.Tokens,
	}
}

type RoomFile struct {
	ID        uuid.UUID `json:"id"`
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewFileRoomFile(s domain.FileRoomFile) RoomFile {
	return RoomFile{
		ID:        s.ID,
		Path:      s.Path,
		Name:      s.Name,
		Size:      s.Size,
		CreatedAt: s.CreatedAt,
	}
}
