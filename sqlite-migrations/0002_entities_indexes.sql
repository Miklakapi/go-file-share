PRAGMA foreign_keys = ON;

CREATE INDEX IF NOT EXISTS idx_rooms_expires_at ON rooms(expires_at);

CREATE INDEX IF NOT EXISTS idx_room_files_room_id ON room_files(room_id);

CREATE INDEX IF NOT EXISTS idx_room_tokens_room_id ON room_tokens(room_id);