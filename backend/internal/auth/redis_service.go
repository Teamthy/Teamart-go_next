package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/teamart/commerce-api/pkg/logger"
)

// RedisService provides Redis integration for auth caching, rate limiting, and OTP
type RedisService struct {
	client *redis.Client
	logger *logger.Logger
}

// NewRedisService creates a new Redis service
func NewRedisService(client *redis.Client, logger *logger.Logger) *RedisService {
	return &RedisService{
		client: client,
		logger: logger,
	}
}

// ===== OTP Redis Operations =====

// StoreOTP stores an OTP in Redis with TTL
func (rs *RedisService) StoreOTP(ctx context.Context, email string, otpHash string, ttl time.Duration) error {
	key := fmt.Sprintf("otp:%s", email)

	err := rs.client.Set(ctx, key, otpHash, ttl).Err()
	if err != nil {
		rs.logger.Errorf("failed to store OTP: %v", err)
		return err
	}

	rs.logger.Debugf("OTP stored for %s with TTL %v", email, ttl)
	return nil
}

// GetOTP retrieves an OTP from Redis
func (rs *RedisService) GetOTP(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("otp:%s", email)

	val, err := rs.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("OTP not found or expired")
	}
	if err != nil {
		rs.logger.Errorf("failed to get OTP: %v", err)
		return "", err
	}

	return val, nil
}

// DeleteOTP deletes an OTP from Redis
func (rs *RedisService) DeleteOTP(ctx context.Context, email string) error {
	key := fmt.Sprintf("otp:%s", email)

	err := rs.client.Del(ctx, key).Err()
	if err != nil {
		rs.logger.Errorf("failed to delete OTP: %v", err)
		return err
	}

	return nil
}

// ===== OTP Attempt Tracking =====

// IncrementOTPAttempts increments OTP attempt counter
func (rs *RedisService) IncrementOTPAttempts(ctx context.Context, email string, maxAttempts int32, window time.Duration) (int32, error) {
	key := fmt.Sprintf("otp_attempts:%s", email)

	val, err := rs.client.Incr(ctx, key).Result()
	if err != nil {
		rs.logger.Errorf("failed to increment OTP attempts: %v", err)
		return 0, err
	}

	// Set expiration on first increment
	if val == 1 {
		rs.client.Expire(ctx, key, window)
	}

	attempts := int32(val)
	if attempts >= maxAttempts {
		rs.logger.Warnf("OTP attempts exceeded for %s: %d/%d", email, attempts, maxAttempts)
	}

	return attempts, nil
}

// GetOTPAttempts gets the current OTP attempt count
func (rs *RedisService) GetOTPAttempts(ctx context.Context, email string) (int32, error) {
	key := fmt.Sprintf("otp_attempts:%s", email)

	val, err := rs.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		rs.logger.Errorf("failed to get OTP attempts: %v", err)
		return 0, err
	}

	return int32(val), nil
}

// ResetOTPAttempts resets the OTP attempt counter
func (rs *RedisService) ResetOTPAttempts(ctx context.Context, email string) error {
	key := fmt.Sprintf("otp_attempts:%s", email)

	err := rs.client.Del(ctx, key).Err()
	if err != nil {
		rs.logger.Errorf("failed to reset OTP attempts: %v", err)
		return err
	}

	return nil
}

// ===== Session Caching =====

// CacheSession stores a session in Redis for fast retrieval
func (rs *RedisService) CacheSession(ctx context.Context, session *Session, ttl time.Duration) error {
	sessionKey := fmt.Sprintf("session:%s", session.ID)
	userSessionsKey := fmt.Sprintf("user_sessions:%d", session.UserID)

	// Serialize session
	sessionData, err := json.Marshal(session)
	if err != nil {
		rs.logger.Errorf("failed to marshal session: %v", err)
		return err
	}

	// Store in pipeline for atomic operation
	pipe := rs.client.Pipeline()
	pipe.Set(ctx, sessionKey, sessionData, ttl)
	pipe.SAdd(ctx, userSessionsKey, session.ID)
	pipe.Expire(ctx, userSessionsKey, ttl)

	_, err = pipe.Exec(ctx)
	if err != nil {
		rs.logger.Errorf("failed to cache session: %v", err)
		return err
	}

	rs.logger.Debugf("session cached: %s for user %d", session.ID, session.UserID)
	return nil
}

// GetCachedSession retrieves a cached session
func (rs *RedisService) GetCachedSession(ctx context.Context, sessionID string) (*Session, error) {
	key := fmt.Sprintf("session:%s", sessionID)

	val, err := rs.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("session not found in cache")
	}
	if err != nil {
		rs.logger.Errorf("failed to get cached session: %v", err)
		return nil, err
	}

	session := &Session{}
	err = json.Unmarshal([]byte(val), session)
	if err != nil {
		rs.logger.Errorf("failed to unmarshal session: %v", err)
		return nil, err
	}

	return session, nil
}

// InvalidateSession removes a session from cache
func (rs *RedisService) InvalidateSession(ctx context.Context, sessionID string, userID int64) error {
	sessionKey := fmt.Sprintf("session:%s", sessionID)
	userSessionsKey := fmt.Sprintf("user_sessions:%d", userID)

	pipe := rs.client.Pipeline()
	pipe.Del(ctx, sessionKey)
	pipe.SRem(ctx, userSessionsKey, sessionID)

	_, err := pipe.Exec(ctx)
	if err != nil {
		rs.logger.Errorf("failed to invalidate session: %v", err)
		return err
	}

	rs.logger.Debugf("session invalidated: %s", sessionID)
	return nil
}

// ===== Token Blacklisting =====

