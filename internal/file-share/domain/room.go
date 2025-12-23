package domain

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type FileRoom struct {
	ID        uuid.UUID
	ExpiresAt time.Time
	Files     map[uuid.UUID]*FileRoomFile

	tokens     map[string]bool
	password   string
	deleteOnce sync.Once
}

func NewFileRoom(hashedPassword string, creatorToken string, lifespan time.Duration) (*FileRoom, error) {
	if hashedPassword == "" {
		return nil, ErrEmptyPasswordHash
	}
	if creatorToken == "" {
		return nil, ErrEmptyCreatorToken
	}
	if lifespan <= 0 {
		return nil, ErrInvalidRoomTTL
	}

	now := time.Now()

	r := &FileRoom{
		ID:        uuid.New(),
		ExpiresAt: now.Add(lifespan),
		Files:     make(map[uuid.UUID]*FileRoomFile),
		tokens:    make(map[string]bool),
		password:  hashedPassword,
	}

	r.tokens[creatorToken] = true
	return r, nil
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

func (r *FileRoom) TokensCount() int {
	return len(r.tokens)
}

func (r *FileRoom) Password() string {
	return r.password
}

func (r *FileRoom) GetFile(id uuid.UUID) (*FileRoomFile, bool) {
	if r.Files == nil {
		return nil, false
	}
	file, ok := r.Files[id]
	return file, ok
}

func (r *FileRoom) ListFiles() []*FileRoomFile {
	if r.Files == nil {
		return nil
	}
	files := make([]*FileRoomFile, 0, len(r.Files))
	for _, f := range r.Files {
		files = append(files, f)
	}
	return files
}

func (r *FileRoom) AddFile(file *FileRoomFile) error {
	if file == nil {
		return ErrInvalidFile
	}
	if r.Files == nil {
		r.Files = make(map[uuid.UUID]*FileRoomFile)
	}
	r.Files[file.ID] = file
	return nil
}

func (r *FileRoom) DeleteFile(id uuid.UUID) (*FileRoomFile, error) {
	if r.Files == nil {
		return nil, ErrFileNotFound
	}

	file, ok := r.Files[id]
	if !ok {
		return nil, ErrFileNotFound
	}

	delete(r.Files, id)
	return file, nil
}

func (r *FileRoom) IsExpired(now time.Time) bool {
	return now.After(r.ExpiresAt)
}

func (r *FileRoom) Delete() (filesToCleanup []*FileRoomFile) {
	r.deleteOnce.Do(func() {
		for _, f := range r.Files {
			filesToCleanup = append(filesToCleanup, f)
		}

		r.Files = nil
		r.tokens = nil
		r.password = ""
	})

	return filesToCleanup
}
