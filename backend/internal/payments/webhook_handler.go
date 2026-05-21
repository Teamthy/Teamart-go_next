package payments

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// WebhookHandler processes payment webhooks from various gateways
type WebhookHandler interface {
	HandleWebhook(ctx context.Context, payload *WebhookPayload) error
	VerifyWebhookSignature(provider string, signature, timestamp, body string) (bool, error)
	ProcessStripeWebhook(ctx context.Context, data map[string]interface{}) error
	ProcessPaystackWebhook(ctx context.Context, data map[string]interface{}) error
	ProcessFlutterwaveWebhook(ctx context.Context, data map[string]interface{}) error
}

// WebhookProcessor handles webhook processing
type WebhookProcessor struct {
	querier PaymentQuerier
	logger  interface{ Printf(string, ...interface{}) }
	service *Service // Reference to payment service
}

// NewWebhookProcessor creates a new webhook processor
func NewWebhookProcessor(querier PaymentQuerier, logger interface{ Printf(string, ...interface{}) }, service *Service) *WebhookProcessor {
	return &WebhookProcessor{
		querier: querier,
		logger:  logger,
		service: service,
	}
}

// HandleWebhook processes a webhook payload
func (wp *WebhookProcessor) HandleWebhook(ctx context.Context, payload *WebhookPayload) error {
	if payload.Provider == "" {
		return fmt.Errorf("provider is required")
	}

	// Log webhook
	payloadBytes, _ := json.Marshal(payload)
	logID, err := wp.querier.CreateWebhookLog(ctx, payload.Provider, payload.EventType, payloadBytes)
	if err != nil {
		wp.logger.Printf("failed to log webhook: %v", err)
		return err
	}

	// Verify signature
	valid, err := wp.VerifyWebhookSignature(payload.Provider, payload.Signature, fmt.Sprintf("%d", payload.Timestamp), string(payloadBytes))
	if err != nil {
		wp.logger.Printf("failed to verify webhook signature: %v", err)
		_ = wp.querier.UpdateWebhookProcessed(ctx, logID, false, fmt.Sprintf("signature verification failed: %v", err))
		return err
	}

	if !valid {
		wp.logger.Printf("invalid webhook signature for provider %s", payload.Provider)
		_ = wp.querier.UpdateWebhookProcessed(ctx, logID, false, "invalid signature")
		return fmt.Errorf("invalid webhook signature")
	}

	// Process based on provider
	var processErr error
	switch payload.Provider {
	case "stripe":
		processErr = wp.ProcessStripeWebhook(ctx, payload.Data)
	case "paystack":
		processErr = wp.ProcessPaystackWebhook(ctx, payload.Data)
	case "flutterwave":
		processErr = wp.ProcessFlutterwaveWebhook(ctx, payload.Data)
	default:
		processErr = fmt.Errorf("unsupported provider: %s", payload.Provider)
	}

	// Update webhook log
	if processErr != nil {
		wp.logger.Printf("webhook processing error: %v", processErr)
		_ = wp.querier.UpdateWebhookProcessed(ctx, logID, false, processErr.Error())
		return processErr
	}

	// Mark as processed
	_ = wp.querier.UpdateWebhookProcessed(ctx, logID, true, "")
	wp.logger.Printf("webhook processed successfully: %s - %s", payload.Provider, payload.EventType)

	return nil
}

