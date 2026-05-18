package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	ID          int64   `json:"id"`
	UserID      int64   `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// HandleCreateOrder handles POST /orders requests
// Example: curl -X POST http://localhost:8080/orders \
//   -H "Content-Type: application/json" \
//   -d '{"user_id":1,"total_amount":199.99,"status":"pending"}'
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
	json.NewEncoder(w).Encode(OrderResponse{
		ID:          output.ID,
		UserID:      output.UserID,
		TotalAmount: output.TotalAmount,
		Status:      output.Status,
		CreatedAt:   output.CreatedAt,
		UpdatedAt:   output.UpdatedAt,
	})
}

// HandleGetOrder handles GET /orders/:id requests
// Example: curl http://localhost:8080/orders/1
func (h *OrderHandler) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.PathValue("id")
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
	json.NewEncoder(w).Encode(OrderResponse{
		ID:          output.ID,
		UserID:      output.UserID,
		TotalAmount: output.TotalAmount,
		Status:      output.Status,
		CreatedAt:   output.CreatedAt,
		UpdatedAt:   output.UpdatedAt,
	})
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
	userIDStr := r.PathValue("user_id")
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
		orderResponses[i] = OrderResponse{
			ID:          order.ID,
			UserID:      order.UserID,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		}
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
	status := r.PathValue("status")
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
		orderResponses[i] = OrderResponse{
			ID:          order.ID,
			UserID:      order.UserID,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		}
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
		orderResponses[i] = OrderResponse{
			ID:          order.ID,
			UserID:      order.UserID,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ListOrdersResponse{
		Orders: orderResponses,
		Limit:  output.Limit,
		Offset: output.Offset,
	})
}

// RegisterOrderRoutes registers all order-related routes
func RegisterOrderRoutes(mux *http.ServeMux, handler *OrderHandler) {
	// Order endpoints
	mux.HandleFunc("POST /orders", handler.HandleCreateOrder)
	mux.HandleFunc("GET /orders", handler.HandleListAllOrders)
	mux.HandleFunc("GET /orders/{id}", handler.HandleGetOrder)
	mux.HandleFunc("GET /users/{user_id}/orders", handler.HandleListOrdersByUser)
	mux.HandleFunc("GET /orders/status/{status}", handler.HandleListOrdersByStatus)
}
