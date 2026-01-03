package eventbus

import (
	"sync"

	"github.com/Miklakapi/go-file-share/internal/file-share/ports"
)

type EventBus struct {
	mu     sync.RWMutex
	subs   map[ports.EventName]map[uint64]*sub
	nextID uint64
}

type sub struct {
	ch chan ports.Event
}

func New() *EventBus {
	return &EventBus{
		subs: make(map[ports.EventName]map[uint64]*sub),
	}
}

func (eb *EventBus) Subscribe(name ports.EventName) (<-chan ports.Event, ports.UnsubscribeFunc, error) {
	s := &sub{ch: make(chan ports.Event, 1)}

	eb.mu.Lock()
	eb.nextID++
	id := eb.nextID

	if eb.subs[name] == nil {
		eb.subs[name] = map[uint64]*sub{}
	}
	eb.subs[name][id] = s
	eb.mu.Unlock()

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			eb.mu.Lock()
			m := eb.subs[name]
			if m != nil {
				delete(m, id)
				if len(m) == 0 {
					delete(eb.subs, name)
				}
			}
			eb.mu.Unlock()
			close(s.ch)
		})
	}

	return s.ch, unsubscribe, nil
}

func (eb *EventBus) Publish(event ports.Event) error {
	eb.mu.RLock()
	subsMap := eb.subs[event.Name]
	subs := make([]*sub, 0, len(subsMap))
	for _, s := range subsMap {
		subs = append(subs, s)
	}
	eb.mu.RUnlock()

	for _, s := range subs {
		if trySend(s.ch, event) {
			continue
		}

		for tryDrainOne(s.ch) {
		}

		_ = trySend(s.ch, event)
	}

	return nil
}

func trySend(ch chan ports.Event, event ports.Event) (sent bool) {
	defer func() {
		_ = recover()
	}()

	select {
	case ch <- event:
		return true
	default:
		return false
	}
}
func tryDrainOne(ch chan ports.Event) (drained bool) {
	defer func() {
		_ = recover()
	}()

	select {
	case <-ch:
		return true
	default:
		return false
	}

}