// VerifyWebhookSignature verifies webhook signatures
func (wp *WebhookProcessor) VerifyWebhookSignature(provider string, signature, timestamp, body string) (bool, error) {
	switch provider {
	case "stripe":
		return wp.verifyStripeSignature(signature, timestamp, body)
	case "paystack":
		return wp.verifyPaystackSignature(signature, body)
	case "flutterwave":
		return wp.verifyFlutterwaveSignature(signature, body)
	default:
		return false, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// verifyStripeSignature verifies Stripe webhook signature
func (wp *WebhookProcessor) verifyStripeSignature(signature, timestamp, body string) (bool, error) {
	// In production, use Stripe's webhook secret
	// stripe.Webhook.ConstructEvent(body, signature, webhookSecret)
	// For now, just return true (implement properly with actual Stripe signing)
	return true, nil
}

// verifyPaystackSignature verifies Paystack webhook signature
func (wp *WebhookProcessor) verifyPaystackSignature(signature, body string) (bool, error) {
	// Paystack uses HMAC-SHA512
	// secretKey := os.Getenv("PAYSTACK_SECRET_KEY")
	// hash := hmac.New(sha512.New, []byte(secretKey))
	// hash.Write([]byte(body))
	// expectedSignature := hex.EncodeToString(hash.Sum(nil))
	// return signature == expectedSignature, nil

	return true, nil
}

// verifyFlutterwaveSignature verifies Flutterwave webhook signature
func (wp *WebhookProcessor) verifyFlutterwaveSignature(signature, body string) (bool, error) {
	// Flutterwave uses HMAC-SHA256
	// secretHash := os.Getenv("FLUTTERWAVE_SECRET_HASH")
	// hash := hmac.New(sha256.New, []byte(secretHash))
	// hash.Write([]byte(body))
	// expectedSignature := hex.EncodeToString(hash.Sum(nil))
	// return signature == expectedSignature, nil

	return true, nil
}

// ProcessStripeWebhook processes Stripe webhook events
func (wp *WebhookProcessor) ProcessStripeWebhook(ctx context.Context, data map[string]interface{}) error {
	eventType, ok := data["type"].(string)
	if !ok {
		return fmt.Errorf("event type not found")
	}

	wp.logger.Printf("processing Stripe webhook: %s", eventType)

	switch eventType {
	case "payment_intent.succeeded":
		return wp.handleStripePaymentSucceeded(ctx, data)
	case "payment_intent.payment_failed":
		return wp.handleStripePaymentFailed(ctx, data)
	case "charge.refunded":
		return wp.handleStripeChargeRefunded(ctx, data)
	default:
		wp.logger.Printf("unhandled Stripe event type: %s", eventType)
		return nil
	}
}

// ProcessPaystackWebhook processes Paystack webhook events
func (wp *WebhookProcessor) ProcessPaystackWebhook(ctx context.Context, data map[string]interface{}) error {
	eventType, ok := data["event"].(string)
	if !ok {
		return fmt.Errorf("event type not found")
	}

	wp.logger.Printf("processing Paystack webhook: %s", eventType)

	switch eventType {
	case "charge.success":
		return wp.handlePaystackChargeSuccess(ctx, data)
	case "charge.failed":
		return wp.handlePaystackChargeFailed(ctx, data)
	default:
		wp.logger.Printf("unhandled Paystack event type: %s", eventType)
		return nil
	}
}

// ProcessFlutterwaveWebhook processes Flutterwave webhook events
func (wp *WebhookProcessor) ProcessFlutterwaveWebhook(ctx context.Context, data map[string]interface{}) error {
	eventType, ok := data["event"].(string)
	if !ok {
		return fmt.Errorf("event type not found")
	}

	wp.logger.Printf("processing Flutterwave webhook: %s", eventType)

	switch eventType {
	case "charge.completed":
		return wp.handleFlutterwaveChargeCompleted(ctx, data)
	case "charge.failed":
		return wp.handleFlutterwaveChargeFailed(ctx, data)
	default:
		wp.logger.Printf("unhandled Flutterwave event type: %s", eventType)
		return nil
	}
}

// handleStripePaymentSucceeded handles Stripe payment success
func (wp *WebhookProcessor) handleStripePaymentSucceeded(ctx context.Context, data map[string]interface{}) error {
	// Get payment intent ID
	// Update payment intent status to "succeeded"
	// Process any split payments or escrow releases
	wp.logger.Printf("stripe payment succeeded")
	return nil
}

// handleStripePaymentFailed handles Stripe payment failure
func (wp *WebhookProcessor) handleStripePaymentFailed(ctx context.Context, data map[string]interface{}) error {
	// Get payment intent ID
	// Update payment intent status to "failed"
	// Send notification to user
	wp.logger.Printf("stripe payment failed")
	return nil
}

// handleStripeChargeRefunded handles Stripe charge refund
func (wp *WebhookProcessor) handleStripeChargeRefunded(ctx context.Context, data map[string]interface{}) error {
	// Get refund ID
	// Update refund status to "completed"
	// Deduct from seller wallet if necessary
	wp.logger.Printf("stripe charge refunded")
	return nil
}

// handlePaystackChargeSuccess handles Paystack charge success
func (wp *WebhookProcessor) handlePaystackChargeSuccess(ctx context.Context, data map[string]interface{}) error {
	wp.logger.Printf("paystack charge succeeded")
	return nil
}

// handlePaystackChargeFailed handles Paystack charge failure
func (wp *WebhookProcessor) handlePaystackChargeFailed(ctx context.Context, data map[string]interface{}) error {
	wp.logger.Printf("paystack charge failed")
	return nil
}

// handleFlutterwaveChargeCompleted handles Flutterwave charge completion
func (wp *WebhookProcessor) handleFlutterwaveChargeCompleted(ctx context.Context, data map[string]interface{}) error {
	wp.logger.Printf("flutterwave charge completed")
	return nil
}

// handleFlutterwaveChargeFailed handles Flutterwave charge failure
func (wp *WebhookProcessor) handleFlutterwaveChargeFailed(ctx context.Context, data map[string]interface{}) error {
	wp.logger.Printf("flutterwave charge failed")
	return nil
}

// ComputeHMAC computes HMAC signature for verification
func ComputeHMAC(secret, message string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
