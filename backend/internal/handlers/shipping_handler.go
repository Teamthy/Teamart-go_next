package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/teamart/commerce-api/internal/shipping"
	"github.com/teamart/commerce-api/pkg/logger"
)

// ShippingHandler handles shipping-related HTTP requests
type ShippingHandler struct {
	service *shipping.Service
	logger  *logger.Logger
}

// NewShippingHandler creates a new shipping handler
func NewShippingHandler(svc *shipping.Service, log *logger.Logger) *ShippingHandler {
	return &ShippingHandler{service: svc, logger: log}
}

// ListShippingProfiles handles GET /api/shipping-profiles
// Query params: merchant_id
// Response: []ShippingProfile
func (h *ShippingHandler) ListShippingProfiles(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/shipping-profiles")

	merchantIDStr := r.URL.Query().Get("merchant_id")
	if merchantIDStr == "" {
		h.writeError(w, http.StatusBadRequest, "merchant_id is required", nil)
		return
	}

	merchantID, err := strconv.ParseInt(merchantIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid merchant_id", err)
		return
	}

	// TODO: Query shipping profiles from database
	profiles := []interface{}{}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"merchant_id": merchantID,
		"profiles":    profiles,
	})
}

// CreateShippingProfile handles POST /api/shipping-profiles
// Request body: {name, zones, rates, carriers, default_carrier}
// Response: ShippingProfile
func (h *ShippingHandler) CreateShippingProfile(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/shipping-profiles")

	var req struct {
		MerchantID     int64       `json:"merchant_id"`
		Name           string      `json:"name"`
		Description    string      `json:"description"`
		Zones          interface{} `json:"zones"`
		Rates          interface{} `json:"rates"`
		Carriers       interface{} `json:"carriers"`
		DefaultCarrier string      `json:"default_carrier"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if req.Name == "" {
		h.writeError(w, http.StatusBadRequest, "name is required", nil)
		return
	}

	// TODO: Create shipping profile in database

	response := map[string]interface{}{
		"id":             1,
		"merchant_id":    req.MerchantID,
		"name":           req.Name,
		"description":    req.Description,
		"is_default":     true,
		"is_active":      true,
		"created_at":     time.Now(),
		"zones":          req.Zones,
		"rates":          req.Rates,
		"carriers":       req.Carriers,
		"default_carrier": req.DefaultCarrier,
	}

	h.writeJSON(w, http.StatusCreated, response)
}

// GetShippingProfile handles GET /api/shipping-profiles/{profile_id}
// Response: ShippingProfile
func (h *ShippingHandler) GetShippingProfile(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/shipping-profiles/{profile_id}")

	profileIDStr := mux.Vars(r)["profile_id"]
	profileID, err := strconv.ParseInt(profileIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid profile_id", err)
		return
	}

	// TODO: Query shipping profile by ID

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id": profileID,
	})
}

// UpdateShippingProfile handles PUT /api/shipping-profiles/{profile_id}
// Request body: {name, zones, rates, carriers, default_carrier}
// Response: ShippingProfile
func (h *ShippingHandler) UpdateShippingProfile(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("PUT /api/shipping-profiles/{profile_id}")

	profileIDStr := mux.Vars(r)["profile_id"]
	profileID, err := strconv.ParseInt(profileIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid profile_id", err)
		return
	}

	var req struct {
		Name           string      `json:"name"`
		Zones          interface{} `json:"zones"`
		Rates          interface{} `json:"rates"`
		Carriers       interface{} `json:"carriers"`
		DefaultCarrier string      `json:"default_carrier"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// TODO: Update shipping profile in database

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":       profileID,
		"name":     req.Name,
		"zones":    req.Zones,
		"rates":    req.Rates,
		"carriers": req.Carriers,
	})
}

