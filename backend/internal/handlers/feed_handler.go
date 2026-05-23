package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/teamart/commerce-api/internal/feed"
	"github.com/teamart/commerce-api/internal/infra/database"
	rec "github.com/teamart/commerce-api/internal/recommendation"
	"github.com/teamart/commerce-api/pkg/logger"
)

type FeedHandler struct {
	svc  *feed.Service
	repo rec.Repository
	db   *database.Pool
	log  *logger.Logger
}

func NewFeedHandler(svc *feed.Service, repo rec.Repository, db *database.Pool, log *logger.Logger) *FeedHandler {
	return &FeedHandler{svc: svc, repo: repo, db: db, log: log}
}

// GET /feed?limit=N
func (h *FeedHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 10
	if l := q.Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}

	items, err := h.svc.GetFeedForUser("anonymous", limit)
	if err != nil {
		h.log.Errorf("get feed: %v", err)
		http.Error(w, "failed to get feed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// POST /feed/candidates   (ingest candidate JSON)
func (h *FeedHandler) IngestCandidate(w http.ResponseWriter, r *http.Request) {
	var c rec.RecommendationCandidate
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	if err := h.repo.SaveCandidate(ctx, h.db, c); err != nil {
		h.log.Errorf("save candidate: %v", err)
		http.Error(w, "failed to save", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func RegisterFeedRoutes(mux Router, h *FeedHandler) {
	mux.HandleFunc("GET /feed", h.GetFeed)
	mux.HandleFunc("POST /feed/candidates", h.IngestCandidate)
}
