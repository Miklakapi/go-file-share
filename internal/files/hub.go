package files

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrRoomNotFound = errors.New("room not found")

type FileHub struct {
	rooms map[uuid.UUID]*FileRoom
	cmds  chan any
}

func NewFileHub() *FileHub {
	return &FileHub{
		rooms: make(map[uuid.UUID]*FileRoom),
		cmds:  make(chan any, 10),
	}
}

func (fH *FileHub) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case c := <-fH.cmds:
			switch cmd := c.(type) {
			case cmdRoomsList:
				fH.roomList(cmd)
			case cmdRoomGet:
				fH.getRoom(cmd)
			case cmdRoomCreate:
				fH.createRoom(cmd)
			case cmdRoomDelete:
				fH.deleteRoom(cmd)
			default:
				fmt.Printf("unknown file hub command: %T\n", c)
			}
		case <-ticker.C:
			fH.ttlCleanup()
		case <-ctx.Done():
			fH.done()
			return
		}
	}
}

func (h *FileHub) GetRoomList(ctx context.Context) ([]RoomSnapshot, error) {
	resp := make(chan roomsSnapResp, 1)

	select {
	case h.cmds <- cmdRoomsList{resp: resp}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	select {
	case r := <-resp:
		return r.rooms, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (h *FileHub) GetRoom(ctx context.Context, id uuid.UUID) (RoomSnapshot, bool, error) {
	resp := make(chan roomSnapResp, 1)

	select {
	case h.cmds <- cmdRoomGet{id: id, resp: resp}:
	case <-ctx.Done():
		return RoomSnapshot{}, false, ctx.Err()
	}

	select {
	case r := <-resp:
		return r.room, r.ok, nil
	case <-ctx.Done():
		return RoomSnapshot{}, false, ctx.Err()
	}
}

func (h *FileHub) CreateRoom(ctx context.Context, hashedPassword, creatorToken string, lifespanSec int) (RoomSnapshot, error) {
	resp := make(chan roomSnapResp, 1)

	select {
	case h.cmds <- cmdRoomCreate{hashedPassword, creatorToken, lifespanSec, resp}:
	case <-ctx.Done():
		return RoomSnapshot{}, ctx.Err()
	}

	select {
	case r := <-resp:
		if r.err != nil {
			return RoomSnapshot{}, r.err
		}
		return r.room, nil
	case <-ctx.Done():
		return RoomSnapshot{}, ctx.Err()
	}
}

func (h *FileHub) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	resp := make(chan errResp, 1)

	select {
	case h.cmds <- cmdRoomDelete{id: id, resp: resp}:
	case <-ctx.Done():
		return ctx.Err()
	}

	select {
	case r := <-resp:
		return r.err
	case <-ctx.Done():
		return ctx.Err()
	}
}
