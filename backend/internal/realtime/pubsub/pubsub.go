package pubsub

// Topic represents a pubsub topic name.
type Topic string

// PubSub is a minimal publish/subscribe interface used by the realtime gateway.
type PubSub interface {
    Subscribe(topic Topic) (<-chan []byte, func(), error)
    Publish(topic Topic, msg []byte) error
}
