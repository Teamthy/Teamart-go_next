package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/teamart/commerce-api/internal/users"
	"github.com/teamart/commerce-api/pkg/logger"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	service *users.Service
	logger  *logger.Logger
}

// NewUserHandler creates a new user HTTP handler
func NewUserHandler(service *users.Service, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// CreateUserRequest represents the HTTP request body for creating a user
type CreateUserRequest struct {
	Email        string `json:"email" binding:"required"`
	Name         string `json:"name" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
}

// UserResponse represents the HTTP response body for a user
type UserResponse struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// HandleCreateUser handles POST /users requests
//
//	Example: curl -X POST http://localhost:8080/users \
//	  -H "Content-Type: application/json" \
//	  -d '{"email":"user@example.com","name":"John Doe","password_hash":"hashed_password"}'
func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling CreateUser request")

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := &users.CreateUserInput{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: req.PasswordHash,
	}

	output, err := h.service.CreateUser(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UserResponse{
		ID:           output.ID,
		Email:        output.Email,
		Name:         output.Name,
		PasswordHash: output.PasswordHash,
		CreatedAt:    output.CreatedAt,
		UpdatedAt:    output.UpdatedAt,
	})
}

// HandleGetUser handles GET /users/:id requests
// Example: curl http://localhost:8080/users/1
func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	h.logger.Debugf("handling GetUser request for user: %d", userID)

	input := &users.GetUserByIDInput{UserID: userID}
	output, err := h.service.GetUserByID(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserResponse{
		ID:           output.ID,
		Email:        output.Email,
		Name:         output.Name,
		PasswordHash: output.PasswordHash,
		CreatedAt:    output.CreatedAt,
		UpdatedAt:    output.UpdatedAt,
	})
}

// HandleGetUserByEmail handles GET /users/email/:email requests
// Example: curl http://localhost:8080/users/email/user@example.com
func (h *UserHandler) HandleGetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	h.logger.Debugf("handling GetUserByEmail request for email: %s", email)

	input := &users.GetUserByEmailInput{Email: email}
	output, err := h.service.GetUserByEmail(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserResponse{
		ID:           output.ID,
		Email:        output.Email,
		Name:         output.Name,
		PasswordHash: output.PasswordHash,
		CreatedAt:    output.CreatedAt,
		UpdatedAt:    output.UpdatedAt,
	})
}

// ListUsersRequest represents the query parameters for listing users
type ListUsersRequest struct {
	Limit  int32
	Offset int32
}

// ListUsersResponse represents the HTTP response body
type ListUsersResponse struct {
	Users  []UserResponse `json:"users"`
	Total  int32          `json:"total"`
	Limit  int32          `json:"limit"`
	Offset int32          `json:"offset"`
}

// HandleListUsers handles GET /users requests with pagination
// Example: curl "http://localhost:8080/users?limit=10&offset=0"
func (h *UserHandler) HandleListUsers(w http.ResponseWriter, r *http.Request) {
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

	h.logger.Debugf("handling ListUsers request with limit: %d, offset: %d", limit, offset)

	input := &users.ListUsersInput{Limit: limit, Offset: offset}
	output, err := h.service.ListUsers(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userResponses := make([]UserResponse, len(output.Users))
	for i, user := range output.Users {
		userResponses[i] = UserResponse{
			ID:           user.ID,
			Email:        user.Email,
			Name:         user.Name,
			PasswordHash: user.PasswordHash,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ListUsersResponse{
		Users:  userResponses,
		Limit:  output.Limit,
		Offset: output.Offset,
	})
}

// UpdateUserRequest represents the HTTP request body for updating a user
type UpdateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

// HandleUpdateUser handles PUT /users/:id requests
//
//	Example: curl -X PUT http://localhost:8080/users/1 \
//	  -H "Content-Type: application/json" \
//	  -d '{"name":"Updated Name"}'
func (h *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Debugf("handling UpdateUser request for user: %d", userID)

	input := &users.UpdateUserInput{UserID: userID, Name: req.Name}
	output, err := h.service.UpdateUser(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserResponse{
		ID:           output.ID,
		Email:        output.Email,
		Name:         output.Name,
		PasswordHash: output.PasswordHash,
		CreatedAt:    output.CreatedAt,
		UpdatedAt:    output.UpdatedAt,
	})
}

// HandleDeleteUser handles DELETE /users/:id requests
// Example: curl -X DELETE http://localhost:8080/users/1
func (h *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	h.logger.Debugf("handling DeleteUser request for user: %d", userID)

	input := &users.DeleteUserInput{UserID: userID}
	err = h.service.DeleteUser(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(mux *http.ServeMux, handler *UserHandler) {
	// User endpoints
	mux.HandleFunc("POST /users", handler.HandleCreateUser)
	mux.HandleFunc("GET /users", handler.HandleListUsers)
	mux.HandleFunc("GET /users/{id}", handler.HandleGetUser)
	mux.HandleFunc("GET /users/email/{email}", handler.HandleGetUserByEmail)
	mux.HandleFunc("PUT /users/{id}", handler.HandleUpdateUser)
	mux.HandleFunc("DELETE /users/{id}", handler.HandleDeleteUser)
}
