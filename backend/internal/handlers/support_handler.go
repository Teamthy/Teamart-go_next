package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/teamart/commerce-api/internal/support"
	"github.com/teamart/commerce-api/pkg/logger"
)

// SupportHandler handles support tickets and messages-related HTTP requests
type SupportHandler struct {
	service *support.Service
	logger  *logger.Logger
}

// NewSupportHandler creates a new support handler
func NewSupportHandler(svc *support.Service, log *logger.Logger) *SupportHandler {
	return &SupportHandler{service: svc, logger: log}
}

// CreateTicket handles POST /api/support/tickets
// Request body: {merchant_id, customer_id, subject, description, category, priority}
// Response: SupportTicket
func (h *SupportHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/support/tickets")

	var req struct {
		MerchantID  int64  `json:"merchant_id"`
		CustomerID  int64  `json:"customer_id"`
		OrderID     *int64 `json:"order_id,omitempty"`
		Subject     string `json:"subject"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Priority    string `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if req.Subject == "" || req.Description == "" {
		h.writeError(w, http.StatusBadRequest, "subject and description are required", nil)
		return
	}

	// TODO: Create ticket in database

	response := map[string]interface{}{
		"id":           1,
		"ticket_number": "TKT-2026-001",
		"merchant_id":  req.MerchantID,
		"customer_id":  req.CustomerID,
		"subject":      req.Subject,
		"status":       "open",
		"priority":     req.Priority,
		"created_at":   time.Now(),
	}

	h.writeJSON(w, http.StatusCreated, response)
}

// GetTicket handles GET /api/support/tickets/{ticket_id}
// Response: SupportTicket with messages
func (h *SupportHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/support/tickets/{ticket_id}")

	ticketIDStr := mux.Vars(r)["ticket_id"]
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid ticket_id", err)
		return
	}

	// TODO: Query ticket from database

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id": ticketID,
	})
}

// ListTickets handles GET /api/support/tickets
// Query params: status, assigned_to, priority, created_after
// Response: []SupportTicket
func (h *SupportHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/support/tickets")

	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")

	// TODO: Query tickets from database with filters

	response := map[string]interface{}{
		"tickets": []interface{}{},
		"status":  status,
		"priority": priority,
	}

	h.writeJSON(w, http.StatusOK, response)
}

// UpdateTicketStatus handles PATCH /api/support/tickets/{ticket_id}/status
// Request body: {status: string}
// Response: SupportTicket
func (h *SupportHandler) UpdateTicketStatus(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("PATCH /api/support/tickets/{ticket_id}/status")

	ticketIDStr := mux.Vars(r)["ticket_id"]
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid ticket_id", err)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// TODO: Update ticket status in database

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":     ticketID,
		"status": req.Status,
	})
}

// AssignTicket handles POST /api/support/tickets/{ticket_id}/assign
// Request body: {assigned_to: int64}
// Response: SupportTicket
func (h *SupportHandler) AssignTicket(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/support/tickets/{ticket_id}/assign")

	ticketIDStr := mux.Vars(r)["ticket_id"]
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid ticket_id", err)
		return
	}

	var req struct {
		AssignedTo int64 `json:"assigned_to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// TODO: Assign ticket in database

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":          ticketID,
		"assigned_to": req.AssignedTo,
		"assigned_at": time.Now(),
	})
}

// AddMessage handles POST /api/support/tickets/{ticket_id}/messages
// Request body: {message: string, is_internal: bool}
// Response: SupportMessage
func (h *SupportHandler) AddMessage(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/support/tickets/{ticket_id}/messages")

	ticketIDStr := mux.Vars(r)["ticket_id"]
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid ticket_id", err)
		return
	}

	var req struct {
		Message    string `json:"message"`
		IsInternal bool   `json:"is_internal"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if req.Message == "" {
		h.writeError(w, http.StatusBadRequest, "message is required", nil)
		return
	}

	// TODO: Create message in database

	response := map[string]interface{}{
		"id":          1,
		"ticket_id":   ticketID,
		"message":     req.Message,
		"is_internal": req.IsInternal,
		"created_at":  time.Now(),
	}

	h.writeJSON(w, http.StatusCreated, response)
}

// GetMessages handles GET /api/support/tickets/{ticket_id}/messages
// Response: []SupportMessage
func (h *SupportHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/support/tickets/{ticket_id}/messages")

	ticketIDStr := mux.Vars(r)["ticket_id"]
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid ticket_id", err)
		return
	}

	// TODO: Query messages from database

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"ticket_id": ticketID,
		"messages":  []interface{}{},
	})
}

// Helper functions

func (h *SupportHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *SupportHandler) writeError(w http.ResponseWriter, statusCode int, message string, err error) {
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

// RegisterSupportRoutes registers support routes
func RegisterSupportRoutes(mux Router, h *SupportHandler) {
	mux.HandleFunc("POST /api/support/tickets", h.CreateTicket)
	mux.HandleFunc("GET /api/support/tickets", h.ListTickets)
	mux.HandleFunc("GET /api/support/tickets/{ticket_id}", h.GetTicket)
	mux.HandleFunc("PATCH /api/support/tickets/{ticket_id}/status", h.UpdateTicketStatus)
	mux.HandleFunc("POST /api/support/tickets/{ticket_id}/assign", h.AssignTicket)
	mux.HandleFunc("POST /api/support/tickets/{ticket_id}/messages", h.AddMessage)
	mux.HandleFunc("GET /api/support/tickets/{ticket_id}/messages", h.GetMessages)
}
