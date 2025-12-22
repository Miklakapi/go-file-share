package files

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

type roomSnapResp struct {
	room RoomSnapshot
	ok   bool
	err  error
}

type roomsSnapResp struct {
	rooms []RoomSnapshot
}

type errResp struct {
	err error
}
