package domain

import "errors"

var (
	ErrEmptyPasswordHash = errors.New("password hash is empty")
	ErrEmptyCreatorToken = errors.New("creator token is empty")
	ErrInvalidRoomTTL    = errors.New("room lifespan must be positive")
	ErrEmptyUploadDir    = errors.New("upload dir is empty")
	ErrEmptyFilename     = errors.New("filename is empty")
	ErrNilReader         = errors.New("file reader is nil")
	ErrFileNotFound      = errors.New("file not found")
	ErrEmptyToken        = errors.New("token is empty")
	ErrTokenNotFound     = errors.New("token not found")
	ErrRoomNotFound      = errors.New("room not found")
	ErrInvalidFile       = errors.New("invalid file")
	ErrRoomAlreadyExists = errors.New("room already exists")
)
