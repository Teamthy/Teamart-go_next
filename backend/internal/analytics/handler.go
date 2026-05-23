package analytics

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

type Handler struct {
	service *AnalyticsService
	logger  *logger.Logger
}

func NewHandler(service *AnalyticsService, logger *logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func RegisterAnalyticsRoutes(router interface {
	HandleFunc(string, func(http.ResponseWriter, *http.Request))
}, handler *Handler) {
	if router == nil || handler == nil {
		return
	}

	router.HandleFunc("POST /api/v1/analytics/events", handler.handleIngestEvent)
	router.HandleFunc("GET /api/v1/analytics/metrics", handler.handleGetMetrics)
	router.HandleFunc("GET /api/v1/analytics/metrics/creator", handler.handleGetCreatorMetrics)
	router.HandleFunc("GET /api/v1/analytics/metrics/marketplace", handler.handleGetMarketplaceMetrics)
}

func (h *Handler) handleIngestEvent(w http.ResponseWriter, r *http.Request) {
	var input AnalyticsEventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Errorf("failed to decode analytics event request: %v", err)
		h.writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.EventType == "" {
		h.writeJSONError(w, http.StatusBadRequest, "event_type is required")
		return
	}

	timestamp := time.Now()
	if input.Timestamp != "" {
		parsed, err := time.Parse(time.RFC3339, input.Timestamp)
		if err == nil {
			timestamp = parsed
		}
	}

	event := &EventRecord{
		EventType: input.EventType,
		UserID:    input.UserID,
		SessionID: input.SessionID,
		Timestamp: timestamp,
		Data:      input.Data,
	}

	if err := h.service.IngestEvent(event); err != nil {
		h.logger.Errorf("failed to ingest analytics event: %v", err)
		h.writeJSONError(w, http.StatusInternalServerError, "failed to ingest analytics event")
		return
	}

	h.writeJSON(w, http.StatusAccepted, map[string]string{"status": "ingested"})
}

func (h *Handler) handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, h.service.GetMetrics())
}

func (h *Handler) handleGetCreatorMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := h.service.GetMetrics().Creator
	h.writeJSON(w, http.StatusOK, metrics)
}

func (h *Handler) handleGetMarketplaceMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := h.service.GetMetrics().Marketplace
	h.writeJSON(w, http.StatusOK, metrics)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *Handler) writeJSONError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
