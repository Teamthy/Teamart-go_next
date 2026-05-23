package websocket

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/teamart/commerce-api/internal/moderation"
	"github.com/teamart/commerce-api/pkg/logger"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 30 * time.Second
	pingPeriod = 20 * time.Second
)

// Server handles websocket connections for realtime clients.
type Server struct {
	hub               *Hub
	authenticator     Authenticator
	moderationService *moderation.ModerationService
	upgrader          websocket.Upgrader
	logger            *logger.Logger
}

// NewServer creates a websocket server bound to a hub, authenticator, and optional moderation service.
func NewServer(hub *Hub, authenticator Authenticator, moderationService *moderation.ModerationService, logger *logger.Logger) *Server {
	return &Server{
		hub:               hub,
		authenticator:     authenticator,
		moderationService: moderationService,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		logger: logger,
	}
}

// Topic is a channel name for websocket pub/sub.
type Topic string

// Client represents a connected websocket client.
type Client struct {
	ID            string
	UserID        int64
	Send          chan []byte
	Subscriptions map[Topic]struct{}
}

// Hub coordinates websocket clients and topic broadcasts.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client
	rooms   map[Topic]map[string]*Client
}

// NewHub creates a new in-memory websocket hub.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		rooms:   make(map[Topic]map[string]*Client),
	}
}

// RegisterClient adds a client to the hub.
func (h *Hub) RegisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client.Subscriptions == nil {
		client.Subscriptions = make(map[Topic]struct{})
	}
	h.clients[client.ID] = client
}

// UnregisterClient removes a client and clears subscriptions.
func (h *Hub) UnregisterClient(clientID string) {
	h.mu.Lock()
	client, ok := h.clients[clientID]
	if !ok {
		h.mu.Unlock()
		return
	}

	for topic := range client.Subscriptions {
		h.removeClientFromTopicLocked(client, topic)
	}

	delete(h.clients, clientID)
	close(client.Send)
	h.mu.Unlock()
}

// SubscribeClient subscribes a client to a topic.
func (h *Hub) SubscribeClient(clientID string, topic Topic) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	client, ok := h.clients[clientID]
	if !ok {
		return fmt.Errorf("client not found")
	}

	if _, subscribed := client.Subscriptions[topic]; subscribed {
		return nil
	}

	client.Subscriptions[topic] = struct{}{}
	if _, ok := h.rooms[topic]; !ok {
		h.rooms[topic] = make(map[string]*Client)
	}
	h.rooms[topic][client.ID] = client
	return nil
}

// UnsubscribeClient unsubscribes a client from a topic.
func (h *Hub) UnsubscribeClient(clientID string, topic Topic) {
	h.mu.Lock()
	defer h.mu.Unlock()

	client, ok := h.clients[clientID]
	if !ok {
		return
	}

	if _, subscribed := client.Subscriptions[topic]; !subscribed {
		return
	}

	delete(client.Subscriptions, topic)
	h.removeClientFromTopicLocked(client, topic)
}

func (h *Hub) removeClientFromTopicLocked(client *Client, topic Topic) {
	if members, ok := h.rooms[topic]; ok {
		delete(members, client.ID)
		if len(members) == 0 {
			delete(h.rooms, topic)
		}
	}
}

// Publish broadcasts a message to all clients subscribed to a topic.
func (h *Hub) Publish(topic Topic, msg []byte) error {
	h.mu.RLock()
	clients, ok := h.rooms[topic]
	if !ok {
		h.mu.RUnlock()
		return nil
	}

	for _, client := range clients {
		payload := make([]byte, len(msg))
		copy(payload, msg)
		select {
		case client.Send <- payload:
		default:
		}
	}
	h.mu.RUnlock()
	return nil
}

// ServeHTTP upgrades the HTTP connection into a websocket after authentication.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	claims, err := s.authenticator.Authenticate(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Errorf("websocket upgrade failed: %v", err)
		return
	}

	client := &Client{
		ID:            newClientID(),
		UserID:        claims.UserID,
		Send:          make(chan []byte, 32),
		Subscriptions: make(map[Topic]struct{}),
	}
	s.hub.RegisterClient(client)

	defer s.hub.UnregisterClient(client.ID)

	go s.writePump(client, conn)
	s.readPump(client, conn)
}

