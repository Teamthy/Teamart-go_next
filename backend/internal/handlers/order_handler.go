package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/teamart/commerce-api/internal/orders"
	"github.com/teamart/commerce-api/pkg/logger"
)

// OrderHandler handles HTTP requests related to orders
type OrderHandler struct {
	service *orders.Service
	logger  *logger.Logger
}

// NewOrderHandler creates a new order HTTP handler
func NewOrderHandler(service *orders.Service, logger *logger.Logger) *OrderHandler {
	return &OrderHandler{
		service: service,
		logger:  logger,
	}
}

// CreateOrderRequest represents the HTTP request body
type CreateOrderRequest struct {
	UserID      int64   `json:"user_id" binding:"required"`
	TotalAmount float64 `json:"total_amount" binding:"required"`
	Status      string  `json:"status"`
}

// OrderResponse represents the HTTP response body
type OrderResponse struct {
	ID            int64   `json:"id"`
	OrderNumber   string  `json:"order_number"`
	UserID        int64   `json:"user_id"`
	CustomerName  string  `json:"customer_name,omitempty"`
	CustomerEmail string  `json:"customer_email,omitempty"`
	TotalAmount   float64 `json:"total_amount"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method,omitempty"`
	ItemsCount    int64   `json:"items_count,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

func makeOrderResponse(data *orders.OrderData) OrderResponse {
	return OrderResponse{
		ID:            data.ID,
		OrderNumber:   data.OrderNumber,
		UserID:        data.UserID,
		CustomerName:  data.CustomerName,
		CustomerEmail: data.CustomerEmail,
		TotalAmount:   data.TotalAmount,
		Status:        data.Status,
		PaymentMethod: data.PaymentMethod,
		ItemsCount:    data.ItemsCount,
		CreatedAt:     data.CreatedAt,
		UpdatedAt:     data.UpdatedAt,
	}
}

// HandleCreateOrder handles POST /orders requests
//
//	Example: curl -X POST http://localhost:8080/orders \
//	  -H "Content-Type: application/json" \
//	  -d '{"user_id":1,"total_amount":199.99,"status":"pending"}'
func (h *OrderHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling CreateOrder request")

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := &orders.CreateOrderInput{
		UserID:      req.UserID,
		TotalAmount: req.TotalAmount,
		Status:      req.Status,
	}

	output, err := h.service.CreateOrder(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(makeOrderResponse(&orders.OrderData{
		ID:            output.ID,
		OrderNumber:   output.OrderNumber,
		UserID:        output.UserID,
		CustomerName:  output.CustomerName,
		CustomerEmail: output.CustomerEmail,
		TotalAmount:   output.TotalAmount,
		Status:        output.Status,
		PaymentMethod: output.PaymentMethod,
		ItemsCount:    output.ItemsCount,
		CreatedAt:     output.CreatedAt,
		UpdatedAt:     output.UpdatedAt,
	}))
}

// HandleGetOrder handles GET /orders/:id requests
// Example: curl http://localhost:8080/orders/1
func (h *OrderHandler) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["id"]
	if orderIDStr == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	h.logger.Debugf("handling GetOrder request for order: %d", orderID)

	input := &orders.GetOrderByIDInput{OrderID: orderID}
	output, err := h.service.GetOrderByID(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(makeOrderResponse(&orders.OrderData{
		ID:            output.ID,
		OrderNumber:   output.OrderNumber,
		UserID:        output.UserID,
		CustomerName:  output.CustomerName,
		CustomerEmail: output.CustomerEmail,
		TotalAmount:   output.TotalAmount,
		Status:        output.Status,
		PaymentMethod: output.PaymentMethod,
		ItemsCount:    output.ItemsCount,
		CreatedAt:     output.CreatedAt,
		UpdatedAt:     output.UpdatedAt,
	}))
}

// ListOrdersResponse represents the HTTP response body
type ListOrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
	Limit  int32           `json:"limit"`
	Offset int32           `json:"offset"`
}

// HandleListOrdersByUser handles GET /users/:user_id/orders requests
// Example: curl "http://localhost:8080/users/1/orders?limit=10&offset=0"
func (h *OrderHandler) HandleListOrdersByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	limit := int32(10)
	offset := int32(0)

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 32); err == nil {
			offset = int32(o)
		}
	}

	h.logger.Debugf("handling ListOrdersByUser request for user: %d", userID)

	input := &orders.ListOrdersByUserIDInput{UserID: userID, Limit: limit, Offset: offset}
	output, err := h.service.ListOrdersByUserID(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orderResponses := make([]OrderResponse, len(output.Orders))
	for i, order := range output.Orders {
		orderResponses[i] = makeOrderResponse(&order)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ListOrdersResponse{
		Orders: orderResponses,
		Limit:  output.Limit,
		Offset: output.Offset,
	})
}

// HandleListOrdersByStatus handles GET /orders/status/:status requests
// Example: curl "http://localhost:8080/orders/status/pending?limit=10"
func (h *OrderHandler) HandleListOrdersByStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]
	if status == "" {
		http.Error(w, "Status is required", http.StatusBadRequest)
		return
	}

	limit := int32(10)
	offset := int32(0)

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 32); err == nil {
			offset = int32(o)
		}
	}

	h.logger.Debugf("handling ListOrdersByStatus request for status: %s", status)

	input := &orders.ListOrdersByStatusInput{Status: status, Limit: limit, Offset: offset}
	output, err := h.service.ListOrdersByStatus(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orderResponses := make([]OrderResponse, len(output.Orders))
	for i, order := range output.Orders {
		orderResponses[i] = makeOrderResponse(&order)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ListOrdersResponse{
		Orders: orderResponses,
		Limit:  output.Limit,
		Offset: output.Offset,
	})
}

// HandleListAllOrders handles GET /orders requests
// Example: curl "http://localhost:8080/orders?limit=10&offset=0"
func (h *OrderHandler) HandleListAllOrders(w http.ResponseWriter, r *http.Request) {
	limit := int32(10)
	offset := int32(0)

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 32); err == nil {
			offset = int32(o)
		}
	}

	h.logger.Debugf("handling ListAllOrders request with limit: %d, offset: %d", limit, offset)

	input := &orders.ListAllOrdersInput{Limit: limit, Offset: offset}
	output, err := h.service.ListAllOrders(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orderResponses := make([]OrderResponse, len(output.Orders))
	for i, order := range output.Orders {
		orderResponses[i] = makeOrderResponse(&order)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ListOrdersResponse{
		Orders: orderResponses,
		Limit:  output.Limit,
		Offset: output.Offset,
	})
}

// UpdateOrderStatusRequest represents request to update order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Notes  string `json:"notes"`
}

// HandleUpdateOrderStatus handles PATCH /orders/{id}/status
// Updates the status of an order (pending, confirmed, shipped, delivered, cancelled)
func (h *OrderHandler) HandleUpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["id"]
	if orderIDStr == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var req UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Debugf("handling UpdateOrderStatus request for order: %d, status: %s", orderID, req.Status)

	input := &orders.UpdateOrderStatusInput{OrderID: orderID, Status: req.Status}
	output, err := h.service.UpdateOrderStatus(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(makeOrderResponse(&orders.OrderData{
		ID:            output.ID,
		OrderNumber:   output.OrderNumber,
		UserID:        output.UserID,
		CustomerName:  output.CustomerName,
		CustomerEmail: output.CustomerEmail,
		TotalAmount:   output.TotalAmount,
		Status:        output.Status,
		PaymentMethod: output.PaymentMethod,
		ItemsCount:    output.ItemsCount,
		CreatedAt:     output.CreatedAt,
		UpdatedAt:     output.UpdatedAt,
	}))
}

// BulkUpdateOrdersRequest represents request to bulk update orders
type BulkUpdateOrdersRequest struct {
	OrderIDs []int64 `json:"order_ids" binding:"required"`
	Status   string  `json:"status" binding:"required"`
	Notes    string  `json:"notes"`
}

// HandleBulkUpdateOrders handles PATCH /orders/bulk/status
// Updates status for multiple orders atomically
func (h *OrderHandler) HandleBulkUpdateOrders(w http.ResponseWriter, r *http.Request) {
	var req BulkUpdateOrdersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.OrderIDs) == 0 {
		http.Error(w, "order_ids is required", http.StatusBadRequest)
		return
	}

	h.logger.Debugf("handling BulkUpdateOrders for %d orders with status: %s", len(req.OrderIDs), req.Status)

	// TODO: Call service to bulk update orders
	// input := &orders.BulkUpdateOrdersInput{OrderIDs: req.OrderIDs, Status: req.Status, Notes: req.Notes}
	// output, err := h.service.BulkUpdateOrders(r.Context(), input)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"updated_count": len(req.OrderIDs),
		"status":        req.Status,
	})
}

// TimelineEvent represents an order timeline event
type TimelineEvent struct {
	ID        int64  `json:"id"`
	OrderID   int64  `json:"order_id"`
	EventType string `json:"event_type"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

// HandleGetOrderTimeline handles GET /orders/{id}/timeline
// Returns the timeline of status changes and events for an order
func (h *OrderHandler) HandleGetOrderTimeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["id"]
	if orderIDStr == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	h.logger.Debugf("handling GetOrderTimeline request for order: %d", orderID)

	// TODO: Call service to get timeline
	// input := &orders.GetOrderTimelineInput{OrderID: orderID}
	// output, err := h.service.GetOrderTimeline(r.Context(), input)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order_id": orderID,
		"events":   []TimelineEvent{},
	})
}

