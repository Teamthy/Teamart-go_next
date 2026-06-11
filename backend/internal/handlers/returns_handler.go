package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/teamart/commerce-api/internal/returns"
	"github.com/teamart/commerce-api/pkg/logger"
)

// ReturnsHandler handles returns and refunds-related HTTP requests
type ReturnsHandler struct {
	service *returns.Service
	logger  *logger.Logger
}

// NewReturnsHandler creates a new returns handler
func NewReturnsHandler(svc *returns.Service, log *logger.Logger) *ReturnsHandler {
	return &ReturnsHandler{service: svc, logger: log}
}

// CreateReturn handles POST /api/returns
// Request body: {order_id, return_type, reason, description, items}
// Response: Return with return_number
func (h *ReturnsHandler) CreateReturn(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/returns")

	var req struct {
		OrderID     int64       `json:"order_id"`
		ReturnType  string      `json:"return_type"` // return, refund, replacement
		Reason      string      `json:"reason"`      // damaged, wrong_item, etc.
		Description string      `json:"description"`
		Items       interface{} `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if req.OrderID == 0 {
		h.writeError(w, http.StatusBadRequest, "order_id is required", nil)
		return
	}

	// TODO: Create return in database

	response := map[string]interface{}{
		"id":            1,
		"order_id":      req.OrderID,
		"return_number": "RET-2026-001",
		"return_type":   req.ReturnType,
		"reason":        req.Reason,
		"status":        "requested",
		"requested_at":  time.Now(),
	}

	h.writeJSON(w, http.StatusCreated, response)
}

// GetReturn handles GET /api/returns/{return_id}
// Response: Return with timeline
func (h *ReturnsHandler) GetReturn(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/returns/{return_id}")

	returnIDStr := mux.Vars(r)["return_id"]
	returnID, err := strconv.ParseInt(returnIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid return_id", err)
		return
	}

	// TODO: Query return from database

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id": returnID,
	})
}

// ListReturns handles GET /api/returns
// Query params: order_id, status, customer_id
// Response: []Return
func (h *ReturnsHandler) ListReturns(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/returns")

	status := r.URL.Query().Get("status")
	orderIDStr := r.URL.Query().Get("order_id")

	// TODO: Query returns from database with filters

	response := map[string]interface{}{
		"returns": []interface{}{},
		"status":  status,
	}
	if orderIDStr != "" {
		response["order_id"] = orderIDStr
	}

	h.writeJSON(w, http.StatusOK, response)
}

// ApproveReturn handles POST /api/returns/{return_id}/approve
// Request body: {notes: string}
// Response: Return
func (h *ReturnsHandler) ApproveReturn(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/returns/{return_id}/approve")

	returnIDStr := mux.Vars(r)["return_id"]
	returnID, err := strconv.ParseInt(returnIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid return_id", err)
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// TODO: Approve return in database

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":     returnID,
		"status": "approved",
		"notes":  req.Notes,
	})
}

// RejectReturn handles POST /api/returns/{return_id}/reject
// Request body: {reason: string}
// Response: Return
func (h *ReturnsHandler) RejectReturn(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/returns/{return_id}/reject")

	returnIDStr := mux.Vars(r)["return_id"]
	returnID, err := strconv.ParseInt(returnIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid return_id", err)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// TODO: Reject return in database

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":     returnID,
		"status": "rejected",
		"reason": req.Reason,
	})
}

// GetReturnTimeline handles GET /api/returns/{return_id}/timeline
// Response: []TimelineEvent
func (h *ReturnsHandler) GetReturnTimeline(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/returns/{return_id}/timeline")

	returnIDStr := mux.Vars(r)["return_id"]
	returnID, err := strconv.ParseInt(returnIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid return_id", err)
		return
	}

	// TODO: Query return timeline from database

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"return_id": returnID,
		"events":    []interface{}{},
	})
}

// Helper functions

func (h *ReturnsHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *ReturnsHandler) writeError(w http.ResponseWriter, statusCode int, message string, err error) {
	if err != nil {
		h.logger.Errorf("%s: %v", message, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   message,
		"details": errorString(err),
	})
}

// RegisterReturnsRoutes registers returns routes
func RegisterReturnsRoutes(mux Router, h *ReturnsHandler) {
	mux.HandleFunc("POST /api/returns", h.CreateReturn)
	mux.HandleFunc("GET /api/returns", h.ListReturns)
	mux.HandleFunc("GET /api/returns/{return_id}", h.GetReturn)
	mux.HandleFunc("POST /api/returns/{return_id}/approve", h.ApproveReturn)
	mux.HandleFunc("POST /api/returns/{return_id}/reject", h.RejectReturn)
	mux.HandleFunc("GET /api/returns/{return_id}/timeline", h.GetReturnTimeline)
}
