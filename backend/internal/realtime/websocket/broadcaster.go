package websocket

// Broadcaster is a helper wrapper for hub broadcast methods.
type Broadcaster struct {
	hub *Hub
}

// NewBroadcaster creates a new broadcaster.
func NewBroadcaster(hub *Hub) *Broadcaster {
	return &Broadcaster{hub: hub}
}

// SendToRoom broadcasts a payload to a room.
func (b *Broadcaster) SendToRoom(roomID string, payload []byte) {
	b.hub.Broadcast(roomID, payload)
}
