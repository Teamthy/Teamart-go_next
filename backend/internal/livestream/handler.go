package livestream

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

type Handler struct {
	service *Service
	logger  *logger.Logger
}

func NewHandler(service *Service, logger *logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func RegisterLivestreamRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("/api/v1/livestreams", handler.handleCollection)
	mux.HandleFunc("/api/v1/livestreams/", handler.handleItem)
}

func (h *Handler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleList(w, r)
	case http.MethodPost:
		h.handleCreate(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleItem(w http.ResponseWriter, r *http.Request) {
	id, suffix := parseStreamPath(r.URL.Path)
	if id == "" {
		writeJSONError(w, http.StatusNotFound, "stream id required")
		return
	}

	switch {
	case suffix == "" && r.Method == http.MethodGet:
		h.handleGet(w, r, id)
	case suffix == "/metadata" && r.Method == http.MethodPatch:
		h.handleUpdateMetadata(w, r, id)
	case suffix == "/state" && r.Method == http.MethodPost:
		h.handleTransition(w, r, id)
	case suffix == "/viewer" && r.Method == http.MethodPost:
		h.handleViewer(w, r, id)
	case suffix == "/engagement" && r.Method == http.MethodPost:
		h.handleEngagement(w, r, id)
	default:
		writeJSONError(w, http.StatusNotFound, "not found")
	}
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID           string      `json:"id"`
		Title        string      `json:"title"`
		ThumbnailURL string      `json:"thumbnail_url,omitempty"`
		Category     string      `json:"category,omitempty"`
		Tags         []string    `json:"tags,omitempty"`
		CreatorID    int64       `json:"creator_id"`
		CreatorName  string      `json:"creator_name,omitempty"`
		CoHosts      []string    `json:"co_hosts,omitempty"`
		ScheduledAt  *string     `json:"scheduled_at,omitempty"`
		State        StreamState `json:"state,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	metadata := StreamMetadata{
		ID:           req.ID,
		Title:        req.Title,
		ThumbnailURL: req.ThumbnailURL,
		Category:     req.Category,
		Tags:         req.Tags,
		CreatorID:    req.CreatorID,
		CreatorName:  req.CreatorName,
		CoHosts:      req.CoHosts,
	}
	if req.ScheduledAt != nil && *req.ScheduledAt != "" {
		scheduledAt, err := parseTime(*req.ScheduledAt)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid scheduled_at format")
			return
		}
		metadata.ScheduledAt = scheduledAt
	}
	info, err := h.service.CreateStream(metadata, req.State)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, info)
}

func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.service.ListStreams())
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request, streamID string) {
	info, ok := h.service.GetStream(streamID)
	if !ok {
		writeJSONError(w, http.StatusNotFound, "stream not found")
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func (h *Handler) handleUpdateMetadata(w http.ResponseWriter, r *http.Request, streamID string) {
	var req struct {
		Title        string   `json:"title,omitempty"`
		ThumbnailURL string   `json:"thumbnail_url,omitempty"`
		Category     string   `json:"category,omitempty"`
		Tags         []string `json:"tags,omitempty"`
		CreatorName  string   `json:"creator_name,omitempty"`
		CoHosts      []string `json:"co_hosts,omitempty"`
		ScheduledAt  *string  `json:"scheduled_at,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	metadata := StreamMetadata{
		Title:        req.Title,
		ThumbnailURL: req.ThumbnailURL,
		Category:     req.Category,
		Tags:         req.Tags,
		CreatorName:  req.CreatorName,
		CoHosts:      req.CoHosts,
	}
	if req.ScheduledAt != nil && *req.ScheduledAt != "" {
		scheduledAt, err := parseTime(*req.ScheduledAt)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid scheduled_at format")
			return
		}
		metadata.ScheduledAt = scheduledAt
	}
	info, err := h.service.UpdateMetadata(streamID, metadata)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func (h *Handler) handleTransition(w http.ResponseWriter, r *http.Request, streamID string) {
	var req struct {
		State StreamState `json:"state"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	info, err := h.service.TransitionState(streamID, req.State)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func (h *Handler) handleViewer(w http.ResponseWriter, r *http.Request, streamID string) {
	var req struct {
		UserID int64  `json:"user_id"`
		Action string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.UserID == 0 || req.Action == "" {
		writeJSONError(w, http.StatusBadRequest, "user_id and action are required")
		return
	}

	var analytics *StreamAnalytics
	var err error
	switch strings.ToLower(req.Action) {
	case "join":
		analytics, err = h.service.AddViewer(streamID, req.UserID)
	case "leave":
		analytics, err = h.service.RemoveViewer(streamID, req.UserID)
	default:
		writeJSONError(w, http.StatusBadRequest, "action must be join or leave")
		return
	}
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, analytics)
}

func (h *Handler) handleEngagement(w http.ResponseWriter, r *http.Request, streamID string) {
	var req struct {
		Type EngagementType `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Type == "" {
		writeJSONError(w, http.StatusBadRequest, "type is required")
		return
	}
	analytics, err := h.service.TrackEngagement(streamID, req.Type)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, analytics)
}

func parseStreamPath(path string) (id string, suffix string) {
	trimmed := strings.TrimPrefix(path, "/api/v1/livestreams/")
	if trimmed == "" {
		return "", ""
	}
	parts := strings.SplitN(trimmed, "/", 2)
	id = parts[0]
	if len(parts) > 1 {
		suffix = "/" + parts[1]
	}
	return id, suffix
}

func parseTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
