package websocket

import "sync"

// RoomManager tracks room membership for websocket connections.
type RoomManager struct {
	mu    sync.RWMutex
	rooms map[string]map[*Connection]bool
}

// NewRoomManager creates a room manager.
func NewRoomManager() *RoomManager {
	return &RoomManager{rooms: make(map[string]map[*Connection]bool)}
}

// Join adds a connection to a room.
func (m *RoomManager) Join(roomID string, conn *Connection) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.rooms[roomID] == nil {
		m.rooms[roomID] = make(map[*Connection]bool)
	}
	m.rooms[roomID][conn] = true
}

// Leave removes a connection from a room.
func (m *RoomManager) Leave(roomID string, conn *Connection) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if members, ok := m.rooms[roomID]; ok {
		delete(members, conn)
		if len(members) == 0 {
			delete(m.rooms, roomID)
		}
	}
}

// Broadcast sends a message to all members of a room.
func (m *RoomManager) Broadcast(roomID string, payload []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for conn := range m.rooms[roomID] {
		select {
		case conn.send <- payload:
		default:
			conn.Close()
		}
	}
}
