package compliance

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ConsentType represents types of consent
type ConsentType string

const (
	ConsentMarketingEmail ConsentType = "marketing_email"
	ConsentAnalytics      ConsentType = "analytics"
	ConsentThirdPartyShar ConsentType = "third_party_sharing"
)

// ConsentRecord represents a user's consent record
type ConsentRecord struct {
	ID             int64
	UserID         int64
	ConsentType    ConsentType
	ConsentVersion string
	Given          bool
	IPAddress      string
	UserAgent      string
	GivenAt        time.Time
	ExpiresAt      *time.Time
}

// DataExport represents a data export request
type DataExport struct {
	ID           int64
	UserID       int64
	Status       string // 'pending', 'processing', 'ready', 'delivered', 'expired'
	ExportFormat string // 'json', 'csv'
	FilePath     string
	FileSize     int64
	RequestedAt  time.Time
	CompletedAt  *time.Time
	ExpiresAt    *time.Time
	DeliveredAt  *time.Time
}

// DataDeletionRequest represents a user's data deletion request
type DataDeletionRequest struct {
	ID              int64
	UserID          int64
	Status          string // 'pending', 'reviewing', 'approved', 'in_progress', 'completed'
	Reason          string
	RequestedAt     time.Time
	ApprovedAt      *time.Time
	ApprovedBy      *int64
	CompletedAt     *time.Time
	GracePeriodEnds *time.Time
}

// ComplianceStorage defines storage interface for compliance
type ComplianceStorage interface {
	SaveConsentRecord(ctx context.Context, record *ConsentRecord) error
	GetConsentRecord(ctx context.Context, userID int64, consentType ConsentType) (*ConsentRecord, error)
	UpdateConsentRecord(ctx context.Context, record *ConsentRecord) error
	GetAllConsents(ctx context.Context, userID int64) ([]*ConsentRecord, error)

	SaveDataExport(ctx context.Context, export *DataExport) error
	GetDataExport(ctx context.Context, exportID int64) (*DataExport, error)
	UpdateDataExport(ctx context.Context, export *DataExport) error
	ListDataExports(ctx context.Context, userID int64) ([]*DataExport, error)

	SaveDeletionRequest(ctx context.Context, request *DataDeletionRequest) error
	GetDeletionRequest(ctx context.Context, requestID int64) (*DataDeletionRequest, error)
	GetPendingDeletionRequest(ctx context.Context, userID int64) (*DataDeletionRequest, error)
	UpdateDeletionRequest(ctx context.Context, request *DataDeletionRequest) error
	ListDeletionRequests(ctx context.Context, status string) ([]*DataDeletionRequest, error)
}

// ComplianceEngine manages GDPR and CCPA compliance
type ComplianceEngine struct {
	storage ComplianceStorage
	config  *ComplianceConfig
}

// ComplianceConfig holds compliance configuration
type ComplianceConfig struct {
	GDPREnabled            bool
	CCPAEnabled            bool
	DataExportExpiry       time.Duration
	DeletionGracePeriod    time.Duration
	DeletionProcessingTime time.Duration
}

// NewComplianceEngine creates a new compliance engine
func NewComplianceEngine(storage ComplianceStorage, config *ComplianceConfig) *ComplianceEngine {
	if config == nil {
		config = &ComplianceConfig{
			GDPREnabled:            true,
			CCPAEnabled:            true,
			DataExportExpiry:       7 * 24 * time.Hour,  // 7 days
			DeletionGracePeriod:    30 * 24 * time.Hour, // 30 days
			DeletionProcessingTime: 30 * 24 * time.Hour, // 30 days
		}
	}

	return &ComplianceEngine{
		storage: storage,
		config:  config,
	}
}

// RequestDataExport creates a data export request
func (c *ComplianceEngine) RequestDataExport(ctx context.Context, userID int64, format string) (*DataExport, error) {
	if userID == 0 {
		return nil, errors.New("user_id is required")
	}

	if format != "json" && format != "csv" {
		return nil, errors.New("format must be 'json' or 'csv'")
	}

	export := &DataExport{
		UserID:       userID,
		Status:       "pending",
		ExportFormat: format,
		RequestedAt:  time.Now(),
	}

	// Set expiry date
	expiresAt := time.Now().Add(c.config.DataExportExpiry)
	export.ExpiresAt = &expiresAt

	if err := c.storage.SaveDataExport(ctx, export); err != nil {
		return nil, fmt.Errorf("failed to save export request: %w", err)
	}

	return export, nil
}

// GetDataExports retrieves all export requests for a user
func (c *ComplianceEngine) GetDataExports(ctx context.Context, userID int64) ([]*DataExport, error) {
	return c.storage.ListDataExports(ctx, userID)
}

