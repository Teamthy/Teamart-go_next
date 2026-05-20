package audit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// AuditEvent represents an audit log entry
type AuditEvent struct {
	ID           int64
	EventID      string // UUID for idempotency
	EventType    string // 'user.created', 'auth.success', 'user.data_export'
	ActorID      *int64
	ActorType    string // 'user', 'system', 'admin'
	ResourceType string
	ResourceID   string
	Action       string // 'create', 'update', 'delete', 'login'
	Status       string // 'success', 'failure'
	IPAddress    string
	UserAgent    string
	Details      map[string]interface{}
	Hash         string // SHA256 of previous entry for chain
	CreatedAt    time.Time
}

// AuditStorage defines storage interface for audit logs
type AuditStorage interface {
	SaveAuditEvent(ctx context.Context, event *AuditEvent) error
	GetAuditEvent(ctx context.Context, id int64) (*AuditEvent, error)
	SearchAuditEvents(ctx context.Context, filter *AuditFilter) ([]*AuditEvent, error)
	GetAuditEventsByUser(ctx context.Context, userID int64, limit int) ([]*AuditEvent, error)
	GetAuditEventsByResource(ctx context.Context, resourceType, resourceID string, limit int) ([]*AuditEvent, error)
	GetLastAuditEvent(ctx context.Context) (*AuditEvent, error)
	VerifyAuditChain(ctx context.Context, startID, endID int64) (bool, error)
}

// AuditFilter defines search filters for audit events
type AuditFilter struct {
	EventType    string
	ActorID      *int64
	ActorType    string
	ResourceType string
	ResourceID   string
	Action       string
	Status       string
	StartTime    time.Time
	EndTime      time.Time
	Limit        int
	Offset       int
}

// AuditLogger logs immutable audit events
type AuditLogger struct {
	storage AuditStorage
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(storage AuditStorage) *AuditLogger {
	return &AuditLogger{
		storage: storage,
	}
}

// LogEvent logs an audit event
func (l *AuditLogger) LogEvent(ctx context.Context, event *AuditEvent) error {
	if event == nil {
		return errors.New("event is required")
	}

	if event.EventID == "" {
		return errors.New("event_id is required")
	}

	if event.EventType == "" {
		return errors.New("event_type is required")
	}

	// Set timestamp if not provided
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	// Calculate hash for audit chain integrity
	hash, err := l.calculateEventHash(event)
	if err != nil {
		return fmt.Errorf("failed to calculate hash: %w", err)
	}
	event.Hash = hash

	// Save event
	if err := l.storage.SaveAuditEvent(ctx, event); err != nil {
		return fmt.Errorf("failed to save audit event: %w", err)
	}

	return nil
}

// SearchEvents searches audit events
func (l *AuditLogger) SearchEvents(ctx context.Context, filter *AuditFilter) ([]*AuditEvent, error) {
	if filter == nil {
		return nil, errors.New("filter is required")
	}

	return l.storage.SearchAuditEvents(ctx, filter)
}

// GetUserAuditTrail retrieves audit trail for a specific user
func (l *AuditLogger) GetUserAuditTrail(ctx context.Context, userID int64, limit int) ([]*AuditEvent, error) {
	if limit == 0 {
		limit = 100
	}

	return l.storage.GetAuditEventsByUser(ctx, userID, limit)
}

// GetResourceAuditTrail retrieves audit trail for a specific resource
func (l *AuditLogger) GetResourceAuditTrail(ctx context.Context, resourceType, resourceID string, limit int) ([]*AuditEvent, error) {
	if resourceType == "" || resourceID == "" {
		return nil, errors.New("resource_type and resource_id are required")
	}

	if limit == 0 {
		limit = 100
	}

	return l.storage.GetAuditEventsByResource(ctx, resourceType, resourceID, limit)
}

// VerifyIntegrity verifies the integrity of the audit chain
func (l *AuditLogger) VerifyIntegrity(ctx context.Context, startID, endID int64) (bool, error) {
	return l.storage.VerifyAuditChain(ctx, startID, endID)
}

// ExportEvents exports audit events
func (l *AuditLogger) ExportEvents(ctx context.Context, filter *AuditFilter) ([]*AuditEvent, error) {
	if filter == nil {
		return nil, errors.New("filter is required")
	}

	// Set a reasonable limit for exports
	if filter.Limit == 0 || filter.Limit > 100000 {
		filter.Limit = 100000
	}

	return l.storage.SearchAuditEvents(ctx, filter)
}

// calculateEventHash calculates SHA256 hash for audit event
func (l *AuditLogger) calculateEventHash(event *AuditEvent) (string, error) {
	// Get previous event's hash for chaining
	lastEvent, err := l.storage.GetLastAuditEvent(context.Background())
	var prevHash string
	if err == nil && lastEvent != nil {
		prevHash = lastEvent.Hash
	}

	// Create hash input from event data
	hashInput := fmt.Sprintf("%s|%s|%s|%s|%v|%s",
		prevHash,
		event.EventID,
		event.EventType,
		event.CreatedAt.String(),
		event.Details,
		event.Action,
	)

	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:]), nil
}

