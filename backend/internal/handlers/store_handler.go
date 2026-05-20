package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/teamart/commerce-api/internal/auth"
	"github.com/teamart/commerce-api/internal/middleware"
	"github.com/teamart/commerce-api/pkg/logger"
)

// StoreHandler handles store management operations
type StoreHandler struct {
	logger *logger.Logger
}

// NewStoreHandler creates a new store handler
func NewStoreHandler(logger *logger.Logger) *StoreHandler {
	return &StoreHandler{
		logger: logger,
	}
}

// CreateStoreRequest represents the request body for store creation
type CreateStoreRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	BannerURL   string `json:"banner_url,omitempty"`
}

// CreateStoreResponse represents the response for store creation
type CreateStoreResponse struct {
	StoreID   int64  `json:"store_id"`
	Name      string `json:"name"`
	OwnerID   int64  `json:"owner_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// GetStoreResponse represents the response for store retrieval
type GetStoreResponse struct {
	StoreID     int64  `json:"store_id"`
	OwnerID     int64  `json:"owner_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	BannerURL   string `json:"banner_url"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// UpdateStoreRequest represents the request body for store update
type UpdateStoreRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
	BannerURL   string `json:"banner_url,omitempty"`
}

// UpdateStoreResponse represents the response for store update
type UpdateStoreResponse struct {
	StoreID   int64  `json:"store_id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
	Message   string `json:"message"`
}

// HandleCreateStore handles POST /stores request to create a new store
//
// This endpoint creates a new store for the authenticated user:
// - Validates user has merchant role
// - Creates store record
// - Publishes store.created event
// - Returns store details
//
//	Example: curl -X POST http://localhost:8080/stores \
//	  -H "Content-Type: application/json" \
//	  -H "Authorization: Bearer <token>" \
//	  -d '{"name":"My Store","description":"My awesome store","category":"electronics"}'
func (h *StoreHandler) HandleCreateStore(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling create store request")

	user, ok := r.Context().Value("user").(*auth.CustomClaims)
	if !ok || user == nil {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		h.respondError(w, "Store name is required", http.StatusBadRequest)
		return
	}

	// TODO: In production, persist to database
	// For now, return mock response
	storeID := int64(1000 + user.UserID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateStoreResponse{
		StoreID:   storeID,
		Name:      req.Name,
		OwnerID:   user.UserID,
		Status:    "active",
		CreatedAt: time.Now().Format(time.RFC3339),
	})

	h.logger.Infof("store created: %d by user %d", storeID, user.UserID)
}

// HandleGetStore handles GET /stores/{store_id} request
//
// This endpoint retrieves store details:
// - Validates store ownership (user can only see own stores)
// - Returns store information
//
//	Example: curl -X GET http://localhost:8080/stores/1000 \
//	  -H "Authorization: Bearer <token>"
func (h *StoreHandler) HandleGetStore(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling get store request")

	user, ok := r.Context().Value("user").(*auth.CustomClaims)
	if !ok || user == nil {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	storeIDStr := r.PathValue("store_id")
	if storeIDStr == "" {
		h.respondError(w, "Missing store_id parameter", http.StatusBadRequest)
		return
	}

	storeID, err := strconv.ParseInt(storeIDStr, 10, 64)
	if err != nil {
		h.respondError(w, "Invalid store_id", http.StatusBadRequest)
		return
	}

	// Get tenant context (user owns their own stores)
	tenantID, ok := middleware.GetTenantID(r)
	if !ok || tenantID != user.UserID {
		h.respondError(w, "Forbidden: you cannot access this store", http.StatusForbidden)
		return
	}

	// TODO: Query store from database
	// For now, return mock response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetStoreResponse{
		StoreID:     storeID,
		OwnerID:     user.UserID,
		Name:        "Sample Store",
		Description: "A sample store",
		Category:    "electronics",
		Status:      "active",
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	})
}

// HandleUpdateStore handles PUT /stores/{store_id} request
//
// This endpoint updates store details:
// - Validates store ownership
// - Updates store information
// - Publishes store.updated event
//
//	Example: curl -X PUT http://localhost:8080/stores/1000 \
//	  -H "Content-Type: application/json" \
//	  -H "Authorization: Bearer <token>" \
//	  -d '{"name":"Updated Store Name"}'
func (h *StoreHandler) HandleUpdateStore(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling update store request")

	user, ok := r.Context().Value("user").(*auth.CustomClaims)
	if !ok || user == nil {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	storeIDStr := r.PathValue("store_id")
	if storeIDStr == "" {
		h.respondError(w, "Missing store_id parameter", http.StatusBadRequest)
		return
	}

	storeID, err := strconv.ParseInt(storeIDStr, 10, 64)
	if err != nil {
		h.respondError(w, "Invalid store_id", http.StatusBadRequest)
		return
	}

	var req UpdateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify ownership (user owns the store)
	ownerID, ok := middleware.GetStoreOwnerID(r)
	if !ok || ownerID != user.UserID {
		h.logger.Warnf("unauthorized store update: user %d attempting to update store %d", user.UserID, storeID)
		h.respondError(w, "Forbidden: you cannot update this store", http.StatusForbidden)
		return
	}

	// TODO: Update store in database
	// Update the provided fields

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UpdateStoreResponse{
		StoreID:   storeID,
		Name:      req.Name,
		Status:    "active",
		UpdatedAt: time.Now().Format(time.RFC3339),
		Message:   "Store updated successfully",
	})

	h.logger.Infof("store updated: %d by user %d", storeID, user.UserID)
}

// HandleListStores handles GET /stores request to list user's stores
func (h *StoreHandler) HandleListStores(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling list stores request")

	user, ok := r.Context().Value("user").(*auth.CustomClaims)
	if !ok || user == nil {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// TODO: Query stores from database for this user
	// For now, return empty list
	stores := []GetStoreResponse{}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"stores": stores,
		"total":  0,
	})

	h.logger.Infof("listed stores for user %d", user.UserID)
}

// ===== Helper Functions =====

func (h *StoreHandler) respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   http.StatusText(statusCode),
		"message": message,
		"code":    statusCode,
	})
}

// RegisterStoreRoutes registers all store management routes
func RegisterStoreRoutes(mux *http.ServeMux, handler *StoreHandler) {
	mux.HandleFunc("GET /stores", handler.HandleListStores)
	mux.HandleFunc("POST /stores", handler.HandleCreateStore)
	mux.HandleFunc("GET /stores/{store_id}", handler.HandleGetStore)
	mux.HandleFunc("PUT /stores/{store_id}", handler.HandleUpdateStore)
}