// CreateShipment handles POST /api/shipments
// Request body: {order_id, carrier, service, weight_grams, dimensions, shipping_address}
// Response: Shipment with tracking_number and label_url
func (h *ShippingHandler) CreateShipment(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/shipments")

	var req struct {
		OrderID          int64       `json:"order_id"`
		Carrier          string      `json:"carrier"`
		Service          string      `json:"carrier_service"`
		WeightGrams      int         `json:"weight_grams"`
		Dimensions       interface{} `json:"dimensions"`
		ShippingAddress  interface{} `json:"shipping_address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if req.OrderID == 0 {
		h.writeError(w, http.StatusBadRequest, "order_id is required", nil)
		return
	}

	// TODO: Create shipment and generate label via carrier API

	response := map[string]interface{}{
		"id":              1,
		"order_id":        req.OrderID,
		"shipment_number": "SHP-2026-001",
		"carrier":         req.Carrier,
		"carrier_service": req.Service,
		"tracking_number": "TRK123456789",
		"status":          "label_generated",
		"label_url":       "https://example.com/labels/label123.pdf",
		"created_at":      time.Now(),
	}

	h.writeJSON(w, http.StatusCreated, response)
}

// GetShipment handles GET /api/shipments/{shipment_id}
// Response: Shipment
func (h *ShippingHandler) GetShipment(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/shipments/{shipment_id}")

	shipmentIDStr := mux.Vars(r)["shipment_id"]
	shipmentID, err := strconv.ParseInt(shipmentIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid shipment_id", err)
		return
	}

	// TODO: Query shipment from database

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id": shipmentID,
	})
}

// GetTracking handles GET /api/track/{tracking_number}
// Response: TrackingInfo with timeline
func (h *ShippingHandler) GetTracking(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/track/{tracking_number}")

	trackingNumber := mux.Vars(r)["tracking_number"]
	if trackingNumber == "" {
		h.writeError(w, http.StatusBadRequest, "tracking_number is required", nil)
		return
	}

	// TODO: Query shipment by tracking number and get events

	response := map[string]interface{}{
		"tracking_number":   trackingNumber,
		"order_number":      "ORD-2026-001",
		"carrier":           "DHL",
		"status":            "in_transit",
		"expected_delivery": "2026-06-12",
		"events": []interface{}{
			map[string]interface{}{
				"status":      "picked_up",
				"location":    "Lagos Distribution Center",
				"event_time":  "2026-06-10T08:00:00Z",
				"description": "Package picked up",
			},
		},
	}

	h.writeJSON(w, http.StatusOK, response)
}

// StreamTracking handles GET /api/track/stream/{tracking_number}
// Server-Sent Events stream of tracking updates
func (h *ShippingHandler) StreamTracking(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/track/stream/{tracking_number} (SSE)")

	trackingNumber := mux.Vars(r)["tracking_number"]
	if trackingNumber == "" {
		h.writeError(w, http.StatusBadRequest, "tracking_number is required", nil)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// TODO: Stream tracking updates
	// This would be connected to a Redis pub/sub or similar

	flusher, ok := w.(http.Flusher)
	if !ok {
		h.writeError(w, http.StatusInternalServerError, "Streaming not supported", nil)
		return
	}

	// Send initial event
	w.Write([]byte("event: tracking_update\n"))
	w.Write([]byte("data: {\"status\":\"connected\",\"tracking_number\":\"" + trackingNumber + "\"}\n\n"))
	flusher.Flush()

	// Keep connection open
	<-r.Context().Done()
}

// Helper functions

func (h *ShippingHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *ShippingHandler) writeError(w http.ResponseWriter, statusCode int, message string, err error) {
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

// RegisterShippingRoutes registers shipping routes on the provided router
func RegisterShippingRoutes(mux Router, h *ShippingHandler) {
	mux.HandleFunc("GET /api/shipping-profiles", h.ListShippingProfiles)
	mux.HandleFunc("POST /api/shipping-profiles", h.CreateShippingProfile)
	mux.HandleFunc("GET /api/shipping-profiles/{profile_id}", h.GetShippingProfile)
	mux.HandleFunc("PUT /api/shipping-profiles/{profile_id}", h.UpdateShippingProfile)

	mux.HandleFunc("POST /api/shipments", h.CreateShipment)
	mux.HandleFunc("GET /api/shipments/{shipment_id}", h.GetShipment)
	mux.HandleFunc("GET /api/track/{tracking_number}", h.GetTracking)
	mux.HandleFunc("GET /api/track/stream/{tracking_number}", h.StreamTracking)
}
