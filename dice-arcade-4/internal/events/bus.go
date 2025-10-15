package events

import (
	"sync"
	"time"
)

// Bus is an abstraction we will keep across implementations
// (e.g. switching to a message queue later)
type Bus interface {
	Publish(topic string, e Event) error
	Subscribe(topic string) (<-chan Event, func(), error) // returns ch + unsubscribe
}

type chanBus struct {
	mu sync.RWMutex
	// subscribers per topic; we also support a "*" wildcard (fan-out for all)
	subs map[string]map[chan Event]struct{}
}
