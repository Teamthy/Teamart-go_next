package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/teamart/commerce-api/internal/admin"
	"github.com/teamart/commerce-api/pkg/logger"
)

type AdminHandler struct {
	svc admin.Service
	log *logger.Logger
}

func NewAdminHandler(s admin.Service, log *logger.Logger) *AdminHandler {
	return &AdminHandler{svc: s, log: log}
}

func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s, err := h.svc.GetDashboard(ctx)
	if err != nil {
		h.log.Errorf("dashboard: %v", err)
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(s)
}

func (h *AdminHandler) ListDisputes(w http.ResponseWriter, r *http.Request) {
	ds, err := h.svc.ListDisputes(context.Background())
	if err != nil {
		h.log.Errorf("list disputes: %v", err)
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(ds)
}

func (h *AdminHandler) CreateDispute(w http.ResponseWriter, r *http.Request) {
	var d admin.Dispute
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	if err := h.svc.CreateDispute(context.Background(), d); err != nil {
		h.log.Errorf("create dispute: %v", err)
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *AdminHandler) ApprovePayout(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := q.Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	if err := h.svc.ApprovePayout(context.Background(), id); err != nil {
		h.log.Errorf("approve payout: %v", err)
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) ListFraudAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.svc.ListFraudAlerts(context.Background())
	if err != nil {
		h.log.Errorf("list alerts: %v", err)
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(alerts)
}

func (h *AdminHandler) VerifyCreator(w http.ResponseWriter, r *http.Request) {
	var v admin.CreatorVerification
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	if err := h.svc.VerifyCreator(context.Background(), v); err != nil {
		h.log.Errorf("verify creator: %v", err)
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) Refund(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := q.Get("dispute_id")
	if id == "" {
		http.Error(w, "missing dispute_id", http.StatusBadRequest)
		return
	}
	if err := h.svc.Refund(context.Background(), id); err != nil {
		h.log.Errorf("refund: %v", err)
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) SuspendAccount(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := q.Get("user_id")
	if id == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}
	if err := h.svc.SuspendAccount(context.Background(), id); err != nil {
		h.log.Errorf("suspend: %v", err)
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// POST /admin/payouts/request?id=<payout_id>&requested_by=<user_id>
func (h *AdminHandler) RequestPayoutApproval(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pid := q.Get("id")
	rb := q.Get("requested_by")
	if pid == "" || rb == "" {
		http.Error(w, "missing params", http.StatusBadRequest)
		return
	}
	// parse requested_by
	var reqBy int64
	if v, err := strconv.ParseInt(rb, 10, 64); err == nil {
		reqBy = v
	}
	notes := r.URL.Query().Get("notes")
	if err := h.svc.RequestPayoutApproval(context.Background(), pid, reqBy, notes); err != nil {
		h.log.Errorf("request payout approval: %v", err)
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *AdminHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.svc.ListAuditLogs(context.Background())
	if err != nil {
		h.log.Errorf("list audit: %v", err)
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(logs)
}

func (h *AdminHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	notes, err := h.svc.ListNotifications(context.Background())
	if err != nil {
		h.log.Errorf("list notifications: %v", err)
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(notes)
}

func RegisterAdminRoutes(mux Router, h *AdminHandler) {
	mux.HandleFunc("GET /admin/dashboard", h.Dashboard)
	mux.HandleFunc("GET /admin/disputes", h.ListDisputes)
	mux.HandleFunc("POST /admin/disputes", h.CreateDispute)
	mux.HandleFunc("POST /admin/payouts/approve", h.ApprovePayout)
	mux.HandleFunc("GET /admin/fraud/alerts", h.ListFraudAlerts)
	mux.HandleFunc("POST /admin/creators/verify", h.VerifyCreator)
	mux.HandleFunc("POST /admin/support/refund", h.Refund)
	mux.HandleFunc("POST /admin/support/suspend", h.SuspendAccount)
}