// BlacklistToken adds a token JTI to the blacklist
func (rs *RedisService) BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", jti)

	err := rs.client.Set(ctx, key, "revoked", ttl).Err()
	if err != nil {
		rs.logger.Errorf("failed to blacklist token: %v", err)
		return err
	}

	return nil
}

// IsTokenBlacklisted checks if a token JTI is blacklisted
func (rs *RedisService) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", jti)

	val, err := rs.client.Exists(ctx, key).Result()
	if err != nil {
		rs.logger.Errorf("failed to check token blacklist: %v", err)
		return false, err
	}

	return val > 0, nil
}

// ===== Rate Limiting =====

// TrackLoginAttempt tracks a login attempt for rate limiting
func (rs *RedisService) TrackLoginAttempt(ctx context.Context, identifier string, window time.Duration) (int32, error) {
	key := fmt.Sprintf("login_attempts:%s", identifier)

	val, err := rs.client.Incr(ctx, key).Result()
	if err != nil {
		rs.logger.Errorf("failed to track login attempt: %v", err)
		return 0, err
	}

	// Set expiration on first increment
	if val == 1 {
		rs.client.Expire(ctx, key, window)
	}

	return int32(val), nil
}

// GetLoginAttempts gets the current login attempt count
func (rs *RedisService) GetLoginAttempts(ctx context.Context, identifier string) (int32, error) {
	key := fmt.Sprintf("login_attempts:%s", identifier)

	val, err := rs.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		rs.logger.Errorf("failed to get login attempts: %v", err)
		return 0, err
	}

	return int32(val), nil
}

// ResetLoginAttempts resets the login attempt counter
func (rs *RedisService) ResetLoginAttempts(ctx context.Context, identifier string) error {
	key := fmt.Sprintf("login_attempts:%s", identifier)

	err := rs.client.Del(ctx, key).Err()
	if err != nil {
		rs.logger.Errorf("failed to reset login attempts: %v", err)
		return err
	}

	return nil
}

// ===== IP-based Rate Limiting =====

// IsIPBlocked checks if an IP is blocked
func (rs *RedisService) IsIPBlocked(ctx context.Context, ipAddress string) (bool, error) {
	key := fmt.Sprintf("blocked_ip:%s", ipAddress)

	val, err := rs.client.Exists(ctx, key).Result()
	if err != nil {
		rs.logger.Errorf("failed to check IP block: %v", err)
		return false, err
	}

	return val > 0, nil
}

// BlockIP blocks an IP address temporarily
func (rs *RedisService) BlockIP(ctx context.Context, ipAddress string, duration time.Duration) error {
	key := fmt.Sprintf("blocked_ip:%s", ipAddress)

	err := rs.client.Set(ctx, key, "blocked", duration).Err()
	if err != nil {
		rs.logger.Errorf("failed to block IP: %v", err)
		return err
	}

	rs.logger.Warnf("IP blocked: %s for %v", ipAddress, duration)
	return nil
}

// ===== Device Trust Cache =====

// CacheDeviceTrust stores device trust information
func (rs *RedisService) CacheDeviceTrust(ctx context.Context, userID int64, deviceID string, trustLevel TrustLevel, ttl time.Duration) error {
	key := fmt.Sprintf("device_trust:%d:%s", userID, deviceID)

	err := rs.client.Set(ctx, key, string(trustLevel), ttl).Err()
	if err != nil {
		rs.logger.Errorf("failed to cache device trust: %v", err)
		return err
	}

	return nil
}

// GetDeviceTrust retrieves cached device trust information
func (rs *RedisService) GetDeviceTrust(ctx context.Context, userID int64, deviceID string) (TrustLevel, error) {
	key := fmt.Sprintf("device_trust:%d:%s", userID, deviceID)

	val, err := rs.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("device trust not found in cache")
	}
	if err != nil {
		rs.logger.Errorf("failed to get device trust: %v", err)
		return "", err
	}

	return TrustLevel(val), nil
}

// ===== Onboarding State Cache =====

// CacheOnboardingState stores temporary onboarding state
func (rs *RedisService) CacheOnboardingState(ctx context.Context, userID int64, state OnboardingState, data map[string]interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("onboarding:%d", userID)

	stateData := map[string]interface{}{
		"state": state,
		"data":  data,
	}

	jsonData, err := json.Marshal(stateData)
	if err != nil {
		return err
	}

	err = rs.client.Set(ctx, key, jsonData, ttl).Err()
	if err != nil {
		rs.logger.Errorf("failed to cache onboarding state: %v", err)
		return err
	}

	return nil
}

// GetOnboardingState retrieves cached onboarding state
func (rs *RedisService) GetOnboardingState(ctx context.Context, userID int64) (OnboardingState, map[string]interface{}, error) {
	key := fmt.Sprintf("onboarding:%d", userID)

	val, err := rs.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil, fmt.Errorf("onboarding state not found in cache")
	}
	if err != nil {
		rs.logger.Errorf("failed to get onboarding state: %v", err)
		return "", nil, err
	}

	stateData := map[string]interface{}{}
	err = json.Unmarshal([]byte(val), &stateData)
	if err != nil {
		return "", nil, err
	}

	state := OnboardingState(stateData["state"].(string))
	data := stateData["data"].(map[string]interface{})

	return state, data, nil
}

// ===== Health Check =====

// Health checks if Redis is available
func (rs *RedisService) Health(ctx context.Context) error {
	err := rs.client.Ping(ctx).Err()
	if err != nil {
		rs.logger.Errorf("Redis health check failed: %v", err)
		return err
	}

	rs.logger.Debugf("Redis health check passed")
	return nil
}

// Close closes the Redis connection
func (rs *RedisService) Close() error {
	return rs.client.Close()
}
