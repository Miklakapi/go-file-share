package ports

import "errors"

var (
	ErrRoomAlreadyExists = errors.New("room already exists")
	ErrRoomNotFound      = errors.New("room not found")

	ErrNilReader      = errors.New("file reader is nil")
	ErrEmptyFilename  = errors.New("filename is empty")
	ErrEmptyUploadDir = errors.New("upload dir is empty")

	ErrInvalidToken      = errors.New("token invalid")
	ErrTokenSignAlgo     = errors.New("unexpected signing method")
	ErrTokenExpired      = errors.New("token expired")
	ErrTokenParse        = errors.New("token parse error")
	ErrTokenRoomMismatch = errors.New("token room mismatch")
)