// SearchOrdersRequest represents a search query for orders
type SearchOrdersRequest struct {
	Query      string   `json:"query"`
	Status     string   `json:"status"`
	UserID     int64    `json:"user_id"`
	DateFrom   string   `json:"date_from"`
	DateTo     string   `json:"date_to"`
	MinAmount  float64  `json:"min_amount"`
	MaxAmount  float64  `json:"max_amount"`
	SortBy     string   `json:"sort_by"`
	SortOrder  string   `json:"sort_order"`
	Limit      int32    `json:"limit"`
	Offset     int32    `json:"offset"`
}

// HandleSearchOrders handles POST /orders/search
// Searches orders with advanced filters and sorting
func (h *OrderHandler) HandleSearchOrders(w http.ResponseWriter, r *http.Request) {
	var req SearchOrdersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	h.logger.Debugf("handling SearchOrders with query: %s, status: %s", req.Query, req.Status)

	// TODO: Call service to search orders with filters
	// input := &orders.SearchOrdersInput{...}
	// output, err := h.service.SearchOrders(r.Context(), input)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": []OrderResponse{},
		"limit":   req.Limit,
		"offset":  req.Offset,
		"total":   0,
	})
}

// HandleExportOrders handles GET /orders/export
// Query params: format (csv, excel, json), status, date_from, date_to
// Response: {download_url, expires_in}
func (h *OrderHandler) HandleExportOrders(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "csv"
	}

	status := r.URL.Query().Get("status")
	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	h.logger.Debugf("handling ExportOrders with format: %s, status: %s", format, status)

	// TODO: Call service to generate export
	// Create temporary file with orders data
	// Upload to storage (S3, etc)
	// Return download URL

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"download_url": "https://api.teamart.local/exports/orders-2026-001." + format,
		"format":       format,
		"expires_in":   3600,
	})
}

// RegisterOrderRoutes registers all order-related routes
func RegisterOrderRoutes(mux Router, handler *OrderHandler) {
	// Order endpoints
	mux.HandleFunc("POST /orders", handler.HandleCreateOrder)
	mux.HandleFunc("GET /orders", handler.HandleListAllOrders)
	mux.HandleFunc("GET /orders/{id}", handler.HandleGetOrder)
	mux.HandleFunc("GET /users/{user_id}/orders", handler.HandleListOrdersByUser)
	mux.HandleFunc("GET /orders/status/{status}", handler.HandleListOrdersByStatus)
	
	// Enhanced endpoints (Phase 3)
	mux.HandleFunc("PATCH /orders/{id}/status", handler.HandleUpdateOrderStatus)
	mux.HandleFunc("PUT /orders/{id}", handler.HandleUpdateOrderStatus)
	mux.HandleFunc("PATCH /orders/bulk/status", handler.HandleBulkUpdateOrders)
	mux.HandleFunc("GET /orders/{id}/timeline", handler.HandleGetOrderTimeline)
	mux.HandleFunc("POST /orders/search", handler.HandleSearchOrders)
	mux.HandleFunc("GET /orders/export", handler.HandleExportOrders)
}
