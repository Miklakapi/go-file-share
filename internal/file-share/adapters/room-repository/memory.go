package roomrepository

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
	rooms map[uuid.UUID]*domain.FileRoom
}

var _ ports.RoomRepository = (*MemoryRepo)(nil)

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		rooms: make(map[uuid.UUID]*domain.FileRoom),
	}
}

func (r *MemoryRepo) Get(ctx context.Context, roomID uuid.UUID) (*domain.FileRoom, bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}

	r.mu.RLock()
	room, ok := r.rooms[roomID]
	r.mu.RUnlock()

	if !ok || room == nil {
		return nil, false, nil
	}

	return cloneRoom(room), true, nil
}

func (r *MemoryRepo) GetByToken(ctx context.Context, roomID uuid.UUID, token string) (*domain.FileRoom, bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}

	if token == "" {
		return nil, false, nil
	}

	r.mu.RLock()
	room, ok := r.rooms[roomID]
	r.mu.RUnlock()

	if !ok || room == nil {
		return nil, false, nil
	}

	if !room.HasToken(token) {
		return nil, false, nil
	}

	return cloneRoom(room), true, nil
}

func (r *MemoryRepo) ListSnapshots(ctx context.Context) ([]domain.RoomSnapshot, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.RoomSnapshot, 0, len(r.rooms))
	for _, room := range r.rooms {
		if room == nil {
			continue
		}
		out = append(out, domain.RoomSnapshot{
			ID:        room.ID,
			ExpiresAt: room.ExpiresAt,
			Files:     len(room.Files),
			Tokens:    room.TokensCount(),
		})
	}

	return out, nil
}

func (r *MemoryRepo) Create(ctx context.Context, room *domain.FileRoom) error {
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

func (r *MemoryRepo) Update(ctx context.Context, room *domain.FileRoom) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if room == nil {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.rooms[room.ID]; !ok {
		return ports.ErrRoomNotFound
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

func (r *MemoryRepo) GetPasswordHash(ctx context.Context, roomID uuid.UUID) (hash string, ok bool, err error) {
	if err := ctx.Err(); err != nil {
		return "", false, err
	}

	r.mu.RLock()
	room, ok := r.rooms[roomID]
	r.mu.RUnlock()

	if !ok || room == nil {
		return "", false, nil
	}

	return room.Password(), true, nil
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

func cloneRoom(src *domain.FileRoom) *domain.FileRoom {
	if src == nil {
		return nil
	}

	dst := &domain.FileRoom{
		ID:        src.ID,
		ExpiresAt: src.ExpiresAt,
		Files:     make(map[uuid.UUID]*domain.FileRoomFile, len(src.Files)),
	}

	for id, f := range src.Files {
		dst.Files[id] = cloneFile(f)
	}

	return dst
}

func cloneFile(f *domain.FileRoomFile) *domain.FileRoomFile {
	if f == nil {
		return nil
	}
	cp := *f
	return &cp
}
