package domain

import (
	"time"

	"github.com/google/uuid"
)

type RoomSnapshot struct {
	ID        uuid.UUID
	ExpiresAt time.Time
	Files     int
	Tokens    int
}
