package websocket

import (
	"sync"
)

// Hub coordinates WebSocket connections and broadcasts.
type Hub struct {
	mu          sync.RWMutex
	connections map[*Connection]bool
	rooms       map[string]map[*Connection]bool
	register    chan *Connection
	unregister  chan *Connection
	broadcast   chan *BroadcastMessage
}

// BroadcastMessage wraps a payload for delivery.
type BroadcastMessage struct {
	RoomID  string
	Payload []byte
}

// NewHub creates a new WebSocket hub.
func NewHub() *Hub {
	return &Hub{
		connections: make(map[*Connection]bool),
		rooms:       make(map[string]map[*Connection]bool),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		broadcast:   make(chan *BroadcastMessage, 128),
	}
}

// Run starts the hub event loop.
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.connections[conn] = true
			if conn.RoomID != "" {
				if h.rooms[conn.RoomID] == nil {
					h.rooms[conn.RoomID] = make(map[*Connection]bool)
				}
				h.rooms[conn.RoomID][conn] = true
			}
			h.mu.Unlock()
		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.connections[conn]; ok {
				delete(h.connections, conn)
			}
			if conn.RoomID != "" {
				if room, ok := h.rooms[conn.RoomID]; ok {
					delete(room, conn)
					if len(room) == 0 {
						delete(h.rooms, conn.RoomID)
					}
				}
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			room, ok := h.rooms[message.RoomID]
			if ok {
				for conn := range room {
					select {
					case conn.send <- message.Payload:
					default:
						conn.Close()
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register adds a connection to the hub.
func (h *Hub) Register(conn *Connection) {
	h.register <- conn
}

// Unregister removes a connection from the hub.
func (h *Hub) Unregister(conn *Connection) {
	h.unregister <- conn
}

// Broadcast sends a message to all connections in a room.
func (h *Hub) Broadcast(roomID string, payload []byte) {
	h.broadcast <- &BroadcastMessage{RoomID: roomID, Payload: payload}
}
