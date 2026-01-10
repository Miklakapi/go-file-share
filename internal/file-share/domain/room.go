package domain

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID        uuid.UUID
	ExpiresAt time.Time
	Files     map[uuid.UUID]*RoomFile

	tokens   map[string]bool
	password string
}

func NewRoom(hashedPassword string, lifespan time.Duration) (*Room, error) {
	if hashedPassword == "" {
		return nil, ErrEmptyPasswordHash
	}
	if lifespan <= 0 {
		return nil, ErrInvalidRoomTTL
	}

	now := time.Now()

	r := &Room{
		ID:        uuid.New(),
		ExpiresAt: now.Add(lifespan),
		Files:     make(map[uuid.UUID]*RoomFile),
		tokens:    make(map[string]bool, 1),
		password:  hashedPassword,
	}

	return r, nil
}

func HydrateRoom(id uuid.UUID, passwordHash string, expiresAt time.Time) *Room {
	return &Room{
		ID:        id,
		ExpiresAt: expiresAt,
		Files:     make(map[uuid.UUID]*RoomFile),
		tokens:    make(map[string]bool),
		password:  passwordHash,
	}
}

func (r *Room) HasToken(token string) bool {
	if token == "" || r.tokens == nil {
		return false
	}
	_, ok := r.tokens[token]
	return ok
}

func (r *Room) AddToken(token string) error {
	if token == "" {
		return ErrEmptyToken
	}
	if r.tokens == nil {
		r.tokens = make(map[string]bool)
	}
	r.tokens[token] = true
	return nil
}

func (r *Room) RemoveToken(token string) error {
	if token == "" {
		return ErrEmptyToken
	}
	if r.tokens == nil {
		r.tokens = make(map[string]bool)
	}
	if _, ok := r.tokens[token]; !ok {
		return ErrTokenNotFound
	}
	delete(r.tokens, token)
	return nil
}

func (r *Room) TokensCount() int {
	return len(r.tokens)
}

func (r *Room) Password() string {
	return r.password
}

func (r *Room) GetFile(id uuid.UUID) (*RoomFile, bool) {
	if r.Files == nil {
		return nil, false
	}
	file, ok := r.Files[id]
	return file, ok
}

func (r *Room) ListFiles() []*RoomFile {
	if r.Files == nil {
		return nil
	}
	files := make([]*RoomFile, 0, len(r.Files))
	for _, f := range r.Files {
		files = append(files, f)
	}
	return files
}

func (r *Room) AddFile(file *RoomFile) error {
	if file == nil {
		return ErrInvalidFile
	}
	if r.Files == nil {
		r.Files = make(map[uuid.UUID]*RoomFile)
	}
	r.Files[file.ID] = file
	return nil
}

func (r *Room) DeleteFile(id uuid.UUID) (*RoomFile, error) {
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

func (r *Room) IsExpired(now time.Time) bool {
	return now.After(r.ExpiresAt)
}

func (r *Room) Clone() *Room {
	if r == nil {
		return nil
	}

	cp := &Room{
		ID:        r.ID,
		ExpiresAt: r.ExpiresAt,
		Files:     make(map[uuid.UUID]*RoomFile, len(r.Files)),
		tokens:    make(map[string]bool, len(r.tokens)),
		password:  r.password,
	}

	for id, f := range r.Files {
		if f == nil {
			continue
		}
		ff := *f
		cp.Files[id] = &ff
	}

	for t := range r.tokens {
		cp.tokens[t] = true
	}

	return cp
}
