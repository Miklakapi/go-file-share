package files

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FileHub struct {
	rooms map[uuid.UUID]*FileRoom
	cmds  chan any
}

func NewFileHub() *FileHub {
	return &FileHub{
		rooms: make(map[uuid.UUID]*FileRoom),
		cmds:  make(chan any, 1),
	}
}

type cmdRegisterRoom struct {
	room *FileRoom
}

type cmdUnregisterRoom struct {
	uuid uuid.UUID
}

func (fH *FileHub) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case c := <-fH.cmds:
			switch cmd := c.(type) {
			case cmdRegisterRoom:
				fH.rooms[cmd.room.UUID] = cmd.room
			case cmdUnregisterRoom:
				room := fH.rooms[cmd.uuid]
				if room != nil {
					room.Delete()
				}
				delete(fH.rooms, cmd.uuid)
			default:
				fmt.Printf("unknown file hub command: %T\n", c)
			}
		case <-ticker.C:
			now := time.Now()
			for uuid, room := range fH.rooms {
				if now.After(room.ExpiresAt) {
					room.Delete()
					delete(fH.rooms, uuid)
				}
			}
		case <-ctx.Done():
			for uuid, room := range fH.rooms {
				room.Delete()
				delete(fH.rooms, uuid)
			}
			return
		}
	}
}

func (h *FileHub) RegisterRoom(ctx context.Context, room *FileRoom) error {
	if room == nil {
		return nil
	}

	select {
	case h.cmds <- cmdRegisterRoom{room: room}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (h *FileHub) UnregisterRoom(ctx context.Context, uuid uuid.UUID) error {
	select {
	case h.cmds <- cmdUnregisterRoom{uuid: uuid}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