// RequestDataDeletion creates a data deletion request
func (c *ComplianceEngine) RequestDataDeletion(ctx context.Context, userID int64, reason string) (*DataDeletionRequest, error) {
	if userID == 0 {
		return nil, errors.New("user_id is required")
	}

	// Check if there's already a pending deletion request
	existing, err := c.storage.GetPendingDeletionRequest(ctx, userID)
	if err == nil && existing != nil {
		return nil, errors.New("deletion request already pending for this user")
	}

	request := &DataDeletionRequest{
		UserID:      userID,
		Status:      "pending",
		Reason:      reason,
		RequestedAt: time.Now(),
	}

	// Set grace period end date
	gracePeriodEnds := time.Now().Add(c.config.DeletionGracePeriod)
	request.GracePeriodEnds = &gracePeriodEnds

	if err := c.storage.SaveDeletionRequest(ctx, request); err != nil {
		return nil, fmt.Errorf("failed to save deletion request: %w", err)
	}

	return request, nil
}

// ApproveDeletion approves a data deletion request
func (c *ComplianceEngine) ApproveDeletion(ctx context.Context, requestID int64, approvedBy int64) error {
	request, err := c.storage.GetDeletionRequest(ctx, requestID)
	if err != nil {
		return fmt.Errorf("deletion request not found: %w", err)
	}

	if request == nil {
		return errors.New("deletion request not found")
	}

	request.Status = "approved"
	request.ApprovedAt = timePtr(time.Now())
	request.ApprovedBy = &approvedBy

	// Set completion date
	completesAt := time.Now().Add(c.config.DeletionProcessingTime)
	request.CompletedAt = &completesAt

	return c.storage.UpdateDeletionRequest(ctx, request)
}

// CancelDeletion cancels a pending deletion request
func (c *ComplianceEngine) CancelDeletion(ctx context.Context, requestID int64) error {
	request, err := c.storage.GetDeletionRequest(ctx, requestID)
	if err != nil {
		return fmt.Errorf("deletion request not found: %w", err)
	}

	if request == nil {
		return errors.New("deletion request not found")
	}

	// Only allow cancellation during grace period
	if request.GracePeriodEnds != nil && time.Now().After(*request.GracePeriodEnds) {
		return errors.New("grace period has expired - deletion cannot be cancelled")
	}

	request.Status = "cancelled"
	return c.storage.UpdateDeletionRequest(ctx, request)
}

// GetDeletionRequest retrieves a deletion request
func (c *ComplianceEngine) GetDeletionRequest(ctx context.Context, requestID int64) (*DataDeletionRequest, error) {
	return c.storage.GetDeletionRequest(ctx, requestID)
}

// SaveConsent saves user consent
func (c *ComplianceEngine) SaveConsent(ctx context.Context, userID int64, consentType ConsentType, given bool, ipAddress, userAgent string) error {
	if userID == 0 {
		return errors.New("user_id is required")
	}

	record := &ConsentRecord{
		UserID:      userID,
		ConsentType: consentType,
		Given:       given,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		GivenAt:     time.Now(),
	}

	return c.storage.SaveConsentRecord(ctx, record)
}

// GetConsent gets user's consent status
func (c *ComplianceEngine) GetConsent(ctx context.Context, userID int64, consentType ConsentType) (bool, error) {
	record, err := c.storage.GetConsentRecord(ctx, userID, consentType)
	if err != nil {
		return false, err
	}

	if record == nil {
		return false, nil
	}

	return record.Given, nil
}

// GetAllConsents gets all user consents
func (c *ComplianceEngine) GetAllConsents(ctx context.Context, userID int64) ([]*ConsentRecord, error) {
	return c.storage.GetAllConsents(ctx, userID)
}

// ExportUserData exports all user data (for GDPR right to data portability)
func (c *ComplianceEngine) ExportUserData(ctx context.Context, userID int64) (map[string]interface{}, error) {
	if userID == 0 {
		return nil, errors.New("user_id is required")
	}

	// This would be implemented to gather all user data from all services
	userData := make(map[string]interface{})
	userData["user_id"] = userID
	userData["export_timestamp"] = time.Now()

	// Include various data categories
	userData["profile"] = nil     // To be populated
	userData["sessions"] = nil    // To be populated
	userData["auth_events"] = nil // To be populated
	userData["consents"] = nil    // To be populated
	userData["audit_log"] = nil   // To be populated

	return userData, nil
}

// DeleteUserData permanently deletes user data (for GDPR right to be forgotten)
func (c *ComplianceEngine) DeleteUserData(ctx context.Context, userID int64) error {
	if userID == 0 {
		return errors.New("user_id is required")
	}

	// This would be implemented to delete user data from all services
	// This is a cascading operation that requires careful coordination

	return nil
}

// Helper function to create a pointer to time.Time
func timePtr(t time.Time) *time.Time {
	return &t
}
