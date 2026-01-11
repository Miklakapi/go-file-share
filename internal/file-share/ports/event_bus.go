package ports

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
	Publish(event Event) error
}

type UnsubscribeFunc func()

type EventSubscriber interface {
	Subscribe(name EventName) (<-chan Event, UnsubscribeFunc, error)
}

type EventPublisherSubscriber interface {
	EventPublisher
	EventSubscriber
}
