package device

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Device represents a user's device
type Device struct {
	ID             int64
	UserID         int64
	DeviceID       string // Unique device identifier
	DeviceType     string // 'mobile', 'desktop', 'tablet'
	OSName         string
	OSVersion      string
	BrowserName    string
	BrowserVersion string
	DeviceModel    string
	DeviceName     string // User-friendly name
	Trusted        bool
	TrustExpiresAt *time.Time
	FirstSeen      time.Time
	LastUsed       time.Time
	IsActive       bool
}

// DeviceRiskScore represents the risk level of a device
type DeviceRiskScore struct {
	DeviceID     int64
	RiskLevel    string // 'trusted', 'low_risk', 'medium_risk', 'high_risk', 'blocked'
	RiskScore    int    // 0-100
	Factors      map[string]interface{}
	CalculatedAt time.Time
}

// DeviceFingerprint represents a device fingerprint
type DeviceFingerprint struct {
	ID                 int64
	DeviceID           int64
	BrowserFingerprint string
	CanvasFingerprint  string
	WebGLFingerprint   string
	ScreenDimensions   string
	ScreenResolution   string
	Timezone           string
	Language           string
	UserAgent          string
	AcceptLanguage     string
	Plugins            []string
	Fonts              []string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// DeviceService manages device tracking and trust
type DeviceService struct {
	storage DeviceStorage
	config  *DeviceConfig
}

// DeviceConfig holds device configuration
type DeviceConfig struct {
	RequireDeviceVerification  bool
	DeviceTrustExpiry          time.Duration
	EnableFingerprinting       bool
	EnableGeoLocation          bool
	FingerprintRefreshInterval time.Duration
}

// DeviceStorage defines storage interface for devices
type DeviceStorage interface {
	SaveDevice(ctx context.Context, device *Device) error
	GetDevice(ctx context.Context, deviceID string) (*Device, error)
	GetDeviceByID(ctx context.Context, id int64) (*Device, error)
	UpdateDevice(ctx context.Context, device *Device) error
	DeleteDevice(ctx context.Context, id int64) error
	ListDevices(ctx context.Context, userID int64) ([]*Device, error)
	ListTrustedDevices(ctx context.Context, userID int64) ([]*Device, error)

	SaveDeviceFingerprint(ctx context.Context, fp *DeviceFingerprint) error
	GetDeviceFingerprint(ctx context.Context, deviceID int64) (*DeviceFingerprint, error)
	UpdateDeviceFingerprint(ctx context.Context, fp *DeviceFingerprint) error

	SaveRiskScore(ctx context.Context, score *DeviceRiskScore) error
	GetRiskScore(ctx context.Context, deviceID int64) (*DeviceRiskScore, error)
}

// NewDeviceService creates a new device service
func NewDeviceService(storage DeviceStorage, config *DeviceConfig) *DeviceService {
	if config == nil {
		config = &DeviceConfig{
			RequireDeviceVerification:  true,
			DeviceTrustExpiry:          90 * 24 * time.Hour, // 90 days
			EnableFingerprinting:       true,
			EnableGeoLocation:          true,
			FingerprintRefreshInterval: 30 * 24 * time.Hour, // 30 days
		}
	}

	return &DeviceService{
		storage: storage,
		config:  config,
	}
}

// RegisterDevice registers a new device for a user
func (s *DeviceService) RegisterDevice(ctx context.Context, userID int64, fingerprint *DeviceFingerprint, name string) (*Device, error) {
	if userID == 0 {
		return nil, errors.New("user_id is required")
	}

	if fingerprint == nil {
		return nil, errors.New("fingerprint is required")
	}

	// Generate device ID from fingerprint
	deviceID := s.generateDeviceID(fingerprint)

	// Check if device already exists
	existing, err := s.storage.GetDevice(ctx, deviceID)
	if err == nil && existing != nil {
		// Device already registered, update last used
		existing.LastUsed = time.Now()
		if err := s.storage.UpdateDevice(ctx, existing); err != nil {
			return nil, fmt.Errorf("failed to update device: %w", err)
		}
		return existing, nil
	}

	// Create new device
	now := time.Now()
	device := &Device{
		UserID:     userID,
		DeviceID:   deviceID,
		DeviceType: fingerprint.ScreenResolution, // Could be improved
		OSName:     extractOSName(fingerprint.UserAgent),
		Trusted:    false,
		FirstSeen:  now,
		LastUsed:   now,
		IsActive:   true,
		DeviceName: name,
	}

	if err := s.storage.SaveDevice(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to save device: %w", err)
	}

	// Save fingerprint
	fingerprint.DeviceID = device.ID
	fingerprint.CreatedAt = now
	fingerprint.UpdatedAt = now
	if err := s.storage.SaveDeviceFingerprint(ctx, fingerprint); err != nil {
		return nil, fmt.Errorf("failed to save fingerprint: %w", err)
	}

	return device, nil
}

// TrustDevice marks a device as trusted
func (s *DeviceService) TrustDevice(ctx context.Context, deviceID int64, expiryDuration *time.Duration) error {
	device, err := s.storage.GetDeviceByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("device not found: %w", err)
	}

	device.Trusted = true
	if expiryDuration == nil {
		expiry := time.Now().Add(s.config.DeviceTrustExpiry)
		device.TrustExpiresAt = &expiry
	} else {
		expiry := time.Now().Add(*expiryDuration)
		device.TrustExpiresAt = &expiry
	}

	return s.storage.UpdateDevice(ctx, device)
}

