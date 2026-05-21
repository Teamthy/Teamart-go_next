package websocket

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 8
)

// Connection represents a single WebSocket client.
type Connection struct {
	ws     *websocket.Conn
	hub    *Hub
	send   chan []byte
	RoomID string
	UserID int64
}

// NewConnection creates a new websocket connection wrapper.
func NewConnection(ws *websocket.Conn, hub *Hub, roomID string, userID int64) *Connection {
	return &Connection{
		ws:     ws,
		hub:    hub,
		send:   make(chan []byte, 256),
		RoomID: roomID,
		UserID: userID,
	}
}

// ReadPump reads messages from the WebSocket connection.
func (c *Connection) ReadPump() {
	defer c.Close()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			return
		}
	}
}

// WritePump writes messages from the hub to the WebSocket connection.
func (c *Connection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.ws.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Close closes the underlying WebSocket connection.
func (c *Connection) Close() {
	c.hub.Unregister(c)
	c.ws.Close()
	close(c.send)
}
