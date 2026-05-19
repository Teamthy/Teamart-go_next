package auth

import (
	"context"
	"testing"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

// TestCreateSession_Success tests successful session creation
func TestCreateSession_Success(t *testing.T) {
	// Setup
	config := &AuthConfig{
		SessionTTL:         24 * time.Hour,
		SessionIdleTimeout: 30 * time.Minute,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)

	input := &CreateSessionInput{
		UserID:     1,
		DeviceID:   "device-123",
		DeviceName: "iPhone 14",
		DeviceType: "mobile",
		IPAddress:  "192.168.1.1",
		UserAgent:  "Mozilla/5.0",
	}

	// Execute
	output, err := service.CreateSession(context.Background(), input)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Session == nil {
		t.Fatal("expected session, got nil")
	}
	if output.Session.UserID != 1 {
		t.Errorf("expected user ID 1, got %d", output.Session.UserID)
	}
	if output.Session.DeviceID != "device-123" {
		t.Errorf("expected device ID device-123, got %s", output.Session.DeviceID)
	}
	if output.Session.IsRevoked {
		t.Error("expected session not to be revoked")
	}
	if output.Session.TrustLevel != TrustLevelUntrusted {
		t.Errorf("expected trust level untrusted, got %s", output.Session.TrustLevel)
	}
}

// TestCreateSession_MissingUserID tests session creation with missing user ID
func TestCreateSession_MissingUserID(t *testing.T) {
	config := &AuthConfig{SessionTTL: 24 * time.Hour}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)

	input := &CreateSessionInput{
		UserID:    0, // Invalid
		DeviceID:  "device-123",
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}

	_, err := service.CreateSession(context.Background(), input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestValidateSession_Success tests successful session validation
func TestValidateSession_Success(t *testing.T) {
	// Setup
	config := &AuthConfig{
		SessionTTL:         24 * time.Hour,
		SessionIdleTimeout: 30 * time.Minute,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)

	// Create a session first
	createInput := &CreateSessionInput{
		UserID:    1,
		DeviceID:  "device-123",
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}
	output, _ := service.CreateSession(context.Background(), createInput)
	session := output.Session

	// Validate the session
	validateInput := &ValidateSessionInput{
		SessionID: session.ID,
		UserID:    1,
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}

	validateOutput, err := service.ValidateSession(context.Background(), validateInput)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !validateOutput.IsValid {
		t.Error("expected session to be valid")
	}
	if validateOutput.Session == nil {
		t.Fatal("expected session, got nil")
	}
}

// TestValidateSession_Expired tests validation of expired session
func TestValidateSession_Expired(t *testing.T) {
	// Setup with very short TTL
	config := &AuthConfig{
		SessionTTL:         1 * time.Millisecond, // Expire immediately
		SessionIdleTimeout: 30 * time.Minute,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)

	// Create a session
	createInput := &CreateSessionInput{
		UserID:    1,
		DeviceID:  "device-123",
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}
	output, _ := service.CreateSession(context.Background(), createInput)
	session := output.Session

	// Wait for session to expire
	time.Sleep(10 * time.Millisecond)

	// Try to validate
	validateInput := &ValidateSessionInput{
		SessionID: session.ID,
		UserID:    1,
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}

	validateOutput, _ := service.ValidateSession(context.Background(), validateInput)

	if validateOutput.IsValid {
		t.Error("expected session to be invalid (expired)")
	}
}

// TestValidateSession_Revoked tests validation of revoked session
func TestValidateSession_Revoked(t *testing.T) {
	config := &AuthConfig{
		SessionTTL:         24 * time.Hour,
		SessionIdleTimeout: 30 * time.Minute,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)

	// Create and revoke a session
	createInput := &CreateSessionInput{
		UserID:    1,
		DeviceID:  "device-123",
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}
	output, _ := service.CreateSession(context.Background(), createInput)
	session := output.Session

	// Revoke the session
	revokeInput := &RevokeSessionInput{
		SessionID: session.ID,
		UserID:    1,
		Reason:    "test_revocation",
	}
	service.RevokeSession(context.Background(), revokeInput)

	// Try to validate revoked session
	validateInput := &ValidateSessionInput{
		SessionID: session.ID,
		UserID:    1,
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}

	validateOutput, _ := service.ValidateSession(context.Background(), validateInput)

	if validateOutput.IsValid {
		t.Error("expected revoked session to be invalid")
	}
}

// TestValidateSession_IPChange tests session validation with IP change
func TestValidateSession_IPChange(t *testing.T) {
	config := &AuthConfig{
		SessionTTL:         24 * time.Hour,
		SessionIdleTimeout: 30 * time.Minute,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)

	// Create session
	createInput := &CreateSessionInput{
		UserID:    1,
		DeviceID:  "device-123",
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}
	output, _ := service.CreateSession(context.Background(), createInput)
	session := output.Session

	// Validate from different IP
	validateInput := &ValidateSessionInput{
		SessionID: session.ID,
		UserID:    1,
		IPAddress: "192.168.1.2", // Different IP
		UserAgent: "Mozilla/5.0",
	}

	validateOutput, _ := service.ValidateSession(context.Background(), validateInput)

	if !validateOutput.IsValid {
		t.Error("expected session to still be valid even with IP change")
	}
	if !validateOutput.RequiresMFA {
		t.Error("expected MFA requirement due to IP change")
	}
}

// TestUpdateSessionActivity updates session last activity
func TestUpdateSessionActivity(t *testing.T) {
	config := &AuthConfig{
		SessionTTL:         24 * time.Hour,
		SessionIdleTimeout: 30 * time.Minute,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)

	// Create session
	createInput := &CreateSessionInput{
		UserID:    1,
		DeviceID:  "device-123",
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}
	output, _ := service.CreateSession(context.Background(), createInput)
	originalLastActivity := output.Session.LastActivityAt

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Update activity
	updateInput := &UpdateSessionActivityInput{
		SessionID: output.Session.ID,
		UserID:    1,
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}
	err := service.UpdateSessionActivity(context.Background(), updateInput)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Fetch updated session
	updated, _ := repo.GetSession(context.Background(), output.Session.ID)
	if updated.LastActivityAt == originalLastActivity {
		t.Error("expected last activity to be updated")
	}
	if updated.LastActivityAt.Before(originalLastActivity) {
		t.Error("last activity should not go backwards")
	}
}

// TestRevokeSession revokes a single session
func TestRevokeSession(t *testing.T) {
	config := &AuthConfig{
		SessionTTL:         24 * time.Hour,
		SessionIdleTimeout: 30 * time.Minute,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)

	// Create session
	createInput := &CreateSessionInput{
		UserID:    1,
		DeviceID:  "device-123",
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}
	output, _ := service.CreateSession(context.Background(), createInput)

	// Revoke session
	revokeInput := &RevokeSessionInput{
		SessionID: output.Session.ID,
		UserID:    1,
		Reason:    "logout",
	}
	err := service.RevokeSession(context.Background(), revokeInput)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify revocation
	revoked, _ := repo.GetSession(context.Background(), output.Session.ID)
	if !revoked.IsRevoked {
		t.Error("expected session to be revoked")
	}
	if revoked.RevokeReason != "logout" {
		t.Errorf("expected revoke reason 'logout', got %s", revoked.RevokeReason)
	}
}

// TestRevokeAllUserSessions revokes all sessions for a user
func TestRevokeAllUserSessions(t *testing.T) {
	config := &AuthConfig{
		SessionTTL:         24 * time.Hour,
		SessionIdleTimeout: 30 * time.Minute,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)

	// Create multiple sessions for the same user
	input1 := &CreateSessionInput{
		UserID:    1,
		DeviceID:  "device-1",
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}
	output1, _ := service.CreateSession(context.Background(), input1)

	input2 := &CreateSessionInput{
		UserID:    1,
		DeviceID:  "device-2",
		IPAddress: "192.168.1.2",
		UserAgent: "Safari/5.0",
	}
	output2, _ := service.CreateSession(context.Background(), input2)

	// Revoke all sessions except one
	revokeAllInput := &RevokeAllUserSessionsInput{
		UserID:          1,
		ExceptSessionID: output1.Session.ID,
		Reason:          "security_breach",
	}
	err := service.RevokeAllUserSessions(context.Background(), revokeAllInput)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify: first session should still be active, second should be revoked
	session1, _ := repo.GetSession(context.Background(), output1.Session.ID)
	session2, _ := repo.GetSession(context.Background(), output2.Session.ID)

	if session1.IsRevoked {
		t.Error("expected first session to NOT be revoked")
	}
	if !session2.IsRevoked {
		t.Error("expected second session to be revoked")
	}
}

// TestGetUserActiveSessions lists active sessions
func TestGetUserActiveSessions(t *testing.T) {
	config := &AuthConfig{
		SessionTTL:         24 * time.Hour,
		SessionIdleTimeout: 30 * time.Minute,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)

	// Create multiple sessions
	for i := 1; i <= 3; i++ {
		input := &CreateSessionInput{
			UserID:    1,
			DeviceID:  "device-" + string(rune(i)),
			IPAddress: "192.168.1.1",
			UserAgent: "Mozilla/5.0",
		}
		service.CreateSession(context.Background(), input)
	}

	// Get active sessions
	output, err := service.GetUserActiveSessions(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Sessions) != 3 {
		t.Errorf("expected 3 sessions, got %d", len(output.Sessions))
	}
}

// TestVerifyDevice_Trusted tests device verification for trusted device
func TestVerifyDevice_Trusted(t *testing.T) {
	config := &AuthConfig{
		SessionTTL:         24 * time.Hour,
		SessionIdleTimeout: 30 * time.Minute,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)

	// Create session with a device
	createInput := &CreateSessionInput{
		UserID:    1,
		DeviceID:  "device-123",
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}
	service.CreateSession(context.Background(), createInput)

	// Trust the device
	trustInput := &TrustDeviceInput{
		UserID:     1,
		DeviceID:   "device-123",
		DeviceName: "My Phone",
	}
	service.TrustDevice(context.Background(), trustInput)

	// Verify the device
	verifyInput := &VerifyDeviceInput{
		UserID:   1,
		DeviceID: "device-123",
	}
	verifyOutput, _ := service.VerifyDevice(context.Background(), verifyInput)

	if !verifyOutput.IsTrusted {
		t.Error("expected device to be trusted")
	}
	if verifyOutput.TrustLevel != TrustLevelTrusted {
		t.Errorf("expected trust level trusted, got %s", verifyOutput.TrustLevel)
	}
}

// TestVerifyDevice_Untrusted tests device verification for unknown device
func TestVerifyDevice_Untrusted(t *testing.T) {
	config := &AuthConfig{
		SessionTTL:         24 * time.Hour,
		SessionIdleTimeout: 30 * time.Minute,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)

	// Try to verify device without creating any sessions first
	verifyInput := &VerifyDeviceInput{
		UserID:   1,
		DeviceID: "device-unknown",
	}
	verifyOutput, _ := service.VerifyDevice(context.Background(), verifyInput)

	if verifyOutput.IsTrusted {
		t.Error("expected device to be untrusted")
	}
	if verifyOutput.TrustLevel != TrustLevelUntrusted {
		t.Errorf("expected trust level untrusted, got %s", verifyOutput.TrustLevel)
	}
}

// TestMarkDeviceSuspicious marks device as suspicious
func TestMarkDeviceSuspicious(t *testing.T) {
	config := &AuthConfig{
		SessionTTL:         24 * time.Hour,
		SessionIdleTimeout: 30 * time.Minute,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)

	// Create and trust a session
	createInput := &CreateSessionInput{
		UserID:    1,
		DeviceID:  "device-123",
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}
	output, _ := service.CreateSession(context.Background(), createInput)

	trustInput := &TrustDeviceInput{
		UserID:   1,
		DeviceID: "device-123",
	}
	service.TrustDevice(context.Background(), trustInput)

	// Mark as suspicious
	err := service.MarkDeviceSuspicious(context.Background(), 1, "device-123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify
	updated, _ := repo.GetSession(context.Background(), output.Session.ID)
	if updated.TrustLevel != TrustLevelPartial {
		t.Errorf("expected trust level partial, got %s", updated.TrustLevel)
	}
	if !updated.RequiresMFA {
		t.Error("expected MFA requirement")
	}
}

// TestSessionLimit enforces maximum concurrent sessions
func TestSessionLimit(t *testing.T) {
	config := &AuthConfig{
		SessionTTL:         24 * time.Hour,
		SessionIdleTimeout: 30 * time.Minute,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)
	service.maxSessions = 2 // Set limit to 2

	// Create 3 sessions - 3rd should cause revocation of oldest
	for i := 1; i <= 3; i++ {
		input := &CreateSessionInput{
			UserID:    1,
			DeviceID:  "device-" + string(rune(i)),
			IPAddress: "192.168.1.1",
			UserAgent: "Mozilla/5.0",
		}
		service.CreateSession(context.Background(), input)
		time.Sleep(1 * time.Millisecond) // Small delay to ensure different timestamps
	}

	// Check active sessions
	sessions, _ := repo.GetUserSessions(context.Background(), 1)
	activeSessions := 0
	for _, s := range sessions {
		if !s.IsRevoked {
			activeSessions++
		}
	}

	if activeSessions > service.maxSessions {
		t.Errorf("expected max %d active sessions, got %d", service.maxSessions, activeSessions)
	}
}

// TestCleanupExpiredSessions removes old sessions
func TestCleanupExpiredSessions(t *testing.T) {
	config := &AuthConfig{
		SessionTTL:         1 * time.Millisecond, // Very short
		SessionIdleTimeout: 30 * time.Minute,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)
	service := NewSessionService(config, mockLogger, repo)

	// Create session
	createInput := &CreateSessionInput{
		UserID:    1,
		DeviceID:  "device-1",
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0",
	}
	service.CreateSession(context.Background(), createInput)

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Cleanup
	err := service.CleanupExpiredSessions(context.Background(), 1*time.Hour)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestSessionRepository_InMemory tests in-memory repository directly
func TestSessionRepository_InMemory(t *testing.T) {
	mockLogger := logger.NewLogger("test", false)
	repo := NewSessionRepositoryMemory(mockLogger)

	session := &Session{
		ID:        "test-session-1",
		UserID:    1,
		DeviceID:  "device-1",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create
	err := repo.CreateSession(context.Background(), session)
	if err != nil {
		t.Fatalf("expected no error creating session, got %v", err)
	}

	// Get
	retrieved, err := repo.GetSession(context.Background(), "test-session-1")
	if err != nil {
		t.Fatalf("expected no error getting session, got %v", err)
	}
	if retrieved.ID != "test-session-1" {
		t.Errorf("expected session ID test-session-1, got %s", retrieved.ID)
	}

	// Session exists
	exists, _ := repo.SessionExists(context.Background(), "test-session-1")
	if !exists {
		t.Error("expected session to exist")
	}

	// Revoke
	repo.RevokeSession(context.Background(), "test-session-1", "test_revoke")
	revoked, _ := repo.GetSession(context.Background(), "test-session-1")
	if !revoked.IsRevoked {
		t.Error("expected session to be revoked")
	}

	// Should not exist after revocation
	exists, _ = repo.SessionExists(context.Background(), "test-session-1")
	if exists {
		t.Error("expected revoked session to not exist")
	}
}
