package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

// FraudDetectionService detects and prevents fraudulent activity
type FraudDetectionService struct {
	config       *AuthConfig
	logger       *logger.Logger
	redisService *RedisService
}

// NewFraudDetectionService creates a new fraud detection service
func NewFraudDetectionService(
	config *AuthConfig,
	logger *logger.Logger,
	redisService *RedisService,
) *FraudDetectionService {
	return &FraudDetectionService{
		config:       config,
		logger:       logger,
		redisService: redisService,
	}
}

// FraudScore represents the fraud risk level
type FraudScore int

const (
	FraudScoreLow      FraudScore = 0
	FraudScoreMedium   FraudScore = 1
	FraudScoreHigh     FraudScore = 2
	FraudScoreCritical FraudScore = 3
)

// FraudIndicator represents a detected fraud indicator
type FraudIndicator struct {
	Type        string
	Severity    FraudScore
	Description string
	Timestamp   time.Time
}

// LoginFraudCheckInput represents input for login fraud check
type LoginFraudCheckInput struct {
	UserID    int64
	Email     string
	IPAddress string
	UserAgent string
	Timestamp time.Time
}

// LoginFraudCheckOutput represents the result of fraud check
type LoginFraudCheckOutput struct {
	IsFraudulent         bool
	FraudScore           FraudScore
	Indicators           []FraudIndicator
	RequiresMFA          bool
	RequiresVerification bool
	BlockLogin           bool
	Message              string
}

// CheckLoginFraud performs comprehensive fraud check for login attempts
func (fs *FraudDetectionService) CheckLoginFraud(ctx context.Context, input *LoginFraudCheckInput) (*LoginFraudCheckOutput, error) {
	fs.logger.Debugf("performing fraud check for user %d from IP %s", input.UserID, input.IPAddress)

	output := &LoginFraudCheckOutput{
		IsFraudulent:         false,
		FraudScore:           FraudScoreLow,
		Indicators:           []FraudIndicator{},
		RequiresMFA:          false,
		RequiresVerification: false,
		BlockLogin:           false,
	}

	// Check 1: IP-based rate limiting
	attempts, err := fs.redisService.GetLoginAttempts(ctx, input.IPAddress)
	if err != nil {
		fs.logger.Warnf("failed to check IP attempts: %v", err)
	} else if attempts > 10 {
		output.Indicators = append(output.Indicators, FraudIndicator{
			Type:        "excessive_ip_attempts",
			Severity:    FraudScoreCritical,
			Description: fmt.Sprintf("IP has %d failed login attempts", attempts),
			Timestamp:   time.Now(),
		})
		output.BlockLogin = true
		output.FraudScore = FraudScoreCritical
	}

	// Check 2: Email-based rate limiting
	emailAttempts, err := fs.redisService.GetLoginAttempts(ctx, input.Email)
	if err != nil {
		fs.logger.Warnf("failed to check email attempts: %v", err)
	} else if emailAttempts > 5 {
		output.Indicators = append(output.Indicators, FraudIndicator{
			Type:        "excessive_email_attempts",
			Severity:    FraudScoreMedium,
			Description: fmt.Sprintf("Email has %d failed login attempts", emailAttempts),
			Timestamp:   time.Now(),
		})
		output.RequiresMFA = true
		output.FraudScore = max(output.FraudScore, FraudScoreMedium)
	}

	// Check 3: Geographical impossibility
	isGeographicallyImpossible := fs.checkGeographicalImpossibility(ctx, input.UserID, input.IPAddress)
	if isGeographicallyImpossible {
		output.Indicators = append(output.Indicators, FraudIndicator{
			Type:        "geographical_impossibility",
			Severity:    FraudScoreHigh,
			Description: "Login from geographically impossible location",
			Timestamp:   time.Now(),
		})
		output.RequiresMFA = true
		output.RequiresVerification = true
		output.FraudScore = max(output.FraudScore, FraudScoreHigh)
	}

	// Check 4: New device/IP detection
	isNewDevice := fs.checkNewDevice(ctx, input.UserID, input.IPAddress)
	if isNewDevice {
		output.Indicators = append(output.Indicators, FraudIndicator{
			Type:        "new_device",
			Severity:    FraudScoreLow,
			Description: "Login from new device or location",
			Timestamp:   time.Now(),
		})
		output.FraudScore = max(output.FraudScore, FraudScoreLow)
	}

	// Check 5: Credential stuffing patterns
	isCredentialStuffing := fs.checkCredentialStuffing(ctx, input.Email, input.IPAddress)
	if isCredentialStuffing {
		output.Indicators = append(output.Indicators, FraudIndicator{
			Type:        "credential_stuffing",
			Severity:    FraudScoreCritical,
			Description: "Potential credential stuffing attack detected",
			Timestamp:   time.Now(),
		})
		output.BlockLogin = true
		output.FraudScore = FraudScoreCritical
	}

	// Set output flags
	if output.FraudScore >= FraudScoreMedium {
		output.IsFraudulent = true
		output.RequiresMFA = true
	}

	if output.BlockLogin {
		output.Message = "Login blocked due to suspicious activity. Please try again later."
		fs.logger.Warnf("login blocked for user %d due to fraud detection", input.UserID)
	} else if output.IsFraudulent {
		output.Message = "Additional verification required due to unusual activity."
		fs.logger.Warnf("suspicious login detected for user %d", input.UserID)
	}

	return output, nil
}

