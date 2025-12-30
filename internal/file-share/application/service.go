package application

import (
	"context"
	"errors"
	"io"
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
	policy      domain.Policy
	now         func() time.Time
}

func NewService(rooms ports.RoomRepository, files ports.FileStore, hasher ports.PasswordHasher, tokenIssuer ports.TokenService, policy domain.Policy) *Service {
	return &Service{
		rooms:       rooms,
		files:       files,
		hasher:      hasher,
		tokenIssuer: tokenIssuer,
		policy:      policy,
		now:         time.Now,
	}
}

func (s *Service) Room(ctx context.Context, id uuid.UUID) (*domain.Room, bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}

	room, ok, err := s.rooms.Get(ctx, id)
	if err != nil {
		return nil, false, err
	}
	if !ok || room == nil {
		return nil, false, nil
	}

	return room, true, nil
}

func (s *Service) CheckRoomAccess(ctx context.Context, id uuid.UUID, token string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	room, ok, err := s.rooms.Get(ctx, id)
	if err != nil {
		return false, err
	}
	if !ok || room == nil || !room.HasToken(token) {
		return false, nil
	}

	return true, nil
}

func (s *Service) Rooms(ctx context.Context) ([]*domain.Room, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return s.rooms.List(ctx)
}

func (s *Service) CreateRoom(ctx context.Context, password string, lifespan time.Duration) (*domain.Room, string, error) {
	if err := ctx.Err(); err != nil {
		return nil, "", err
	}

	password = strings.TrimSpace(password)
	if password == "" {
		return nil, "", domain.ErrEmptyPassword
	}

	if lifespan <= 0 {
		lifespan = s.policy.DefaultRoomTTL
	}
	if s.policy.MaxRoomLifespan > 0 && lifespan > s.policy.MaxRoomLifespan {
		return nil, "", domain.ErrRoomLifespanTooLong
	}

	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return nil, "", err
	}

	room, err := domain.NewRoom(hashedPassword, lifespan)
	if err != nil {
		return nil, "", err
	}

	token, _, err := s.tokenIssuer.Issue(ctx, room.ID, lifespan)
	if err != nil {
		return nil, "", err
	}
	if err := room.AddToken(token); err != nil {
		return nil, "", err
	}

	if err := s.rooms.Create(ctx, room); err != nil {
		return nil, "", err
	}

	return room, token, nil
}

func (s *Service) DeleteRoom(ctx context.Context, id uuid.UUID, token string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return domain.ErrEmptyToken
	}

	room, ok, err := s.rooms.Get(ctx, id)
	if err != nil {
		return err
	}
	if !ok || room == nil || !room.HasToken(token) {
		return domain.ErrRoomNotFound
	}

	paths, err := s.rooms.Delete(ctx, id)
	if err != nil {
		return err
	}

	var joined error
	for _, path := range paths {
		if path == "" {
			continue
		}
		if err := s.files.Delete(ctx, path); err != nil {
			joined = errors.Join(joined, err)
		}
	}

	return joined
}

func (s *Service) AuthRoom(ctx context.Context, id uuid.UUID, password string, lifespan time.Duration) (string, time.Time, error) {
	if err := ctx.Err(); err != nil {
		return "", time.Time{}, err
	}

	password = strings.TrimSpace(password)
	if password == "" {
		return "", time.Time{}, domain.ErrEmptyPassword
	}

	if lifespan <= 0 {
		lifespan = s.policy.DefaultTokenTTL
	}
	if s.policy.MaxTokenLifespan > 0 && lifespan > s.policy.MaxTokenLifespan {
		return "", time.Time{}, domain.ErrTokenLifespanTooLong
	}

	room, ok, err := s.rooms.Get(ctx, id)
	if err != nil {
		return "", time.Time{}, err
	}
	if !ok || room == nil {
		return "", time.Time{}, domain.ErrRoomNotFound
	}

	hash := room.Password()
	if !s.hasher.Verify(password, hash) {
		return "", time.Time{}, domain.ErrInvalidPassword
	}

	token, expiresAt, err := s.tokenIssuer.Issue(ctx, id, lifespan)
	if err != nil {
		return "", time.Time{}, err
	}

	if err := s.rooms.AddToken(ctx, id, token); err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func (s *Service) LogoutRoom(ctx context.Context, id uuid.UUID, token string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return domain.ErrEmptyToken
	}

	ok, err := s.rooms.RemoveToken(ctx, id, token)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrTokenNotFound
	}

	return nil
}