func (s *Server) readPump(client *Client, conn *websocket.Conn) {
	conn.SetReadLimit(1024 * 8)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var msg clientMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			s.sendError(client, "invalid_message", err.Error())
			continue
		}

		switch msg.Type {
		case "subscribe":
			if msg.Topic == "" {
				s.sendError(client, "subscribe_error", "topic is required")
				continue
			}
			if err := s.hub.SubscribeClient(client.ID, Topic(msg.Topic)); err != nil {
				s.sendError(client, "subscribe_error", err.Error())
				continue
			}
			s.sendAck(client, "subscribed", msg.Topic)
		case "unsubscribe":
			if msg.Topic == "" {
				s.sendError(client, "unsubscribe_error", "topic is required")
				continue
			}
			s.hub.UnsubscribeClient(client.ID, Topic(msg.Topic))
			s.sendAck(client, "unsubscribed", msg.Topic)
		case "publish", "chat_message":
			if msg.Topic == "" {
				s.sendError(client, "publish_error", "topic is required")
				continue
			}
			if msg.Payload == nil {
				s.sendError(client, "publish_error", "payload is required")
				continue
			}
			if strings.HasPrefix(msg.Topic, "chat:") || msg.Type == "chat_message" {
				if err := s.handleChatPublish(client, msg.Topic, msg.Payload); err != nil {
					s.sendError(client, "chat_error", err.Error())
					continue
				}
				s.sendAck(client, "message_sent", msg.Topic)
				continue
			}
			payload, err := json.Marshal(msg.Payload)
			if err != nil {
				s.sendError(client, "publish_error", err.Error())
				continue
			}
			if err := s.hub.Publish(Topic(msg.Topic), payload); err != nil {
				s.sendError(client, "publish_error", err.Error())
				continue
			}
			s.sendAck(client, "published", msg.Topic)
		case "ping":
			s.sendPong(client)
		default:
			s.sendError(client, "unsupported_type", "unsupported message type")
		}
	}
}

func (s *Server) writePump(client *Client, conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case msg, ok := <-client.Send:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (s *Server) sendAck(client *Client, eventType, topic string) {
	s.sendEvent(client, outgoingEvent{Type: eventType, Topic: topic, Success: true})
}

func (s *Server) sendPong(client *Client) {
	s.sendEvent(client, outgoingEvent{Type: "pong"})
}

func (s *Server) sendError(client *Client, eventType, message string) {
	s.sendEvent(client, outgoingEvent{Type: eventType, Error: message})
}

func (s *Server) handleChatPublish(client *Client, topic string, rawPayload json.RawMessage) error {
	var payload struct {
		UserID int64          `json:"user_id,omitempty"`
		Text   string         `json:"text"`
		Meta   map[string]any `json:"meta,omitempty"`
	}
	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		return err
	}
	if payload.Text == "" {
		return fmt.Errorf("text is required")
	}
	if payload.UserID == 0 {
		payload.UserID = client.UserID
	}
	if s.moderationService != nil {
		decision := s.moderationService.EvaluateMessage(payload.UserID, payload.Text, topic)
		if !decision.Allowed {
			s.sendEvent(client, outgoingEvent{Type: "chat_rejected", Error: "message rejected by moderation", Payload: mustMarshalJSON(decision)})
			return nil
		}
		if decision.ShadowBan {
			s.sendEvent(client, outgoingEvent{Type: "chat_shadow", Payload: rawPayload})
			return nil
		}
	}
	if err := s.hub.Publish(Topic(topic), rawPayload); err != nil {
		return err
	}
	return nil
}

func mustMarshalJSON(value interface{}) json.RawMessage {
	data, _ := json.Marshal(value)
	return data
}

func (s *Server) sendEvent(client *Client, event outgoingEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	select {
	case client.Send <- data:
	default:
	}
}

type clientMessage struct {
	Type    string          `json:"type"`
	Topic   string          `json:"topic,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type outgoingEvent struct {
	Type    string          `json:"type"`
	Topic   string          `json:"topic,omitempty"`
	Success bool            `json:"success,omitempty"`
	Error   string          `json:"error,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

func newClientID() string {
	return fmt.Sprintf("client-%d-%d", time.Now().UnixNano(), rand.Int63())
}
