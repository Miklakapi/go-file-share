package files

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrFileNotFound  = errors.New("file not found")
	ErrEmptyToken    = errors.New("token is empty")
	ErrTokenNotFound = errors.New("token not found")
)

type FileRoom struct {
	ID        uuid.UUID
	ExpiresAt time.Time
	Files     map[uuid.UUID]*FileRoomFile

	tokens     map[string]bool
	password   string
	deleteOnce sync.Once
}

func NewFileRoom(hashedPassword string, creatorToken string, lifespanSec int) (*FileRoom, error) {
	if hashedPassword == "" {
		return nil, errors.New("password hash cannot be empty")
	}
	if creatorToken == "" {
		return nil, errors.New("creator token cannot be empty")
	}
	if lifespanSec <= 0 {
		return nil, errors.New("lifespan must be positive")
	}

	fileRoom := &FileRoom{
		ID:        uuid.New(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(lifespanSec)),
		Files:     make(map[uuid.UUID]*FileRoomFile),
		tokens:    make(map[string]bool),
		password:  hashedPassword,
	}

	fileRoom.tokens[creatorToken] = true

	return fileRoom, nil
}

func (r *FileRoom) HasToken(token string) bool {
	if token == "" || r.tokens == nil {
		return false
	}
	_, ok := r.tokens[token]
	return ok
}

func (r *FileRoom) AddToken(token string) error {
	if token == "" {
		return ErrEmptyToken
	}
	if r.tokens == nil {
		return ErrTokenNotFound
	}
	r.tokens[token] = true
	return nil
}

func (r *FileRoom) RemoveToken(token string) error {
	if token == "" {
		return ErrEmptyToken
	}
	if r.tokens == nil {
		return ErrTokenNotFound
	}
	if _, ok := r.tokens[token]; !ok {
		return ErrTokenNotFound
	}
	delete(r.tokens, token)
	return nil
}

func (r *FileRoom) AddFile(uploadDir string, name string, reader io.Reader) (*FileRoomFile, error) {
	file, err := NewFileRoomFile(uploadDir, name, reader)
	if err != nil {
		return nil, err
	}
	r.Files[file.ID] = file
	return file, nil
}

func (r *FileRoom) GetFile(id uuid.UUID) (*FileRoomFile, bool) {
	file, ok := r.Files[id]
	return file, ok
}

func (r *FileRoom) ListFiles() []*FileRoomFile {
	files := make([]*FileRoomFile, 0, len(r.Files))
	for _, file := range r.Files {
		files = append(files, file)
	}
	return files
}

func (r *FileRoom) DeleteFile(id uuid.UUID) error {
	file, ok := r.Files[id]
	if !ok {
		return ErrFileNotFound
	}

	if err := file.Delete(); err != nil {
		return err
	}
	delete(r.Files, id)
	return nil
}

func (r *FileRoom) IsExpired(now time.Time) bool {
	return now.After(r.ExpiresAt)
}

func (r *FileRoom) Delete() error {
	var joined error

	r.deleteOnce.Do(func() {
		r.tokens = nil
		r.password = ""

		for _, file := range r.Files {
			if err := file.Delete(); err != nil {
				joined = errors.Join(joined, err)
			}
		}

		r.Files = nil
	})

	return joined
}