// CheckSignupFraud performs fraud check for new signups
func (fs *FraudDetectionService) CheckSignupFraud(ctx context.Context, email string, ipAddress string) (*LoginFraudCheckOutput, error) {
	output := &LoginFraudCheckOutput{
		IsFraudulent: false,
		FraudScore:   FraudScoreLow,
		Indicators:   []FraudIndicator{},
	}

	// Check 1: Bot signup detection - too many signups from same IP
	recentSignups, err := fs.redisService.GetLoginAttempts(ctx, "signup:"+ipAddress)
	if err == nil && recentSignups > 10 {
		output.Indicators = append(output.Indicators, FraudIndicator{
			Type:        "bot_signup",
			Severity:    FraudScoreHigh,
			Description: fmt.Sprintf("Excessive signups from IP (%d)", recentSignups),
			Timestamp:   time.Now(),
		})
		output.FraudScore = FraudScoreHigh
		output.IsFraudulent = true
	}

	// Check 2: Disposable email detection
	if fs.isDisposableEmail(email) {
		output.Indicators = append(output.Indicators, FraudIndicator{
			Type:        "disposable_email",
			Severity:    FraudScoreMedium,
			Description: "Signup using disposable email address",
			Timestamp:   time.Now(),
		})
		output.FraudScore = max(output.FraudScore, FraudScoreMedium)
	}

	// Check 3: Recent IP-based account creation
	recentAccounts, err := fs.redisService.GetLoginAttempts(ctx, "accounts:"+ipAddress)
	if err == nil && recentAccounts > 5 {
		output.Indicators = append(output.Indicators, FraudIndicator{
			Type:        "account_farming",
			Severity:    FraudScoreMedium,
			Description: fmt.Sprintf("Multiple accounts from same IP (%d)", recentAccounts),
			Timestamp:   time.Now(),
		})
		output.FraudScore = max(output.FraudScore, FraudScoreMedium)
	}

	if output.FraudScore >= FraudScoreMedium {
		output.IsFraudulent = true
	}

	return output, nil
}

// TrackRefreshTokenPattern tracks refresh token patterns to detect abuse
func (fs *FraudDetectionService) TrackRefreshTokenPattern(ctx context.Context, userID int64) (bool, error) {
	key := fmt.Sprintf("refresh_pattern:%d", userID)

	count, err := fs.redisService.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// Set expiration on first increment (1 hour window)
	if count == 1 {
		fs.redisService.client.Expire(ctx, key, 1*time.Hour)
	}

	// If more than 30 refreshes in 1 hour, it's suspicious
	if count > 30 {
		fs.logger.Warnf("suspicious refresh pattern detected for user %d (%d refreshes)", userID, count)
		return true, nil
	}

	return false, nil
}

// ===== Helper Methods =====

func (fs *FraudDetectionService) checkGeographicalImpossibility(ctx context.Context, userID int64, currentIP string) bool {
	// This is a simplified check - in production, use IP geolocation API
	// and check if the distance traveled is physically impossible
	return false
}

func (fs *FraudDetectionService) checkNewDevice(ctx context.Context, userID int64, ipAddress string) bool {
	// Check if this IP/device combination has been seen before
	key := fmt.Sprintf("device_ip:%d:%s", userID, ipAddress)
	exists, err := fs.redisService.client.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return exists == 0
}

func (fs *FraudDetectionService) checkCredentialStuffing(ctx context.Context, email string, ipAddress string) bool {
	// Check if there are many failed attempts across different accounts from same IP
	key := fmt.Sprintf("attempted_emails:%s", ipAddress)
	count, err := fs.redisService.client.SCard(ctx, key).Result()
	if err != nil {
		return false
	}
	return count > 20 // More than 20 different email attempts = likely stuffing
}

func (fs *FraudDetectionService) isDisposableEmail(email string) bool {
	// Check against list of known disposable email domains
	disposableDomains := map[string]bool{
		"tempmail.com":      true,
		"guerrillamail.com": true,
		"mailinator.com":    true,
		"10minutemail.com":  true,
		"trashmail.com":     true,
		// Add more as needed
	}

	// Extract domain from email
	parts := len(email)
	atIndex := -1
	for i := 0; i < parts; i++ {
		if email[i] == '@' {
			atIndex = i
			break
		}
	}

	if atIndex == -1 {
		return false
	}

	domain := email[atIndex+1:]
	return disposableDomains[domain]
}

func max(a, b FraudScore) FraudScore {
	if a > b {
		return a
	}
	return b
}

// IsIPBlocked checks if an IP is blocked due to fraud
func (fs *FraudDetectionService) IsIPBlocked(ctx context.Context, ipAddress string) (bool, error) {
	return fs.redisService.IsIPBlocked(ctx, ipAddress)
}

// BlockIP blocks an IP due to fraud
func (fs *FraudDetectionService) BlockIP(ctx context.Context, ipAddress string, duration time.Duration) error {
	return fs.redisService.BlockIP(ctx, ipAddress, duration)
}

// GetFraudScoreDescription returns a human-readable description of fraud score
func GetFraudScoreDescription(score FraudScore) string {
	switch score {
	case FraudScoreLow:
		return "Low risk"
	case FraudScoreMedium:
		return "Medium risk - additional verification may be required"
	case FraudScoreHigh:
		return "High risk - requires multi-factor authentication"
	case FraudScoreCritical:
		return "Critical risk - login blocked"
	default:
		return "Unknown risk"
	}
}

// GetHTTPStatus returns appropriate HTTP status for fraud score
func GetHTTPStatus(score FraudScore) int {
	switch score {
	case FraudScoreLow:
		return http.StatusOK
	case FraudScoreMedium:
		return http.StatusAccepted // 202
	case FraudScoreHigh:
		return http.StatusForbidden
	case FraudScoreCritical:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
