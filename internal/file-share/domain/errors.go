package domain

import "errors"

var (
	ErrEmptyPassword     = errors.New("password is empty")
	ErrEmptyPasswordHash = errors.New("password hash is empty")
	ErrInvalidPassword   = errors.New("password is invalid")
	ErrInvalidRoomTTL    = errors.New("room lifespan must be positive")

	ErrTokenNotFound        = errors.New("token not found or already revoked")
	ErrEmptyToken           = errors.New("token is empty")
	ErrTokenLifespanTooLong = errors.New("token lifespan too long")

	ErrFileNotFound = errors.New("file not found")
	ErrInvalidFile  = errors.New("invalid file")

	ErrRoomLifespanTooLong = errors.New("room lifespan too long")
	ErrRoomNotFound        = errors.New("room not found")
)
