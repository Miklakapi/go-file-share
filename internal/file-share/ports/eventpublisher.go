package ports

import "context"

type EventPublisher interface {
	Publish(ctx context.Context, event any) error
}
