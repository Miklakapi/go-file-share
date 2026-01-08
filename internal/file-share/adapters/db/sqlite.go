package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SqliteDB struct {
	Conn *sql.DB
}

func NewSqlite(path string) (*SqliteDB, error) {
	dsn := path

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("pragma foreign_keys: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &SqliteDB{Conn: db}, nil
}
