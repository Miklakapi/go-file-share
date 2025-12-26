package application

import (
	"time"
)

type Settings struct {
	DefaultRoomTTL   time.Duration
	DefaultTokenTTL  time.Duration
	MaxFiles         int
	MaxRoomBytes     int64
	MaxRoomLifespan  time.Duration
	MaxTokenLifespan time.Duration
	CleanupInterval  time.Duration
}

func NewSettings(
	defaultRoomTTL time.Duration,
	defaultTokenTTL time.Duration,
	maxFiles int,
	maxRoomBytes int64,
	maxRoomLifespan time.Duration,
	maxTokenLifespan time.Duration,
	cleanupInterval time.Duration,
) Settings {
	return Settings{
		DefaultRoomTTL:   defaultRoomTTL,
		DefaultTokenTTL:  defaultTokenTTL,
		MaxFiles:         maxFiles,
		MaxRoomBytes:     maxRoomBytes,
		MaxRoomLifespan:  maxRoomLifespan,
		MaxTokenLifespan: maxTokenLifespan,
		CleanupInterval:  cleanupInterval,
	}
}
