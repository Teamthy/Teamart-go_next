package pubsub

import (
    "sync"
)

// InMemoryBroker is a simple in-memory pub/sub broker suitable for tests and single-process use.
type InMemoryBroker struct {
    mu   sync.RWMutex
    subs map[Topic]map[chan []byte]struct{}
}

// NewInMemoryBroker constructs a new InMemoryBroker.
func NewInMemoryBroker() *InMemoryBroker {
    return &InMemoryBroker{
        subs: make(map[Topic]map[chan []byte]struct{}),
    }
}

// Subscribe subscribes to a topic and returns a receive-only channel and an unsubscribe function.
func (b *InMemoryBroker) Subscribe(topic Topic) (<-chan []byte, func(), error) {
    ch := make(chan []byte, 16)
    b.mu.Lock()
    defer b.mu.Unlock()
    m, ok := b.subs[topic]
    if !ok {
        m = make(map[chan []byte]struct{})
        b.subs[topic] = m
    }
    m[ch] = struct{}{}

    unsub := func() {
        b.mu.Lock()
        defer b.mu.Unlock()
        if m, ok := b.subs[topic]; ok {
            delete(m, ch)
            close(ch)
            if len(m) == 0 {
                delete(b.subs, topic)
            }
        }
    }

    return ch, unsub, nil
}

// Publish sends a message to all subscribers of a topic.
func (b *InMemoryBroker) Publish(topic Topic, msg []byte) error {
    b.mu.RLock()
    defer b.mu.RUnlock()
    m, ok := b.subs[topic]
    if !ok {
        return nil
    }
    for ch := range m {
        // non-blocking send
        select {
        case ch <- msg:
        default:
        }
    }
    return nil
}
