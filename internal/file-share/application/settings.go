package application

import (
	"time"
)

type Settings struct {
	DefaultRoomTTL   time.Duration
	TokenTTL         time.Duration
	MaxFiles         int
	MaxRoomBytes     int64
	MaxRoomLifespan  time.Duration
	MaxTokenLifespan time.Duration
	CleanupInterval  time.Duration
}

func NewSettings(
	defaultRoomTTL time.Duration,
	tokenTTL time.Duration,
	maxFiles int,
	maxRoomBytes int64,
	maxRoomLifespan time.Duration,
	maxTokenLifespan time.Duration,
	cleanupInterval time.Duration,
) Settings {
	return Settings{
		DefaultRoomTTL:   defaultRoomTTL,
		TokenTTL:         tokenTTL,
		MaxFiles:         maxFiles,
		MaxRoomBytes:     maxRoomBytes,
		MaxRoomLifespan:  maxRoomLifespan,
		MaxTokenLifespan: maxTokenLifespan,
		CleanupInterval:  cleanupInterval,
	}
}