func (s *Service) File(ctx context.Context, roomId, fileId uuid.UUID, token string) (*domain.RoomFile, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return nil, domain.ErrEmptyToken
	}

	room, ok, err := s.rooms.Get(ctx, roomId)
	if err != nil {
		return nil, err
	}
	if !ok || room == nil || !room.HasToken(token) {
		return nil, domain.ErrRoomNotFound
	}

	f, ok := room.GetFile(fileId)
	if !ok || f == nil {
		return nil, domain.ErrFileNotFound
	}

	return f, nil
}

func (s *Service) DownloadFile(ctx context.Context, roomId, fileId uuid.UUID, token string) (*domain.RoomFile, io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return nil, nil, domain.ErrEmptyToken
	}

	room, ok, err := s.rooms.Get(ctx, roomId)
	if err != nil {
		return nil, nil, err
	}
	if !ok || room == nil || !room.HasToken(token) {
		return nil, nil, domain.ErrRoomNotFound
	}

	file, ok := room.GetFile(fileId)
	if !ok || file == nil {
		return nil, nil, domain.ErrFileNotFound
	}

	rc, err := s.files.Open(ctx, file.Path)
	if err != nil {
		return nil, nil, err
	}

	return file, rc, nil
}

func (s *Service) Files(ctx context.Context, id uuid.UUID, token string) ([]*domain.RoomFile, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return nil, domain.ErrEmptyToken
	}

	room, ok, err := s.rooms.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok || room == nil || !room.HasToken(token) {
		return nil, domain.ErrRoomNotFound
	}

	files := room.ListFiles()
	return files, nil
}

func (s *Service) UploadFile(ctx context.Context, roomId uuid.UUID, token string, filename string, r io.Reader) (*domain.RoomFile, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return nil, domain.ErrEmptyToken
	}

	filename = strings.TrimSpace(filename)
	if filename == "" {
		return nil, ports.ErrEmptyFilename
	}
	if r == nil {
		return nil, ports.ErrNilReader
	}

	path, size, err := s.files.Save(ctx, s.policy.UploadDir, filename, r)
	if err != nil {
		return nil, err
	}

	now := s.now()
	meta, err := domain.NewRoomFile(path, filename, size, now)
	if err != nil {
		_ = s.files.Delete(ctx, path)
		return nil, err
	}

	ok, err := s.rooms.AddFileByToken(ctx, roomId, token, meta)
	if err != nil {
		_ = s.files.Delete(ctx, path)
		return nil, err
	}
	if !ok {
		_ = s.files.Delete(ctx, path)
		return nil, domain.ErrRoomNotFound
	}

	return meta, nil
}

func (s *Service) DeleteFile(ctx context.Context, roomId, fileId uuid.UUID, token string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return domain.ErrEmptyToken
	}

	path, ok, err := s.rooms.DeleteFileByToken(ctx, roomId, fileId, token)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrFileNotFound
	}

	if err := s.files.Delete(ctx, path); err != nil {
		return err
	}

	return nil
}

func (s *Service) CleanupExpired(ctx context.Context) ([]uuid.UUID, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	expired, err := s.rooms.DeleteExpired(ctx, s.now())
	if err != nil {
		return nil, err
	}

	var deleted []uuid.UUID
	var joined error

	for _, item := range expired {
		deleted = append(deleted, item.RoomID)

		for _, path := range item.Paths {
			if strings.TrimSpace(path) == "" {
				continue
			}
			if err := s.files.Delete(ctx, path); err != nil {
				joined = errors.Join(joined, err)
			}
		}
	}

	return deleted, joined
}
