package commerce

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("/api/v1/commerce/streams/", handler.handleStreamCommerce)
}

func (h *Handler) handleStreamCommerce(w http.ResponseWriter, r *http.Request) {
	streamID, suffix := parseCommercePath(r.URL.Path)
	if streamID == "" {
		writeJSONError(w, http.StatusNotFound, "stream id required")
		return
	}

	switch {
	case suffix == "/cart" && r.Method == http.MethodPost:
		h.handleAddToCart(w, r, streamID)
	case suffix == "/cart" && r.Method == http.MethodGet:
		h.handleGetCart(w, r, streamID)
	case suffix == "/checkout" && r.Method == http.MethodPost:
		h.handleCheckout(w, r, streamID)
	case suffix == "/purchases" && r.Method == http.MethodGet:
		h.handleListPurchases(w, r, streamID)
	default:
		writeJSONError(w, http.StatusNotFound, "not found")
	}
}

type addToCartRequest struct {
	UserID    int64 `json:"user_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

func (h *Handler) handleAddToCart(w http.ResponseWriter, r *http.Request, streamID string) {
	var req addToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cart, err := h.service.AddToCart(streamID, req.UserID, req.ProductID, req.Quantity)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, cart)
}

type checkoutRequest struct {
	UserID        int64  `json:"user_id"`
	AffiliateID   *int64 `json:"affiliate_id,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request, streamID string) {
	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	receipt, err := h.service.PurchaseCart(streamID, req.UserID, req.AffiliateID, req.PaymentMethod)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, receipt)
}

func (h *Handler) handleGetCart(w http.ResponseWriter, r *http.Request, streamID string) {
	userID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil || userID == 0 {
		writeJSONError(w, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	cart, err := h.service.GetCart(streamID, userID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, cart)
}

func (h *Handler) handleListPurchases(w http.ResponseWriter, r *http.Request, streamID string) {
	receipts := h.service.ListPurchases(streamID)
	writeJSON(w, http.StatusOK, receipts)
}

func parseCommercePath(path string) (id string, suffix string) {
	trimmed := strings.TrimPrefix(path, "/api/v1/commerce/streams/")
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
