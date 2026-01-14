package redisrepository

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	db *redis.Client
}

func New(db *redis.Client) *RedisRepo {
	return &RedisRepo{
		db: db,
	}
}

func (r *RedisRepo) Get(ctx context.Context, roomID uuid.UUID) (*domain.Room, bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}

	k := roomKey(roomID)

	m, err := r.db.HGetAll(ctx, k).Result()
	if err != nil {
		return nil, false, err
	}
	if len(m) == 0 {
		return nil, false, nil
	}

	expiresAtStr, ok := m["expires_at"]
	if !ok {
		return nil, false, nil
	}
	expiresAt, err := strconv.ParseInt(expiresAtStr, 10, 64)
	if err != nil {
		return nil, false, err
	}

	room := domain.HydrateRoom(
		roomID,
		m["password_hash"],
		time.Unix(expiresAt, 0),
	)

	tokens, err := r.db.SMembers(ctx, k+":tokens").Result()
	if err != nil {
		return nil, false, err
	}
	for _, t := range tokens {
		if err := room.AddToken(t); err != nil {
			return nil, false, err
		}
	}

	files, err := r.db.HGetAll(ctx, k+":files").Result()
	if err != nil {
		return nil, false, err
	}
	for _, raw := range files {
		var f domain.RoomFile
		if err := json.Unmarshal([]byte(raw), &f); err == nil {
			_ = room.AddFile(&f)
		}
	}

	return room, true, nil
}

func (r *RedisRepo) List(ctx context.Context) ([]*domain.Room, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	iter := r.db.Scan(ctx, 0, "room:*", 0).Iterator()
	rooms := make([]*domain.Room, 0, 50)

	for iter.Next(ctx) {
		key := iter.Val()

		if strings.HasSuffix(key, ":tokens") || strings.HasSuffix(key, ":files") {
			continue
		}

		m, err := r.db.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		if len(m) == 0 {
			continue
		}

		roomIDStr := strings.TrimPrefix(key, "room:")
		roomID, err := uuid.Parse(roomIDStr)
		if err != nil {
			continue
		}

		expiresAtStr, ok := m["expires_at"]
		if !ok {
			continue
		}
		expiresAt, err := strconv.ParseInt(expiresAtStr, 10, 64)
		if err != nil {
			continue
		}

		room := domain.HydrateRoom(
			roomID,
			m["password_hash"],
			time.Unix(expiresAt, 0),
		)

		tokens, err := r.db.SMembers(ctx, key+":tokens").Result()
		if err != nil {
			return nil, err
		}
		for _, t := range tokens {
			_ = room.AddToken(t)
		}

		files, err := r.db.HGetAll(ctx, key+":files").Result()
		if err != nil {
			return nil, err
		}
		for _, raw := range files {
			var f domain.RoomFile
			if err := json.Unmarshal([]byte(raw), &f); err == nil {
				_ = room.AddFile(&f)
			}
		}

		rooms = append(rooms, room)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *RedisRepo) Create(ctx context.Context, room *domain.Room) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if room == nil {
		return nil
	}

	k := roomKey(room.ID)

	err := r.db.Watch(ctx, func(tx *redis.Tx) error {
		exists, err := tx.Exists(ctx, k).Result()
		if err != nil {
			return err
		}
		if exists > 0 {
			return ports.ErrRoomAlreadyExists
		}

		nowSec := time.Now().Unix()

		_, err = tx.TxPipelined(ctx, func(p redis.Pipeliner) error {
			p.HSet(ctx, k,
				"password_hash", room.Password(),
				"expires_at", room.ExpiresAt.Unix(),
				"created_at", nowSec,
			)

			tokens := room.ListTokens()
			if len(tokens) > 0 {
				args := make([]any, 0, len(tokens))
				for _, t := range tokens {
					if t == "" {
						continue
					}
					args = append(args, t)
				}
				if len(args) > 0 {
					p.SAdd(ctx, tokensKey(room.ID), args...)
				}
			}

			return nil
		})
		return err
	}, k)

	return err
}

