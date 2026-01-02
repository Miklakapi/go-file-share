package eventbus

import "github.com/Miklakapi/go-file-share/internal/file-share/ports"

type EventBus struct {
	ports.EventPublisherSubscriber
}

func New() *EventBus {
	return &EventBus{}
}
