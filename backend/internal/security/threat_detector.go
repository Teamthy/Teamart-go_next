package security

import (
	"context"
	"fmt"
	"time"
)

// SecurityEvent represents a security event
type SecurityEvent struct {
	ID              int64
	EventType       string // 'brute_force', 'credential_stuffing', 'token_abuse', etc.
	Severity        string // 'info', 'warning', 'critical'
	UserID          *int64
	IPAddress       string
	Details         map[string]interface{}
	TriggeredAction string // 'block_ip', 'revoke_sessions', 'mfa_required'
	Resolved        bool
	CreatedAt       time.Time
	ResolvedAt      *time.Time
}

// SecurityAlert represents a security alert
type SecurityAlert struct {
	ID             int64
	AlertType      string
	Severity       string
	Title          string
	Description    string
	AffectedUsers  int
	ActionRequired bool
	Resolved       bool
	AssignedTo     *int64
	CreatedAt      time.Time
	ResolvedAt     *time.Time
}

// ThreatDetector detects security threats in real-time
type ThreatDetector struct {
	storage ThreatStorage
	config  *ThreatDetectorConfig
}

// ThreatDetectorConfig holds threat detector configuration
type ThreatDetectorConfig struct {
	BruteForceThreshold      int
	BruteForceWindow         time.Duration
	CredentialStuffingWindow time.Duration
	CredentialStuffingLimit  int
	ImpossibleTravelMinSpeed int // km/h
	RiskScoreThreshold       int
}

// ThreatStorage defines storage interface for threat detection
type ThreatStorage interface {
	SaveSecurityEvent(ctx context.Context, event *SecurityEvent) error
	GetSecurityEvents(ctx context.Context, userID *int64, limit int) ([]*SecurityEvent, error)
	SaveSecurityAlert(ctx context.Context, alert *SecurityAlert) error
	GetActiveAlerts(ctx context.Context) ([]*SecurityAlert, error)
	ResolveAlert(ctx context.Context, alertID int64) error
	GetLoginAttempts(ctx context.Context, ipOrEmail string, duration time.Duration) (int, error)
	SaveLoginAttempt(ctx context.Context, email, ipAddress string, success bool) error
}

// NewThreatDetector creates a new threat detector
func NewThreatDetector(storage ThreatStorage, config *ThreatDetectorConfig) *ThreatDetector {
	if config == nil {
		config = &ThreatDetectorConfig{
			BruteForceThreshold:      10,
			BruteForceWindow:         5 * time.Minute,
			CredentialStuffingWindow: 10 * time.Minute,
			CredentialStuffingLimit:  100,
			ImpossibleTravelMinSpeed: 900, // km/h
			RiskScoreThreshold:       70,
		}
	}

	return &ThreatDetector{
		storage: storage,
		config:  config,
	}
}

// DetectBruteForce detects brute force attacks
func (d *ThreatDetector) DetectBruteForce(ctx context.Context, email, ipAddress string) (bool, error) {
	// Count failed login attempts in the last 5 minutes
	attempts, err := d.storage.GetLoginAttempts(ctx, email, d.config.BruteForceWindow)
	if err != nil {
		return false, fmt.Errorf("failed to get login attempts: %w", err)
	}

	if attempts >= d.config.BruteForceThreshold {
		// Record security event
		event := &SecurityEvent{
			EventType:       "brute_force",
			Severity:        "critical",
			IPAddress:       ipAddress,
			Details:         map[string]interface{}{"attempts": attempts, "email": email},
			TriggeredAction: "block_ip",
			CreatedAt:       time.Now(),
		}

		if err := d.storage.SaveSecurityEvent(ctx, event); err != nil {
			return false, fmt.Errorf("failed to save security event: %w", err)
		}

		// Create alert
		alert := &SecurityAlert{
			AlertType:      "brute_force",
			Severity:       "critical",
			Title:          "Brute Force Attack Detected",
			Description:    fmt.Sprintf("%d failed login attempts from IP %s", attempts, ipAddress),
			ActionRequired: true,
			CreatedAt:      time.Now(),
		}

		if err := d.storage.SaveSecurityAlert(ctx, alert); err != nil {
			return false, fmt.Errorf("failed to save alert: %w", err)
		}

		return true, nil
	}

	return false, nil
}

// DetectCredentialStuffing detects credential stuffing attacks
func (d *ThreatDetector) DetectCredentialStuffing(ctx context.Context, ipAddress string) (bool, error) {
	// Count login attempts from single IP in short window
	attempts, err := d.storage.GetLoginAttempts(ctx, ipAddress, d.config.CredentialStuffingWindow)
	if err != nil {
		return false, fmt.Errorf("failed to get login attempts: %w", err)
	}

	if attempts >= d.config.CredentialStuffingLimit {
		// Record security event
		event := &SecurityEvent{
			EventType:       "credential_stuffing",
			Severity:        "critical",
			IPAddress:       ipAddress,
			Details:         map[string]interface{}{"attempts": attempts},
			TriggeredAction: "block_ip",
			CreatedAt:       time.Now(),
		}

		if err := d.storage.SaveSecurityEvent(ctx, event); err != nil {
			return false, fmt.Errorf("failed to save security event: %w", err)
		}

		// Create alert
		alert := &SecurityAlert{
			AlertType:      "credential_stuffing",
			Severity:       "critical",
			Title:          "Credential Stuffing Attack Detected",
			Description:    fmt.Sprintf("%d login attempts from IP %s", attempts, ipAddress),
			ActionRequired: true,
			CreatedAt:      time.Now(),
		}

		if err := d.storage.SaveSecurityAlert(ctx, alert); err != nil {
			return false, fmt.Errorf("failed to save alert: %w", err)
		}

		return true, nil
	}

	return false, nil
}

// RecordLoginAttempt records a login attempt for threat analysis
func (d *ThreatDetector) RecordLoginAttempt(ctx context.Context, email, ipAddress string, success bool) error {
	return d.storage.SaveLoginAttempt(ctx, email, ipAddress, success)
}

// GetSecurityEvents retrieves security events for a user
func (d *ThreatDetector) GetSecurityEvents(ctx context.Context, userID *int64, limit int) ([]*SecurityEvent, error) {
	return d.storage.GetSecurityEvents(ctx, userID, limit)
}

// GetActiveAlerts retrieves all active security alerts
func (d *ThreatDetector) GetActiveAlerts(ctx context.Context) ([]*SecurityAlert, error) {
	return d.storage.GetActiveAlerts(ctx)
}

// ResolveAlert marks an alert as resolved
func (d *ThreatDetector) ResolveAlert(ctx context.Context, alertID int64) error {
	return d.storage.ResolveAlert(ctx, alertID)
}
