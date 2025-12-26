package domain

import "errors"

var (
	ErrEmptyPassword     = errors.New("password is empty")
	ErrEmptyPasswordHash = errors.New("password hash is empty")
	ErrEmptyCreatorToken = errors.New("creator token is empty")
	ErrInvalidRoomTTL    = errors.New("room lifespan must be positive")

	ErrTokenNotFound = errors.New("token not found")
	ErrEmptyToken    = errors.New("token is empty")

	ErrFileNotFound = errors.New("file not found")
	ErrInvalidFile  = errors.New("invalid file")

	ErrRoomLifespanTooLong = errors.New("room lifespan too long")
	ErrRoomNotFound        = errors.New("room not found")
)
