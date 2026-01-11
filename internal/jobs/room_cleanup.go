package jobs

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	fileShare "github.com/Miklakapi/go-file-share/internal/file-share/application"
	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
	"github.com/google/uuid"
)

type RoomCleanupJob struct {
	fileShareService *fileShare.Service
	eventPublisher   ports.EventPublisher
	cleanupInterval  time.Duration
}

func New(fileShareService *fileShare.Service, eventPublisher ports.EventPublisher, cleanupInterval time.Duration) *RoomCleanupJob {
	return &RoomCleanupJob{
		fileShareService: fileShareService,
		eventPublisher:   eventPublisher,
		cleanupInterval:  cleanupInterval,
	}
}

func (r *RoomCleanupJob) Run(ctx context.Context) (func(), error) {
	closeChannel := make(chan struct{}, 1)
	var once sync.Once

	var close = func() {
		once.Do(func() {
			close(closeChannel)
		})
	}

	go r.cleanup(ctx, closeChannel, r.cleanupInterval)

	return close, nil
}

func (r *RoomCleanupJob) cleanup(ctx context.Context, close chan struct{}, duration time.Duration) {
	cleanupTicker := time.NewTicker(duration)
	defer cleanupTicker.Stop()

	for {
		select {
		case <-cleanupTicker.C:
			deletedRooms, err := r.fileShareService.CleanupExpired(ctx)
			if err != nil {
				log.Printf("file error: %v\n", err)
				continue
			}
			if len(deletedRooms) == 0 {
				continue
			}

			deletedRoomsString := uuidsToString(deletedRooms)
			if err := r.eventPublisher.Publish(ports.Event{Name: ports.EventRoomDelete, Data: deletedRoomsString}); err != nil {
				log.Printf("file error: %v\n", err)
				continue
			}
		case <-close:
			return

		case <-ctx.Done():
			return
		}
	}
}

func uuidsToString(uuids []uuid.UUID) string {
	var b strings.Builder
	for i, u := range uuids {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(u.String())
	}
	return b.String()
}
