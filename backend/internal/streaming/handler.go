package streaming

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/teamart/commerce-api/pkg/logger"
)

type Handler struct {
	service *Service
	logger  *logger.Logger
}

func NewHandler(service *Service, logger *logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func RegisterStreamingRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("/api/v1/streaming", handler.handleCollection)
	mux.HandleFunc("/api/v1/streaming/", handler.handleItem)
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
	case suffix == "/ingest" && r.Method == http.MethodPost:
		h.handleIngest(w, r, id)
	case suffix == "/transcode" && r.Method == http.MethodPost:
		h.handleTranscode(w, r, id)
	case suffix == "/playback" && r.Method == http.MethodGet:
		h.handlePlayback(w, r, id)
	case suffix == "/stop" && r.Method == http.MethodPost:
		h.handleStop(w, r, id)
	default:
		writeJSONError(w, http.StatusNotFound, "not found")
	}
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID    string `json:"id,omitempty"`
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		writeJSONError(w, http.StatusBadRequest, "title is required")
		return
	}

	info, err := h.service.CreateSession(req.Title, req.ID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, info)
}

func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.service.ListSessions())
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request, id string) {
	info, ok := h.service.GetSession(id)
	if !ok {
		writeJSONError(w, http.StatusNotFound, "stream not found")
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func (h *Handler) handleIngest(w http.ResponseWriter, r *http.Request, id string) {
	info, err := h.service.StartIngest(id)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func (h *Handler) handleTranscode(w http.ResponseWriter, r *http.Request, id string) {
	info, err := h.service.StartTranscoding(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func (h *Handler) handlePlayback(w http.ResponseWriter, r *http.Request, id string) {
	info, ok := h.service.GetSession(id)
	if !ok {
		writeJSONError(w, http.StatusNotFound, "stream not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"playback_url":  info.PlaybackURL,
		"hls_directory": info.HLSDirectory,
		"cdn_provider":  info.CDNProvider,
	})
}

func (h *Handler) handleStop(w http.ResponseWriter, r *http.Request, id string) {
	info, err := h.service.StopStream(id)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func parseStreamPath(path string) (id string, suffix string) {
	trimmed := strings.TrimPrefix(path, "/api/v1/streaming/")
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

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
