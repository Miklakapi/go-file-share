PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS rooms (
  id            TEXT PRIMARY KEY,
  password_hash TEXT NOT NULL,
  expires_at    INTEGER NOT NULL,
  created_at    INTEGER NOT NULL
    DEFAULT (CAST(strftime('%s','now') AS INTEGER))
);

CREATE TABLE IF NOT EXISTS room_files (
  id         TEXT PRIMARY KEY,
  room_id    TEXT NOT NULL,
  path       TEXT NOT NULL,
  name       TEXT NOT NULL,
  size       INTEGER NOT NULL,
  created_at INTEGER NOT NULL,

  FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS room_tokens (
  room_id    TEXT NOT NULL,
  token      TEXT NOT NULL,
  created_at INTEGER NOT NULL
    DEFAULT (CAST(strftime('%s','now') AS INTEGER)),

  PRIMARY KEY (room_id, token),
  FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);