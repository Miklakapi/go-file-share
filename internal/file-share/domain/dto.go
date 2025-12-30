package domain

import (
	"github.com/google/uuid"
)

type ExpiredCleanup struct {
	RoomID uuid.UUID
	Paths  []string
}
