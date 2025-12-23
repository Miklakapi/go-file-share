package adapters

import "errors"

var (
	ErrRoomAlreadyExists = errors.New("room already exists")
	ErrRoomNotFound      = errors.New("room not found")

	ErrNilReader      = errors.New("file reader is nil")
	ErrEmptyFilename  = errors.New("filename is empty")
	ErrEmptyUploadDir = errors.New("upload dir is empty")
)
