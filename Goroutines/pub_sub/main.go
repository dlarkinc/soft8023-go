package main

import (
	"fmt"
	"sync"
	"time"
)

// Subscriber receives messages via a personal channel
type Subscriber struct {
	ID int
	C  chan string
}

// Publisher sends messages to all subscribers
type Publisher struct {
	subscribers []Subscriber
	mu          sync.Mutex
}

// AddSubscriber creates and registers a new subscriber
func (p *Publisher) AddSubscriber(id int) Subscriber {
	sub := Subscriber{ID: id, C: make(chan string)}
	p.mu.Lock()
	p.subscribers = append(p.subscribers, sub)
	p.mu.Unlock()
	return sub
}

// Publish sends a message to all subscribers (non-blocking)
func (p *Publisher) Publish(msg string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, sub := range p.subscribers {
		go func(c chan string) {
			c <- msg
		}(sub.C)
	}
}

func main() {
	p := &Publisher{}

	// Add two subscribers
	sub1 := p.AddSubscriber(1)
	sub2 := p.AddSubscriber(2)

	// Each subscriber runs in its own goroutine
	go func() {
		for msg := range sub1.C {
			fmt.Printf("[Sub %d] %s\n", sub1.ID, msg)
		}
	}()
	go func() {
		for msg := range sub2.C {
			fmt.Printf("[Sub %d] %s\n", sub2.ID, msg)
		}
	}()

	// Publish some messages
	for i := 1; i <= 3; i++ {
		p.Publish(fmt.Sprintf("Update %d", i))
		time.Sleep(500 * time.Millisecond)
	}

	// Let subscribers finish printing
	time.Sleep(time.Second)
}
