package config

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Mode string
	Port string

	UploadDir string
	PublicDir string

	DefaultRoomTTL time.Duration
	TokenTTL       time.Duration

	MaxFiles         int
	MaxRoomBytes     int64
	MaxRoomLifespan  time.Duration
	MaxTokenLifespan time.Duration
	CleanupInterval  time.Duration

	JWTSecret []byte
}

func Load() (Config, error) {
	cfg := Config{}

	cfg.Mode = getEnv("MODE", "dev")
	cfg.Port = getEnv("PORT", "8080")

	cfg.UploadDir = getEnv("UPLOAD_DIR", "./uploads")
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		return cfg, fmt.Errorf("cannot create upload dir: %w", err)
	}
	cfg.PublicDir = getEnv("PUBLIC_DIR", "./public")
	if err := os.MkdirAll(cfg.PublicDir, 0755); err != nil {
		return cfg, fmt.Errorf("cannot create public dir: %w", err)
	}

	var err error
	cfg.DefaultRoomTTL, err = parseDurationEnv("ROOM_TTL", "10m")
	if err != nil {
		return cfg, err
	}
	cfg.TokenTTL, err = parseDurationEnv("TOKEN_TTL", "10m")
	if err != nil {
		return cfg, err
	}
	cfg.MaxRoomLifespan, err = parseDurationEnv("MAX_ROOM_LIFESPAN", "60m")
	if err != nil {
		return cfg, err
	}
	cfg.MaxTokenLifespan, err = parseDurationEnv("MAX_TOKEN_LIFESPAN", "60m")
	if err != nil {
		return cfg, err
	}
	cfg.CleanupInterval, err = parseDurationEnv("CLEANUP_INTERVAL", "30s")
	if err != nil {
		return cfg, err
	}

	cfg.MaxFiles, err = parseIntEnv("MAX_FILES", 30)
	if err != nil {
		return cfg, err
	}
	if cfg.MaxFiles <= 0 {
		return cfg, fmt.Errorf("MAX_FILES must be positive")
	}
	maxRoomMB, err := parseIntEnv("MAX_ROOM_MEGABYTES", 50)
	if err != nil {
		return cfg, err
	}
	if maxRoomMB <= 0 {
		return cfg, fmt.Errorf("MAX_ROOM_MEGABYTES must be positive")
	}
	cfg.MaxRoomBytes = int64(maxRoomMB) * 1024 * 1024

	if cfg.DefaultRoomTTL > cfg.MaxRoomLifespan {
		return cfg, fmt.Errorf("ROOM_TTL (%s) cannot exceed MAX_ROOM_LIFESPAN (%s)", cfg.DefaultRoomTTL, cfg.MaxRoomLifespan)
	}
	if cfg.TokenTTL > cfg.MaxTokenLifespan {
		return cfg, fmt.Errorf("TOKEN_TTL (%s) cannot exceed MAX_TOKEN_LIFESPAN (%s)", cfg.TokenTTL, cfg.MaxTokenLifespan)
	}
	if cfg.CleanupInterval <= 0 {
		return cfg, fmt.Errorf("CLEANUP_INTERVAL must be positive")
	}

	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return cfg, fmt.Errorf("cannot generate jwt secret: %w", err)
	}
	cfg.JWTSecret = secret

	return cfg, nil
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func parseIntEnv(key string, def int) (int, error) {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return def, nil
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return n, nil
}

func parseDurationEnv(key, def string) (time.Duration, error) {
	val := getEnv(key, def)
	d, err := time.ParseDuration(val)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return d, nil
}
