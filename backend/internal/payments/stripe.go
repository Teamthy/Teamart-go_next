package payments

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// StripeGatewayImpl implements PaymentGatewayProvider for Stripe
type StripeGatewayImpl struct {
	secretKey  string
	publicKey  string
	httpClient *http.Client
	logger     interface{ Printf(string, ...interface{}) }
}

// NewStripeGatewayImpl creates a new Stripe gateway implementation
func NewStripeGatewayImpl(secretKey, publicKey string, logger interface{ Printf(string, ...interface{}) }) *StripeGatewayImpl {
	return &StripeGatewayImpl{
		secretKey: secretKey,
		publicKey: publicKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// CreatePaymentIntent creates a Stripe payment intent
func (sg *StripeGatewayImpl) CreatePaymentIntent(ctx context.Context, input *CreatePaymentIntentInput) (*PaymentIntentResult, error) {
	if input.Amount <= 0 {
		return nil, fmt.Errorf("invalid amount: %.2f", input.Amount)
	}

	// Convert amount to cents
	amountCents := int64(input.Amount * 100)

	// Prepare request to Stripe API
	// POST https://api.stripe.com/v1/payment_intents
	payload := fmt.Sprintf(
		"amount=%d&currency=%s&metadata[order_id]=%d&metadata[user_id]=%d",
		amountCents,
		strings.ToLower(input.Currency),
		input.OrderID,
		input.UserID,
	)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.stripe.com/v1/payment_intents", strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(sg.secretKey, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := sg.httpClient.Do(req)
	if err != nil {
		sg.logger.Printf("stripe API error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("stripe error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// Parse response
	var stripeResponse map[string]interface{}
	if err := json.Unmarshal(body, &stripeResponse); err != nil {
		return nil, err
	}

	intentID, _ := stripeResponse["id"].(string)
	clientSecret, _ := stripeResponse["client_secret"].(string)

	result := &PaymentIntentResult{
		Provider:             "stripe",
		ProviderIntentID:     intentID,
		ProviderClientSecret: &clientSecret,
		Status:               "pending",
	}

	sg.logger.Printf("stripe payment intent created: %s", intentID)
	return result, nil
}

// ProcessPayment confirms and processes a Stripe payment
func (sg *StripeGatewayImpl) ProcessPayment(ctx context.Context, input *ProcessPaymentInput) (*PaymentResult, error) {
	// POST https://api.stripe.com/v1/payment_intents/{id}/confirm

	payload := fmt.Sprintf("payment_method=%s", input.ProviderToken)

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("https://api.stripe.com/v1/payment_intents/%d/confirm", input.PaymentIntentID), strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(sg.secretKey, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := sg.httpClient.Do(req)
	if err != nil {
		sg.logger.Printf("stripe payment confirmation error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var stripeResponse map[string]interface{}
	if err := json.Unmarshal(body, &stripeResponse); err != nil {
		return nil, err
	}

	status, _ := stripeResponse["status"].(string)
	success := status == "succeeded"

	result := &PaymentResult{
		Success: success,
		Status:  status,
		Message: "Payment processed",
	}

	if !success {
		result.Message = "Payment declined"
	}

	sg.logger.Printf("stripe payment processed: status=%s", status)
	return result, nil
}

// RefundPayment refunds a Stripe payment
func (sg *StripeGatewayImpl) RefundPayment(ctx context.Context, input *RefundPaymentInput) (*RefundResult, error) {
	// POST https://api.stripe.com/v1/refunds

	amountCents := int64(input.Amount * 100)
	payload := fmt.Sprintf("amount=%d&charge_id=%d&reason=%s", amountCents, input.PaymentIntentID, input.Reason)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.stripe.com/v1/refunds", strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(sg.secretKey, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := sg.httpClient.Do(req)
	if err != nil {
		sg.logger.Printf("stripe refund error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var stripeResponse map[string]interface{}
	if err := json.Unmarshal(body, &stripeResponse); err != nil {
		return nil, err
	}

	refundID, _ := stripeResponse["id"].(string)
	status, _ := stripeResponse["status"].(string)

	result := &RefundResult{
		Success:          status == "succeeded",
		ProviderRefundID: refundID,
		Status:           status,
		Amount:           input.Amount,
		Message:          "Refund processed",
	}

	sg.logger.Printf("stripe refund processed: %s", refundID)
	return result, nil
}

// VerifyWebhookSignature verifies Stripe webhook signature
func (sg *StripeGatewayImpl) VerifyWebhookSignature(signature, timestamp, payload string) (bool, error) {
	// Stripe uses HMAC-SHA256 for webhook signing
	// Verify using the webhook secret

	expected := ComputeHMAC(sg.secretKey, timestamp+"."+payload)
	return signature == expected, nil
}

// GetPaymentStatus retrieves payment status from Stripe
func (sg *StripeGatewayImpl) GetPaymentStatus(ctx context.Context, providerPaymentID string) (*PaymentStatus, error) {
	// GET https://api.stripe.com/v1/payment_intents/{id}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.stripe.com/v1/payment_intents/%s", providerPaymentID), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(sg.secretKey, "")

	resp, err := sg.httpClient.Do(req)
	if err != nil {
		sg.logger.Printf("stripe status check error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var stripeResponse map[string]interface{}
	if err := json.Unmarshal(body, &stripeResponse); err != nil {
		return nil, err
	}

	status, _ := stripeResponse["status"].(string)
	amount, _ := stripeResponse["amount"].(float64)
	currency, _ := stripeResponse["currency"].(string)

	paymentStatus := &PaymentStatus{
		ProviderPaymentID: providerPaymentID,
		Status:            status,
		Amount:            amount / 100, // Convert from cents
		Currency:          currency,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	return paymentStatus, nil
}
