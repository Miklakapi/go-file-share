package sqliterepository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Miklakapi/go-file-share/internal/file-share/domain"
	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"github.com/google/uuid"
)

const defaultInLimit = 900

type SqliteRepo struct {
	db      *sql.DB
	inLimit int
}

var _ ports.RoomRepository = (*SqliteRepo)(nil)

func New(ctx context.Context, db *sql.DB) *SqliteRepo {
	r := &SqliteRepo{
		db:      db,
		inLimit: defaultInLimit,
	}

	r.inLimit = readMaxSQLVars(ctx, db, defaultInLimit)

	if r.inLimit > 50 {
		r.inLimit -= 50
	}
	if r.inLimit <= 0 {
		r.inLimit = defaultInLimit
	}

	return r
}

func (r *SqliteRepo) Get(ctx context.Context, roomID uuid.UUID) (*domain.Room, bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}

	roomIdString := roomID.String()

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = tx.Rollback() }()

	var (
		idStr        string
		passwordHash string
		expiresAtSec int64
	)

	err = tx.QueryRowContext(ctx, `
		SELECT id, password_hash, expires_at
		FROM rooms
		WHERE id = ?
		LIMIT 1
	`, roomIdString).Scan(&idStr, &passwordHash, &expiresAtSec)

	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, false, err
	}

	room := domain.HydrateRoom(id, passwordHash, time.Unix(expiresAtSec, 0))

	tokenRows, err := tx.QueryContext(ctx, `
		SELECT token
		FROM room_tokens
		WHERE room_id = ?
	`, roomIdString)
	if err != nil {
		return nil, false, err
	}
	defer tokenRows.Close()

	for tokenRows.Next() {
		var t string
		if err := tokenRows.Scan(&t); err != nil {
			return nil, false, err
		}
		if err := room.AddToken(t); err != nil {
			return nil, false, err
		}
	}
	if err := tokenRows.Err(); err != nil {
		return nil, false, err
	}

	fileRows, err := tx.QueryContext(ctx, `
		SELECT id, path, name, size, created_at
		FROM room_files
		WHERE room_id = ?
	`, roomIdString)
	if err != nil {
		return nil, false, err
	}
	defer fileRows.Close()

	for fileRows.Next() {
		var (
			fileIDStr    string
			path         string
			name         string
			size         int64
			createdAtSec int64
		)
		if err := fileRows.Scan(&fileIDStr, &path, &name, &size, &createdAtSec); err != nil {
			return nil, false, err
		}

		fid, err := uuid.Parse(fileIDStr)
		if err != nil {
			return nil, false, err
		}

		_ = room.AddFile(&domain.RoomFile{
			ID:        fid,
			Path:      path,
			Name:      name,
			Size:      size,
			CreatedAt: time.Unix(createdAtSec, 0),
		})
	}
	if err := fileRows.Err(); err != nil {
		return nil, false, err
	}

	if err := tx.Commit(); err != nil {
		return nil, false, err
	}

	return room, true, nil
}

