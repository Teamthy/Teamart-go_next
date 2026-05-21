package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/teamart/commerce-api/internal/realtime/pubsub"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 16 * 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Authenticator interface {
	AuthenticateRequest(*http.Request) (int64, error)
}

type websocketMessage struct {
	Type      string          `json:"type"`
	Topic     string          `json:"topic,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	RequestID string          `json:"request_id,omitempty"`
}

type serverMessage struct {
	Type      string      `json:"type"`
	Topic     string      `json:"topic,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
	Error     string      `json:"error,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

type WebsocketServer struct {
	hub           *Hub
	authenticator Authenticator
}

func NewWebsocketServer(hub *Hub, authenticator Authenticator) *WebsocketServer {
	return &WebsocketServer{hub: hub, authenticator: authenticator}
}

func (s *WebsocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, err := s.authenticator.AuthenticateRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clientID := fmt.Sprintf("%d-%d", userID, time.Now().UnixNano())
	client := &Client{
		ID:            clientID,
		UserID:        userID,
		Subscriptions: make(map[pubsub.Topic]struct{}),
		Send:          make(chan interface{}, 256),
	}

	s.hub.RegisterClient(client)

	initialTopic := r.URL.Query().Get("topic")
	if initialTopic != "" {
		_ = s.hub.SubscribeClientToTopic(r.Context(), clientID, pubsub.Topic(initialTopic))
	}

	go s.writePump(ws, client)
	s.readPump(ws, client)
}

func (s *WebsocketServer) readPump(ws *websocket.Conn, client *Client) {
	defer func() {
		s.hub.UnregisterClient(client.ID)
		ws.Close()
	}()

	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, raw, err := ws.ReadMessage()
		if err != nil {
			return
		}

		var msg websocketMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			s.sendClientError(client, "invalid_message", err.Error(), "")
			continue
		}

		s.handleClientMessage(ws, client, &msg)
	}
}

func (s *WebsocketServer) writePump(ws *websocket.Conn, client *Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.Close()
	}()

	for {
		select {
		case payload, ok := <-client.Send:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(payload)
			if err != nil {
				continue
			}

			if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (s *WebsocketServer) handleClientMessage(ws *websocket.Conn, client *Client, message *websocketMessage) {
	topic := pubsub.Topic(message.Topic)

	switch message.Type {
	case "subscribe":
		if topic == "" {
			s.sendClientError(client, "missing_topic", "topic is required for subscribe", message.RequestID)
			return
		}
		if err := s.hub.SubscribeClientToTopic(context.Background(), client.ID, topic); err != nil {
			s.sendClientError(client, "subscribe_failed", err.Error(), message.RequestID)
			return
		}
		s.sendClientAck(client, "subscribed", topic, message.RequestID)

	case "unsubscribe":
		if topic == "" {
			s.sendClientError(client, "missing_topic", "topic is required for unsubscribe", message.RequestID)
			return
		}
		if err := s.hub.UnsubscribeClientFromTopic(client.ID, topic); err != nil {
			s.sendClientError(client, "unsubscribe_failed", err.Error(), message.RequestID)
			return
		}
		s.sendClientAck(client, "unsubscribed", topic, message.RequestID)

	case "publish":
		if topic == "" {
			s.sendClientError(client, "missing_topic", "topic is required for publish", message.RequestID)
			return
		}
		var payload interface{}
		if len(message.Payload) > 0 {
			if err := json.Unmarshal(message.Payload, &payload); err != nil {
				s.sendClientError(client, "invalid_payload", err.Error(), message.RequestID)
				return
			}
		}
		if err := s.hub.Broadcast(context.Background(), topic, payload); err != nil {
			s.sendClientError(client, "publish_failed", err.Error(), message.RequestID)
			return
		}
		s.sendClientAck(client, "published", topic, message.RequestID)

	case "typing":
		if topic == "" {
			s.sendClientError(client, "missing_topic", "topic is required for typing events", message.RequestID)
			return
		}
		var payload struct {
			Typing bool `json:"typing"`
		}
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			s.sendClientError(client, "invalid_payload", err.Error(), message.RequestID)
			return
		}
		if err := s.hub.Broadcast(context.Background(), topic, map[string]interface{}{
			"type":      "typing",
			"user_id":   client.UserID,
			"typing":    payload.Typing,
			"timestamp": time.Now().UTC(),
		}); err != nil {
			s.sendClientError(client, "typing_failed", err.Error(), message.RequestID)
			return
		}
		s.sendClientAck(client, "typing", topic, message.RequestID)

	case "ping":
		s.sendClientAck(client, "pong", topic, message.RequestID)

	default:
		s.sendClientError(client, "unknown_type", "unsupported message type", message.RequestID)
	}
}

func (s *WebsocketServer) sendClientAck(client *Client, messageType string, topic pubsub.Topic, requestID string) {
	client.Send <- serverMessage{Type: messageType, Topic: string(topic), RequestID: requestID}
}

func (s *WebsocketServer) sendClientError(client *Client, code, message, requestID string) {
	client.Send <- serverMessage{Type: "error", Error: fmt.Sprintf("%s: %s", code, message), RequestID: requestID}
}
