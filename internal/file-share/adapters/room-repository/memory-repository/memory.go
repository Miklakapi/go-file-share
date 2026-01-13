package memoryrepository

import (
	"context"
	"sync"
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"github.com/google/uuid"
)

type MemoryRepo struct {
	mu    sync.RWMutex
	rooms map[uuid.UUID]*domain.Room
}

func New() *MemoryRepo {
	return &MemoryRepo{
		rooms: make(map[uuid.UUID]*domain.Room),
	}
}

func (r *MemoryRepo) Get(ctx context.Context, roomID uuid.UUID) (*domain.Room, bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}

	r.mu.RLock()
	room, ok := r.rooms[roomID]
	var cp *domain.Room
	if ok && room != nil {
		cp = room.Clone()
	}
	r.mu.RUnlock()

	if cp == nil {
		return nil, false, nil
	}

	return cp, true, nil
}

func (r *MemoryRepo) List(ctx context.Context) ([]*domain.Room, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*domain.Room, 0, len(r.rooms))
	for _, room := range r.rooms {
		if room == nil {
			continue
		}
		out = append(out, room.Clone())
	}
	return out, nil
}

func (r *MemoryRepo) Create(ctx context.Context, room *domain.Room) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if room == nil {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.rooms[room.ID]; ok {
		return ports.ErrRoomAlreadyExists
	}

	r.rooms[room.ID] = room.Clone()
	return nil
}

func (r *MemoryRepo) Delete(ctx context.Context, roomID uuid.UUID) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	room, ok := r.rooms[roomID]
	if !ok {
		return nil, ports.ErrRoomNotFound
	}

	paths := make([]string, 0)
	if room != nil {
		paths = make([]string, 0, len(room.Files))
		for _, f := range room.Files {
			if f == nil {
				continue
			}
			if f.Path != "" {
				paths = append(paths, f.Path)
			}
		}
	}

	delete(r.rooms, roomID)
	return paths, nil
}

func (r *MemoryRepo) DeleteExpired(ctx context.Context, now time.Time) ([]domain.ExpiredCleanup, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]domain.ExpiredCleanup, 0)
	for id, room := range r.rooms {
		if room == nil {
			delete(r.rooms, id)
			out = append(out, domain.ExpiredCleanup{
				RoomID: id,
				Paths:  nil,
			})
			continue
		}

		if !now.After(room.ExpiresAt) {
			continue
		}

		paths := make([]string, 0, len(room.Files))
		for _, f := range room.Files {
			if f == nil {
				continue
			}
			if f.Path != "" {
				paths = append(paths, f.Path)
			}
		}

		delete(r.rooms, id)
		out = append(out, domain.ExpiredCleanup{
			RoomID: id,
			Paths:  paths,
		})
	}

	return out, nil
}

func (r *MemoryRepo) RemoveToken(ctx context.Context, roomID uuid.UUID, token string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if token == "" {
		return false, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	room, ok := r.rooms[roomID]
	if !ok || room == nil {
		return false, nil
	}

	if !room.HasToken(token) {
		return false, nil
	}

	if err := room.RemoveToken(token); err != nil {
		return false, err
	}

	return true, nil
}

func (r *MemoryRepo) AddToken(ctx context.Context, roomID uuid.UUID, token string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if token == "" {
		return domain.ErrEmptyToken
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	room, ok := r.rooms[roomID]
	if !ok || room == nil {
		return ports.ErrRoomNotFound
	}

	return room.AddToken(token)
}

func (r *MemoryRepo) AddFileByToken(ctx context.Context, roomID uuid.UUID, token string, file *domain.RoomFile) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if file == nil {
		return false, domain.ErrInvalidFile
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	room, ok := r.rooms[roomID]
	if !ok || room == nil {
		return false, nil
	}

	if !room.HasToken(token) {
		return false, nil
	}

	cp := *file
	if room.Files == nil {
		room.Files = make(map[uuid.UUID]*domain.RoomFile)
	}
	room.Files[cp.ID] = &cp

	return true, nil
}

func (r *MemoryRepo) DeleteFileByToken(ctx context.Context, roomID, fileID uuid.UUID, token string) (string, bool, error) {
	if err := ctx.Err(); err != nil {
		return "", false, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	room, ok := r.rooms[roomID]
	if !ok || room == nil {
		return "", false, nil
	}

	if !room.HasToken(token) {
		return "", false, nil
	}

	f, ok := room.Files[fileID]
	if !ok || f == nil {
		return "", false, nil
	}

	path := f.Path
	delete(room.Files, fileID)

	return path, true, nil
}
