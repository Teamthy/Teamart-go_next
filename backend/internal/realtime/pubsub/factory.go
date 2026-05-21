package pubsub

import (
	"context"
	"fmt"
	"strings"
)

// NewPubSubFromConfig returns a PubSub implementation based on a config string.
// Supported values (case-insensitive): "redis", "nats", "kafka", "memory".
// By default this returns the in-memory broker.
func NewPubSubFromConfig(ctx context.Context, backend string) (PubSub, error) {
	switch strings.ToLower(strings.TrimSpace(backend)) {
	case "redis":
		return nil, fmt.Errorf("redis adapter not available in this build (use -tags pubsub_redis)")
	case "nats":
		return nil, fmt.Errorf("nats adapter not available in this build (use -tags pubsub_nats)")
	case "kafka":
		return nil, fmt.Errorf("kafka adapter not available in this build (use -tags pubsub_kafka)")
	case "memory":
		fallthrough
	default:
		return NewInMemoryPubSub(), nil
	}
}
