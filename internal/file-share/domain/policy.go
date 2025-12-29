package domain

import (
	"time"
)

type Policy struct {
	DefaultRoomTTL   time.Duration
	DefaultTokenTTL  time.Duration
	MaxFiles         int
	MaxRoomBytes     int64
	MaxRoomLifespan  time.Duration
	MaxTokenLifespan time.Duration
	UploadDir        string
}

func NewPolicy(
	defaultRoomTTL time.Duration,
	defaultTokenTTL time.Duration,
	maxFiles int,
	maxRoomBytes int64,
	maxRoomLifespan time.Duration,
	maxTokenLifespan time.Duration,
	uploadDir string,
) Policy {
	return Policy{
		DefaultRoomTTL:   defaultRoomTTL,
		DefaultTokenTTL:  defaultTokenTTL,
		MaxFiles:         maxFiles,
		MaxRoomBytes:     maxRoomBytes,
		MaxRoomLifespan:  maxRoomLifespan,
		MaxTokenLifespan: maxTokenLifespan,
		UploadDir:        uploadDir,
	}
}
