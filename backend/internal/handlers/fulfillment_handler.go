package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/teamart/commerce-api/internal/fulfillment"
	"github.com/teamart/commerce-api/pkg/logger"
)

// FulfillmentHandler handles fulfillment and inventory-related HTTP requests
type FulfillmentHandler struct {
	service *fulfillment.Service
	logger  *logger.Logger
}

// NewFulfillmentHandler creates a new fulfillment handler
func NewFulfillmentHandler(svc *fulfillment.Service, log *logger.Logger) *FulfillmentHandler {
	return &FulfillmentHandler{service: svc, logger: log}
}

// GetDashboard handles GET /api/fulfillment/dashboard
// Query params: merchant_id
// Response: FulfillmentDashboard
func (h *FulfillmentHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/fulfillment/dashboard")

	merchantIDStr := r.URL.Query().Get("merchant_id")
	if merchantIDStr == "" {
		h.writeError(w, http.StatusBadRequest, "merchant_id is required", nil)
		return
	}

	// TODO: Query dashboard metrics from database

	response := map[string]interface{}{
		"merchant_id": merchantIDStr,
		"pending_orders": 0,
		"pending_shipments": 0,
		"returns_pending": 0,
		"inventory_low": 0,
		"recent_orders": []interface{}{},
		"metrics": map[string]interface{}{
			"avg_fulfillment_time": "24h",
			"fulfillment_rate": 0.95,
			"return_rate": 0.02,
		},
	}

	h.writeJSON(w, http.StatusOK, response)
}

// CreateBatch handles POST /api/fulfillment/batch
// Request body: {merchant_id, order_ids, ship_by_date}
// Response: FulfillmentBatch
func (h *FulfillmentHandler) CreateBatch(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/fulfillment/batch")

	var req struct {
		MerchantID  int64   `json:"merchant_id"`
		OrderIDs    []int64 `json:"order_ids"`
		ShipByDate  string  `json:"ship_by_date"`
		ShippingLabelURL string `json:"shipping_label_url,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if len(req.OrderIDs) == 0 {
		h.writeError(w, http.StatusBadRequest, "order_ids is required", nil)
		return
	}

	// TODO: Create batch in database

	response := map[string]interface{}{
		"id":              1,
		"batch_number":    "BATCH-2026-001",
		"merchant_id":     req.MerchantID,
		"order_ids":       req.OrderIDs,
		"status":          "created",
		"items_count":     len(req.OrderIDs),
		"ship_by_date":    req.ShipByDate,
		"created_at":      time.Now(),
	}

	h.writeJSON(w, http.StatusCreated, response)
}

// GetBatch handles GET /api/fulfillment/batch/{batch_id}
// Response: FulfillmentBatch with items
func (h *FulfillmentHandler) GetBatch(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/fulfillment/batch/{batch_id}")

	batchIDStr := mux.Vars(r)["batch_id"]
	batchID, err := strconv.ParseInt(batchIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid batch_id", err)
		return
	}

	// TODO: Query batch from database

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id": batchID,
	})
}

// ListBatches handles GET /api/fulfillment/batch
// Query params: merchant_id, status
// Response: []FulfillmentBatch
func (h *FulfillmentHandler) ListBatches(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/fulfillment/batch")

	merchantIDStr := r.URL.Query().Get("merchant_id")
	status := r.URL.Query().Get("status")

	// TODO: Query batches from database

	response := map[string]interface{}{
		"merchant_id": merchantIDStr,
		"status":      status,
		"batches":     []interface{}{},
	}

	h.writeJSON(w, http.StatusOK, response)
}

// ReserveStock handles POST /api/fulfillment/reserve
// Request body: {order_id, items}
// Response: Reservation
func (h *FulfillmentHandler) ReserveStock(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/fulfillment/reserve")

	var req struct {
		OrderID int64 `json:"order_id"`
		Items   []struct {
			ProductID int64 `json:"product_id"`
			Quantity  int64 `json:"quantity"`
		} `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if req.OrderID == 0 || len(req.Items) == 0 {
		h.writeError(w, http.StatusBadRequest, "order_id and items are required", nil)
		return
	}

	// TODO: Create reservation in database

	response := map[string]interface{}{
		"id":           1,
		"order_id":     req.OrderID,
		"status":       "reserved",
		"items_count":  len(req.Items),
		"reserved_at":  time.Now(),
		"expires_at":   time.Now().Add(30 * 24 * time.Hour),
	}

	h.writeJSON(w, http.StatusCreated, response)
}

// UpdateInventory handles POST /api/fulfillment/inventory
// Request body: {product_id, quantity_delta, reason}
// Response: InventoryUpdate
func (h *FulfillmentHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/fulfillment/inventory")

	var req struct {
		ProductID     int64  `json:"product_id"`
		QuantityDelta int64  `json:"quantity_delta"`
		Reason        string `json:"reason"`
		Reference     string `json:"reference,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if req.ProductID == 0 {
		h.writeError(w, http.StatusBadRequest, "product_id is required", nil)
		return
	}

	// TODO: Update inventory in database

	response := map[string]interface{}{
		"product_id":     req.ProductID,
		"quantity_delta": req.QuantityDelta,
		"reason":         req.Reason,
		"updated_at":     time.Now(),
	}

	h.writeJSON(w, http.StatusOK, response)
}

// GetInventory handles GET /api/fulfillment/inventory/{product_id}
// Response: InventoryStatus
func (h *FulfillmentHandler) GetInventory(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/fulfillment/inventory/{product_id}")

	productIDStr := mux.Vars(r)["product_id"]
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid product_id", err)
		return
	}

	// TODO: Query inventory from database

	response := map[string]interface{}{
		"product_id":     productID,
		"quantity":       0,
		"reserved":       0,
		"available":      0,
		"low_threshold":  10,
		"status":         "in_stock",
	}

	h.writeJSON(w, http.StatusOK, response)
}

// Helper functions

func (h *FulfillmentHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *FulfillmentHandler) writeError(w http.ResponseWriter, statusCode int, message string, err error) {
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

// RegisterFulfillmentRoutes registers fulfillment routes
func RegisterFulfillmentRoutes(mux Router, h *FulfillmentHandler) {
	mux.HandleFunc("GET /api/fulfillment/dashboard", h.GetDashboard)
	mux.HandleFunc("POST /api/fulfillment/batch", h.CreateBatch)
	mux.HandleFunc("GET /api/fulfillment/batch/{batch_id}", h.GetBatch)
	mux.HandleFunc("GET /api/fulfillment/batch", h.ListBatches)
	mux.HandleFunc("POST /api/fulfillment/reserve", h.ReserveStock)
	mux.HandleFunc("POST /api/fulfillment/inventory", h.UpdateInventory)
	mux.HandleFunc("GET /api/fulfillment/inventory/{product_id}", h.GetInventory)
}
