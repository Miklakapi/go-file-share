package ports

import "context"

type EventName string

const (
	EventRoomCreate EventName = "RoomCreate"
	EventRoomDelete EventName = "RoomDelete"
)

type Event struct {
	Name EventName
	Data any
}

type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}

type UnsubscribeFunc func()

type EventSubscriber interface {
	Subscribe(ctx context.Context, name EventName) (<-chan Event, UnsubscribeFunc, error)
}

type EventPublisherSubscriber interface {
	EventPublisher
	EventSubscriber
}
