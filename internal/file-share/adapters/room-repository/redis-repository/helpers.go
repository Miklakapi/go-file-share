package redisrepository

import "github.com/google/uuid"

func roomKey(roomID uuid.UUID) string {
	return "room:" + roomID.String()
}

func tokensKey(roomID uuid.UUID) string {
	return roomKey(roomID) + ":tokens"
}

func filesKey(roomID uuid.UUID) string {
	return roomKey(roomID) + ":files"
}
