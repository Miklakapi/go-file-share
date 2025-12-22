package files

import (
	"time"

	"github.com/google/uuid"
)

type cmdRoomsList struct {
	resp chan roomsSnapResp
}

type cmdRoomGet struct {
	id   uuid.UUID
	resp chan roomSnapResp
}

type cmdRoomCreate struct {
	hashedPassword string
	creatorToken   string
	lifespanSec    int
	resp           chan roomSnapResp
}

type cmdRoomDelete struct {
	id   uuid.UUID
	resp chan errResp
}

func (h *FileHub) roomList(cmd cmdRoomsList) {
	rooms := make([]RoomSnapshot, 0, len(h.rooms))
	for _, room := range h.rooms {
		rooms = append(rooms, RoomSnapshot{
			ID:        room.ID,
			ExpiresAt: room.ExpiresAt,
			Files:     len(room.Files),
			Tokens:    len(room.tokens),
		})
	}
	cmd.resp <- roomsSnapResp{rooms: rooms}
}

func (h *FileHub) getRoom(cmd cmdRoomGet) {
	room, ok := h.rooms[cmd.id]
	if !ok {
		cmd.resp <- roomSnapResp{ok: false}
		return
	}

	cmd.resp <- roomSnapResp{
		ok: true,
		room: RoomSnapshot{
			ID:        room.ID,
			ExpiresAt: room.ExpiresAt,
			Files:     len(room.Files),
			Tokens:    len(room.tokens),
		},
	}
}

func (h *FileHub) createRoom(cmd cmdRoomCreate) {
	room, err := NewFileRoom(cmd.hashedPassword, cmd.creatorToken, cmd.lifespanSec)
	if err != nil {
		cmd.resp <- roomSnapResp{ok: false, err: err}
		return
	}
	h.rooms[room.ID] = room

	cmd.resp <- roomSnapResp{
		ok: true,
		room: RoomSnapshot{
			ID:        room.ID,
			ExpiresAt: room.ExpiresAt,
			Files:     len(room.Files),
			Tokens:    len(room.tokens),
		},
	}
}

func (h *FileHub) deleteRoom(cmd cmdRoomDelete) {
	room, ok := h.rooms[cmd.id]
	if !ok {
		cmd.resp <- errResp{err: ErrRoomNotFound}
		return
	}

	room.Delete()
	delete(h.rooms, cmd.id)

	cmd.resp <- errResp{err: nil}
}

func (h *FileHub) ttlCleanup() {
	now := time.Now()
	for id, room := range h.rooms {
		if room.IsExpired(now) {
			room.Delete()
			delete(h.rooms, id)
		}
	}
}

func (h *FileHub) done() {
	for id, room := range h.rooms {
		room.Delete()
		delete(h.rooms, id)
	}
}