func (r *SqliteRepo) List(ctx context.Context) ([]*domain.Room, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	rows, err := tx.QueryContext(ctx, `
		SELECT id, password_hash, expires_at
		FROM rooms
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roomByID := make(map[string]*domain.Room, 50)
	roomIDs := make([]string, 0, 50)

	for rows.Next() {
		var (
			idStr        string
			passwordHash string
			expiresAtSec int64
		)
		if err := rows.Scan(&idStr, &passwordHash, &expiresAtSec); err != nil {
			return nil, err
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		room := domain.HydrateRoom(id, passwordHash, time.Unix(expiresAtSec, 0))
		roomByID[idStr] = room
		roomIDs = append(roomIDs, idStr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(roomIDs) == 0 {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return []*domain.Room{}, nil
	}

	chunks := chunkStrings(roomIDs, r.inLimit)
	for _, ch := range chunks {
		q := fmt.Sprintf(`
			SELECT room_id, id, path, name, size, created_at
			FROM room_files
			WHERE room_id IN (%s)
		`, makePlaceholders(len(ch)))

		fRows, err := tx.QueryContext(ctx, q, argsFromStrings(ch)...)
		if err != nil {
			return nil, err
		}

		for fRows.Next() {
			var (
				roomIDStr    string
				fileIDStr    string
				path         string
				name         string
				size         int64
				createdAtSec int64
			)
			if err := fRows.Scan(&roomIDStr, &fileIDStr, &path, &name, &size, &createdAtSec); err != nil {
				_ = fRows.Close()
				return nil, err
			}

			room := roomByID[roomIDStr]
			if room == nil {
				continue
			}

			fid, err := uuid.Parse(fileIDStr)
			if err != nil {
				_ = fRows.Close()
				return nil, err
			}

			_ = room.AddFile(&domain.RoomFile{
				ID:        fid,
				Path:      path,
				Name:      name,
				Size:      size,
				CreatedAt: time.Unix(createdAtSec, 0),
			})
		}
		if err := fRows.Err(); err != nil {
			_ = fRows.Close()
			return nil, err
		}
		_ = fRows.Close()
	}

	for _, ch := range chunks {
		q := fmt.Sprintf(`
			SELECT room_id, token
			FROM room_tokens
			WHERE room_id IN (%s)
		`, makePlaceholders(len(ch)))

		tRows, err := tx.QueryContext(ctx, q, argsFromStrings(ch)...)
		if err != nil {
			return nil, err
		}

		for tRows.Next() {
			var (
				roomIDStr string
				token     string
			)
			if err := tRows.Scan(&roomIDStr, &token); err != nil {
				_ = tRows.Close()
				return nil, err
			}

			room := roomByID[roomIDStr]
			if room == nil {
				continue
			}
			if err := room.AddToken(token); err != nil {
				_ = tRows.Close()
				return nil, err
			}
		}
		if err := tRows.Err(); err != nil {
			_ = tRows.Close()
			return nil, err
		}
		_ = tRows.Close()
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	out := make([]*domain.Room, 0, len(roomIDs))
	for _, idStr := range roomIDs {
		if room := roomByID[idStr]; room != nil {
			out = append(out, room)
		}
	}

	return out, nil
}

func (r *SqliteRepo) Create(ctx context.Context, room *domain.Room) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if room == nil {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO rooms (id, password_hash, expires_at, created_at)
		VALUES (?, ?, ?, CAST(strftime('%s','now') AS INTEGER))
	`, room.ID.String(), room.Password(), room.ExpiresAt.Unix())
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ports.ErrRoomAlreadyExists
		}
		return err
	}

	tokens := room.ListTokens()
	if len(tokens) > 0 {
		stmt, err := tx.PrepareContext(ctx, `
			INSERT INTO room_tokens (room_id, token, created_at)
			VALUES (?, ?, CAST(strftime('%s','now') AS INTEGER))
		`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, t := range tokens {
			if t == "" {
				continue
			}
			_, err := stmt.ExecContext(ctx, room.ID.String(), t)
			if err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *SqliteRepo) Delete(ctx context.Context, roomID uuid.UUID) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	panic("TODO")
}

func (r *SqliteRepo) DeleteExpired(ctx context.Context, now time.Time) ([]domain.ExpiredCleanup, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	panic("TODO")
}

func (r *SqliteRepo) RemoveToken(ctx context.Context, roomID uuid.UUID, token string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	panic("TODO")
}

func (r *SqliteRepo) AddToken(ctx context.Context, roomID uuid.UUID, token string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	panic("TODO")
}

func (r *SqliteRepo) AddFileByToken(ctx context.Context, roomID uuid.UUID, token string, file *domain.RoomFile) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	panic("TODO")
}

func (r *SqliteRepo) DeleteFileByToken(ctx context.Context, roomID, fileID uuid.UUID, token string) (string, bool, error) {
	if err := ctx.Err(); err != nil {
		return "", false, err
	}

	panic("TODO")
}
