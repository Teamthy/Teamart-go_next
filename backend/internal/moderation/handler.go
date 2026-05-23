package moderation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

type Handler struct {
	service *ModerationService
	logger  *logger.Logger
}

func NewHandler(service *ModerationService, logger *logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func RegisterModerationRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("/api/v1/moderation/users", handler.handleUsers)
	mux.HandleFunc("/api/v1/moderation/users/", handler.handleUser)
}

func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleListStatus(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleUser(w http.ResponseWriter, r *http.Request) {
	userID, action := parseUserPath(r.URL.Path)
	if userID == 0 {
		writeJSONError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	switch strings.ToLower(action) {
	case "status":
		if r.Method != http.MethodGet {
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.handleStatus(w, r, userID)
	case "moderator":
		h.handleModerator(w, r, userID)
	case "block":
		h.handleBlock(w, r, userID)
	case "shadowban":
		h.handleShadowBan(w, r, userID)
	case "mute":
		h.handleMute(w, r, userID)
	case "permissions":
		h.handlePermissions(w, r, userID)
	default:
		writeJSONError(w, http.StatusNotFound, "not found")
	}
}

func (h *Handler) handleListStatus(w http.ResponseWriter, r *http.Request) {
	writeJSONError(w, http.StatusNotImplemented, "bulk status listing is not supported")
}

func (h *Handler) handleStatus(w http.ResponseWriter, r *http.Request, userID int64) {
	status := h.service.GetUserStatus(userID)
	writeJSON(w, http.StatusOK, status)
}

func (h *Handler) handleModerator(w http.ResponseWriter, r *http.Request, userID int64) {
	switch r.Method {
	case http.MethodPost:
		h.service.AddModerator(userID)
		writeJSON(w, http.StatusOK, map[string]string{"status": "moderator added"})
	case http.MethodDelete:
		h.service.RemoveModerator(userID)
		writeJSON(w, http.StatusOK, map[string]string{"status": "moderator removed"})
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleBlock(w http.ResponseWriter, r *http.Request, userID int64) {
	switch r.Method {
	case http.MethodPost:
		h.service.BlockUser(userID)
		writeJSON(w, http.StatusOK, map[string]string{"status": "user blocked"})
	case http.MethodDelete:
		h.service.UnblockUser(userID)
		writeJSON(w, http.StatusOK, map[string]string{"status": "user unblocked"})
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleShadowBan(w http.ResponseWriter, r *http.Request, userID int64) {
	switch r.Method {
	case http.MethodPost:
		h.service.ShadowBanUser(userID)
		writeJSON(w, http.StatusOK, map[string]string{"status": "user shadow-banned"})
	case http.MethodDelete:
		h.service.UnshadowBanUser(userID)
		writeJSON(w, http.StatusOK, map[string]string{"status": "user removed from shadow ban"})
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleMute(w http.ResponseWriter, r *http.Request, userID int64) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		DurationSeconds int `json:"duration_seconds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.DurationSeconds <= 0 {
		writeJSONError(w, http.StatusBadRequest, "duration_seconds must be greater than zero")
		return
	}
	h.service.MuteUser(userID, time.Duration(req.DurationSeconds)*time.Second)
	writeJSON(w, http.StatusOK, map[string]string{"status": "user muted"})
}

func (h *Handler) handlePermissions(w http.ResponseWriter, r *http.Request, userID int64) {
	switch r.Method {
	case http.MethodPost:
		var req struct {
			Permissions []string `json:"permissions"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		h.service.ConfigureUserPermissions(userID, req.Permissions)
		writeJSON(w, http.StatusOK, map[string]string{"status": "permissions configured"})
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func parseUserPath(path string) (int64, string) {
	trimmed := strings.TrimPrefix(path, "/api/v1/moderation/users/")
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) == 0 || parts[0] == "" {
		return 0, ""
	}
	userID := int64(0)
	fmt.Sscan(parts[0], &userID)
	if len(parts) == 1 {
		return userID, ""
	}
	return userID, parts[1]
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