// RevokeTrust removes trust from a device
func (s *DeviceService) RevokeTrust(ctx context.Context, deviceID int64) error {
	device, err := s.storage.GetDeviceByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("device not found: %w", err)
	}

	device.Trusted = false
	device.TrustExpiresAt = nil

	return s.storage.UpdateDevice(ctx, device)
}

// IsTrusted checks if a device is currently trusted
func (s *DeviceService) IsTrusted(ctx context.Context, deviceID int64) (bool, error) {
	device, err := s.storage.GetDeviceByID(ctx, deviceID)
	if err != nil {
		return false, fmt.Errorf("device not found: %w", err)
	}

	if !device.Trusted {
		return false, nil
	}

	// Check if trust has expired
	if device.TrustExpiresAt != nil && time.Now().After(*device.TrustExpiresAt) {
		// Trust has expired, revoke it
		_ = s.RevokeTrust(ctx, deviceID)
		return false, nil
	}

	return true, nil
}

// ListTrustedDevices lists all trusted devices for a user
func (s *DeviceService) ListTrustedDevices(ctx context.Context, userID int64) ([]*Device, error) {
	return s.storage.ListTrustedDevices(ctx, userID)
}

// ListAllDevices lists all devices for a user
func (s *DeviceService) ListAllDevices(ctx context.Context, userID int64) ([]*Device, error) {
	return s.storage.ListDevices(ctx, userID)
}

// RemoveDevice removes a device
func (s *DeviceService) RemoveDevice(ctx context.Context, deviceID int64) error {
	return s.storage.DeleteDevice(ctx, deviceID)
}

// UpdateDeviceName updates the display name of a device
func (s *DeviceService) UpdateDeviceName(ctx context.Context, deviceID int64, name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	device, err := s.storage.GetDeviceByID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("device not found: %w", err)
	}

	device.DeviceName = name
	device.LastUsed = time.Now()

	return s.storage.UpdateDevice(ctx, device)
}

// generateDeviceID generates a unique device ID from fingerprint
func (s *DeviceService) generateDeviceID(fp *DeviceFingerprint) string {
	data := fmt.Sprintf("%s|%s|%s|%s", fp.UserAgent, fp.ScreenDimensions, fp.Timezone, fp.Language)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// extractOSName extracts OS name from user agent
func extractOSName(userAgent string) string {
	// Simple heuristic - in production, use a proper user agent parser
	if contains(userAgent, "Windows") {
		return "Windows"
	}
	if contains(userAgent, "Mac") {
		return "macOS"
	}
	if contains(userAgent, "Linux") {
		return "Linux"
	}
	if contains(userAgent, "iPhone") || contains(userAgent, "iPad") {
		return "iOS"
	}
	if contains(userAgent, "Android") {
		return "Android"
	}
	return "Unknown"
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr))
}
