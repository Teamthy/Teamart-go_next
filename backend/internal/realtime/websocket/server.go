package websocket

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/teamart/commerce-api/internal/livestream"
	"github.com/teamart/commerce-api/internal/realtime/gateway"
	"github.com/teamart/commerce-api/internal/realtime/pubsub"
	"github.com/teamart/commerce-api/pkg/logger"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 30 * time.Second
	pingPeriod = 20 * time.Second
)

// Server handles websocket connections for realtime clients.
type Server struct {
	hub           *gateway.Hub
	authenticator Authenticator
	streamService *livestream.Service
	upgrader      websocket.Upgrader
	logger        *logger.Logger
}

// NewServer creates a websocket server bound to a hub, authenticator, and optional livestream service.
func NewServer(hub *gateway.Hub, authenticator Authenticator, streamService *livestream.Service, logger *logger.Logger) *Server {
	return &Server{
		hub:           hub,
		authenticator: authenticator,
		streamService: streamService,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		logger: logger,
	}
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

	client := &gateway.Client{
		ID:            newClientID(),
		UserID:        claims.UserID,
		Send:          make(chan []byte, 32),
		Subscriptions: make(map[pubsub.Topic]struct{}),
	}
	s.hub.RegisterClient(client)

	defer s.hub.UnregisterClient(client.ID)

	go s.writePump(client, conn)
	s.readPump(client, conn)
}

func (s *Server) readPump(client *gateway.Client, conn *websocket.Conn) {
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
			if err := s.hub.SubscribeClient(client.ID, pubsub.Topic(msg.Topic)); err != nil {
				s.sendError(client, "subscribe_error", err.Error())
				continue
			}
			s.sendAck(client, "subscribed", msg.Topic)
		case "unsubscribe":
			if msg.Topic == "" {
				s.sendError(client, "unsubscribe_error", "topic is required")
				continue
			}
			s.hub.UnsubscribeClient(client.ID, pubsub.Topic(msg.Topic))
			s.sendAck(client, "unsubscribed", msg.Topic)
		case "publish":
			if msg.Topic == "" {
				s.sendError(client, "publish_error", "topic is required")
				continue
			}
			if msg.Payload == nil {
				s.sendError(client, "publish_error", "payload is required")
				continue
			}
			payload, err := json.Marshal(msg.Payload)
			if err != nil {
				s.sendError(client, "publish_error", err.Error())
				continue
			}
			if err := s.hub.Publish(pubsub.Topic(msg.Topic), payload); err != nil {
				s.sendError(client, "publish_error", err.Error())
				continue
			}
			s.sendAck(client, "published", msg.Topic)
		case "livestream_event":
			if s.streamService == nil {
				s.sendError(client, "livestream_error", "livestream service unavailable")
				continue
			}
			if err := s.handleLivestreamEvent(client, msg.Payload); err != nil {
				s.sendError(client, "livestream_error", err.Error())
			}
		case "ping":
			s.sendPong(client)
		default:
			s.sendError(client, "unsupported_type", "unsupported message type")
		}
	}
}

func (s *Server) writePump(client *gateway.Client, conn *websocket.Conn) {
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

func (s *Server) sendAck(client *gateway.Client, eventType, topic string) {
	s.sendEvent(client, outgoingEvent{Type: eventType, Topic: topic, Success: true})
}

func (s *Server) sendPong(client *gateway.Client) {
	s.sendEvent(client, outgoingEvent{Type: "pong"})
}

func (s *Server) sendError(client *gateway.Client, eventType, message string) {
	s.sendEvent(client, outgoingEvent{Type: eventType, Error: message})
}

func (s *Server) handleLivestreamEvent(client *gateway.Client, payload json.RawMessage) error {
	var event struct {
		Action         string                    `json:"action"`
		StreamID       string                    `json:"stream_id"`
		UserID         int64                     `json:"user_id,omitempty"`
		EngagementType livestream.EngagementType `json:"engagement_type,omitempty"`
		State          livestream.StreamState    `json:"state,omitempty"`
	}
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}
	if event.StreamID == "" {
		return fmt.Errorf("stream_id is required")
	}

	switch strings.ToLower(event.Action) {
	case "viewer_join":
		analytics, err := s.streamService.AddViewer(event.StreamID, event.UserID)
		if err != nil {
			return err
		}
		s.sendEvent(client, outgoingEvent{Type: "livestream_update", Payload: mustMarshalJSON(analytics)})
	case "viewer_leave":
		analytics, err := s.streamService.RemoveViewer(event.StreamID, event.UserID)
		if err != nil {
			return err
		}
		s.sendEvent(client, outgoingEvent{Type: "livestream_update", Payload: mustMarshalJSON(analytics)})
	case "track_engagement":
		analytics, err := s.streamService.TrackEngagement(event.StreamID, event.EngagementType)
		if err != nil {
			return err
		}
		s.sendEvent(client, outgoingEvent{Type: "livestream_update", Payload: mustMarshalJSON(analytics)})
	case "transition_state":
		info, err := s.streamService.TransitionState(event.StreamID, event.State)
		if err != nil {
			return err
		}
		s.sendEvent(client, outgoingEvent{Type: "livestream_update", Payload: mustMarshalJSON(info)})
	default:
		return fmt.Errorf("unsupported livestream action")
	}
	return nil
}

func mustMarshalJSON(value interface{}) json.RawMessage {
	data, _ := json.Marshal(value)
	return data
}

func (s *Server) sendEvent(client *gateway.Client, event outgoingEvent) {
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