// LogUserCreated logs user creation event
func (l *AuditLogger) LogUserCreated(ctx context.Context, eventID string, userID int64, email string, ipAddress string) error {
	event := &AuditEvent{
		EventID:      eventID,
		EventType:    "user.created",
		ActorType:    "user",
		ResourceType: "user",
		ResourceID:   fmt.Sprintf("%d", userID),
		Action:       "create",
		Status:       "success",
		IPAddress:    ipAddress,
		Details: map[string]interface{}{
			"email": email,
		},
	}

	return l.LogEvent(ctx, event)
}

// LogAuthSuccess logs successful authentication
func (l *AuditLogger) LogAuthSuccess(ctx context.Context, eventID string, userID int64, ipAddress, userAgent string) error {
	event := &AuditEvent{
		EventID:      eventID,
		EventType:    "auth.login_success",
		ActorID:      &userID,
		ActorType:    "user",
		ResourceType: "session",
		Action:       "login",
		Status:       "success",
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}

	return l.LogEvent(ctx, event)
}

// LogAuthFailure logs failed authentication
func (l *AuditLogger) LogAuthFailure(ctx context.Context, eventID, reason, ipAddress, userAgent string) error {
	event := &AuditEvent{
		EventID:      eventID,
		EventType:    "auth.login_failure",
		ActorType:    "user",
		ResourceType: "session",
		Action:       "login",
		Status:       "failure",
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Details: map[string]interface{}{
			"reason": reason,
		},
	}

	return l.LogEvent(ctx, event)
}

// LogDataExport logs data export requests
func (l *AuditLogger) LogDataExport(ctx context.Context, eventID string, userID int64, ipAddress string) error {
	event := &AuditEvent{
		EventID:      eventID,
		EventType:    "user.data_export",
		ActorID:      &userID,
		ActorType:    "user",
		ResourceType: "user",
		ResourceID:   fmt.Sprintf("%d", userID),
		Action:       "export",
		Status:       "success",
		IPAddress:    ipAddress,
	}

	return l.LogEvent(ctx, event)
}

// LogDeletion logs data deletion events
func (l *AuditLogger) LogDeletion(ctx context.Context, eventID string, userID int64, reason string, ipAddress string) error {
	event := &AuditEvent{
		EventID:      eventID,
		EventType:    "user.deleted",
		ActorID:      &userID,
		ActorType:    "user",
		ResourceType: "user",
		ResourceID:   fmt.Sprintf("%d", userID),
		Action:       "delete",
		Status:       "success",
		IPAddress:    ipAddress,
		Details: map[string]interface{}{
			"reason": reason,
		},
	}

	return l.LogEvent(ctx, event)
}
