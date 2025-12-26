package application

import (
	"context"
	"strings"
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"github.com/google/uuid"
)

type Service struct {
	rooms       ports.RoomRepository
	files       ports.FileStore
	hasher      ports.PasswordHasher
	tokenIssuer ports.TokenService
	settings    Settings
	now         func() time.Time
}

func NewService(rooms ports.RoomRepository, files ports.FileStore, hasher ports.PasswordHasher, tokenIssuer ports.TokenService, settings Settings) *Service {
	return &Service{
		rooms:       rooms,
		files:       files,
		hasher:      hasher,
		tokenIssuer: tokenIssuer,
		settings:    settings,
		now:         time.Now,
	}
}

func (s *Service) Room(ctx context.Context, id uuid.UUID) (domain.RoomSnapshot, bool, error) {
	if err := ctx.Err(); err != nil {
		return domain.RoomSnapshot{}, false, err
	}

	room, ok, err := s.rooms.Get(ctx, id)
	if err != nil {
		return domain.RoomSnapshot{}, false, err
	}
	if !ok || room == nil {
		return domain.RoomSnapshot{}, false, nil
	}

	snap := domain.RoomSnapshot{
		ID:        room.ID,
		ExpiresAt: room.ExpiresAt,
		Files:     len(room.Files),
		Tokens:    room.TokensCount(),
	}

	return snap, true, nil
}

func (s *Service) Rooms(ctx context.Context) ([]domain.RoomSnapshot, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return s.rooms.ListSnapshots(ctx)
}

func (s *Service) CreateRoom(ctx context.Context, password string, lifespan time.Duration, roomID uuid.UUID) (domain.RoomSnapshot, string, error) {
	if err := ctx.Err(); err != nil {
		return domain.RoomSnapshot{}, "", err
	}

	password = strings.TrimSpace(password)
	if password == "" {
		return domain.RoomSnapshot{}, "", domain.ErrEmptyPassword
	}

	if lifespan <= 0 {
		lifespan = s.settings.DefaultRoomTTL
	}
	if s.settings.MaxRoomLifespan > 0 && lifespan > s.settings.MaxRoomLifespan {
		return domain.RoomSnapshot{}, "", domain.ErrRoomLifespanTooLong
	}

	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return domain.RoomSnapshot{}, "", err
	}

	token, _, err := s.tokenIssuer.Issue(ctx, roomID, lifespan)
	if err != nil {
		return domain.RoomSnapshot{}, "", err
	}

	room, err := domain.NewFileRoom(hashedPassword, token, lifespan)
	if err != nil {
		return domain.RoomSnapshot{}, "", err
	}

	if err := s.rooms.Create(ctx, room); err != nil {
		return domain.RoomSnapshot{}, "", err
	}

	snap := domain.RoomSnapshot{
		ID:        room.ID,
		ExpiresAt: room.ExpiresAt,
		Files:     len(room.Files),
		Tokens:    room.TokensCount(),
	}

	return snap, token, nil
}

func (s *Service) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return s.rooms.Delete(ctx, id)
}

func (s *Service) AuthRoom(ctx context.Context, id uuid.UUID, password string) (token string, expiresAt time.Time, err error) {
	panic("TODO")
}

func (s *Service) LogoutRoom(ctx context.Context, id uuid.UUID, token string) error {
	panic("TODO")
}

func (s *Service) CleanupExpired(ctx context.Context) ([]uuid.UUID, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return s.rooms.DeleteExpired(ctx, s.now())
}
