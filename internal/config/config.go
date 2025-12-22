package config

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port           string
	UploadDir      string
	PublicDir      string
	DefaultRoomTTL time.Duration
	TokenTTL       time.Duration
	Mode           string
	JWTSecret      []byte
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
	roomTTL, err := time.ParseDuration(getEnv("ROOM_TTL", "15m"))
	if err != nil {
		return cfg, fmt.Errorf("invalid ROOM_TTL: %w", err)
	}
	cfg.DefaultRoomTTL = roomTTL
	tokenTTL, err := time.ParseDuration(getEnv("TOKEN_TTL", "15m"))
	if err != nil {
		return cfg, fmt.Errorf("invalid TOKEN_TTL: %w", err)
	}
	cfg.TokenTTL = tokenTTL
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