func (r *RedisRepo) Delete(ctx context.Context, roomID uuid.UUID) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	kRoom := roomKey(roomID)
	kFiles := filesKey(roomID)
	kTokens := tokensKey(roomID)

	ex, err := r.db.Exists(ctx, kRoom).Result()
	if err != nil {
		return nil, err
	}
	if ex == 0 {
		return nil, ports.ErrRoomNotFound
	}

	files, err := r.db.HGetAll(ctx, kFiles).Result()
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(files))
	for _, raw := range files {
		var f domain.RoomFile
		if err := json.Unmarshal([]byte(raw), &f); err == nil {
			if f.Path != "" {
				paths = append(paths, f.Path)
			}
		}
	}

	if err := r.db.Del(ctx, kRoom, kFiles, kTokens).Err(); err != nil {
		return nil, err
	}

	return paths, nil
}

func (r *RedisRepo) DeleteExpired(ctx context.Context, now time.Time) ([]domain.ExpiredCleanup, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	nowSec := now.Unix()

	iter := r.db.Scan(ctx, 0, "room:*", 0).Iterator()

	out := make([]domain.ExpiredCleanup, 0, 50)

	for iter.Next(ctx) {
		key := iter.Val()

		if strings.HasSuffix(key, ":tokens") || strings.HasSuffix(key, ":files") {
			continue
		}

		expStr, err := r.db.HGet(ctx, key, "expires_at").Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			return nil, err
		}

		exp, err := strconv.ParseInt(expStr, 10, 64)
		if err != nil {
			continue
		}

		if exp >= nowSec {
			continue
		}

		roomIDStr := strings.TrimPrefix(key, "room:")
		roomID, err := uuid.Parse(roomIDStr)
		if err != nil {
			continue
		}

		fk := key + ":files"
		files, err := r.db.HGetAll(ctx, fk).Result()
		if err != nil {
			return nil, err
		}

		paths := make([]string, 0, len(files))
		for _, raw := range files {
			var f domain.RoomFile
			if err := json.Unmarshal([]byte(raw), &f); err == nil {
				if f.Path != "" {
					paths = append(paths, f.Path)
				}
			}
		}

		if err := r.db.Del(ctx, key, fk, key+":tokens").Err(); err != nil {
			return nil, err
		}

		out = append(out, domain.ExpiredCleanup{
			RoomID: roomID,
			Paths:  paths,
		})
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *RedisRepo) RemoveToken(ctx context.Context, roomID uuid.UUID, token string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	if token == "" {
		return false, nil
	}

	removed, err := r.db.SRem(ctx, tokensKey(roomID), token).Result()
	if err != nil {
		return false, err
	}
	return removed > 0, nil
}

func (r *RedisRepo) AddToken(ctx context.Context, roomID uuid.UUID, token string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if token == "" {
		return domain.ErrEmptyToken
	}

	exists, err := r.db.Exists(ctx, roomKey(roomID)).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return ports.ErrRoomNotFound
	}

	return r.db.SAdd(ctx, tokensKey(roomID), token).Err()
}

func (r *RedisRepo) AddFileByToken(ctx context.Context, roomID uuid.UUID, token string, file *domain.RoomFile) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	if file == nil {
		return false, domain.ErrInvalidFile
	}

	exists, err := r.db.Exists(ctx, roomKey(roomID)).Result()
	if err != nil {
		return false, err
	}
	if exists == 0 {
		return false, nil
	}

	ok, err := r.db.SIsMember(ctx, tokensKey(roomID), token).Result()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	raw, err := json.Marshal(file)
	if err != nil {
		return false, err
	}

	if err := r.db.HSet(ctx, filesKey(roomID), file.ID.String(), string(raw)).Err(); err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisRepo) DeleteFileByToken(ctx context.Context, roomID, fileID uuid.UUID, token string) (string, bool, error) {
	if err := ctx.Err(); err != nil {
		return "", false, err
	}

	exists, err := r.db.Exists(ctx, roomKey(roomID)).Result()
	if err != nil {
		return "", false, err
	}
	if exists == 0 {
		return "", false, nil
	}

	ok, err := r.db.SIsMember(ctx, tokensKey(roomID), token).Result()
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}

	fk := filesKey(roomID)
	field := fileID.String()

	raw, err := r.db.HGet(ctx, fk, field).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", false, nil
		}
		return "", false, err
	}

	var f domain.RoomFile
	if err := json.Unmarshal([]byte(raw), &f); err != nil {
		_, _ = r.db.HDel(ctx, fk, field).Result()
		return "", false, err
	}

	removed, err := r.db.HDel(ctx, fk, field).Result()
	if err != nil {
		return "", false, err
	}
	if removed == 0 {
		return "", false, nil
	}

	return f.Path, true, nil
}

func (r *RedisRepo) WipeAll(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return r.db.FlushDB(ctx).Err()
}
