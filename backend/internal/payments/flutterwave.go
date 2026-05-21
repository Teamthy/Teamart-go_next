package payments

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// FlutterwaveGatewayImpl implements PaymentGatewayProvider for Flutterwave
type FlutterwaveGatewayImpl struct {
	secretKey  string
	publicKey  string
	httpClient *http.Client
	logger     interface{ Printf(string, ...interface{}) }
}

// NewFlutterwaveGatewayImpl creates a new Flutterwave gateway implementation
func NewFlutterwaveGatewayImpl(secretKey, publicKey string, logger interface{ Printf(string, ...interface{}) }) *FlutterwaveGatewayImpl {
	return &FlutterwaveGatewayImpl{
		secretKey: secretKey,
		publicKey: publicKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// CreatePaymentIntent creates a Flutterwave payment link
func (fg *FlutterwaveGatewayImpl) CreatePaymentIntent(ctx context.Context, input *CreatePaymentIntentInput) (*PaymentIntentResult, error) {
	if input.Amount <= 0 {
		return nil, fmt.Errorf("invalid amount: %.2f", input.Amount)
	}

	// Prepare request to Flutterwave API
	// POST https://api.flutterwave.com/v3/payments

	payload := map[string]interface{}{
		"amount":       fmt.Sprintf("%.2f", input.Amount),
		"currency":     strings.ToUpper(input.Currency),
		"tx_ref":       fmt.Sprintf("order_%d_%d", input.OrderID, time.Now().UnixNano()),
		"redirect_url": "https://yourapp.com/payment/callback",
		"customer": map[string]interface{}{
			"id": input.UserID,
		},
		"customizations": map[string]interface{}{
			"title":       "Teamart Payment",
			"description": "Order payment",
		},
		"meta": map[string]interface{}{
			"order_id": input.OrderID,
		},
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.flutterwave.com/v3/payments", strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+fg.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := fg.httpClient.Do(req)
	if err != nil {
		fg.logger.Printf("flutterwave API error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("flutterwave error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var flutterwaveResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &flutterwaveResponse); err != nil {
		return nil, err
	}

	data, _ := flutterwaveResponse["data"].(map[string]interface{})
	paymentID, _ := data["id"].(string)
	link, _ := data["link"].(string)

	result := &PaymentIntentResult{
		Provider:         "flutterwave",
		ProviderIntentID: paymentID,
		AuthURL:          &link,
		Status:           "pending",
	}

	fg.logger.Printf("flutterwave payment link created: %s", paymentID)
	return result, nil
}

// ProcessPayment verifies and processes a Flutterwave payment
func (fg *FlutterwaveGatewayImpl) ProcessPayment(ctx context.Context, input *ProcessPaymentInput) (*PaymentResult, error) {
	// GET https://api.flutterwave.com/v3/transactions/{id}/verify

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.flutterwave.com/v3/transactions/%s/verify", input.ProviderToken), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+fg.secretKey)

	resp, err := fg.httpClient.Do(req)
	if err != nil {
		fg.logger.Printf("flutterwave verification error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var flutterwaveResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &flutterwaveResponse); err != nil {
		return nil, err
	}

	data, _ := flutterwaveResponse["data"].(map[string]interface{})
	status, _ := data["status"].(string)
	success := status == "successful"

	result := &PaymentResult{
		Success: success,
		Status:  status,
		Message: "Payment verified",
	}

	if !success {
		result.Message = "Payment verification failed"
	}

	fg.logger.Printf("flutterwave payment verified: %s", status)
	return result, nil
}

// RefundPayment refunds a Flutterwave payment
func (fg *FlutterwaveGatewayImpl) RefundPayment(ctx context.Context, input *RefundPaymentInput) (*RefundResult, error) {
	// POST https://api.flutterwave.com/v3/transactions/{id}/refund

	payload := map[string]interface{}{
		"amount": fmt.Sprintf("%.2f", input.Amount),
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("https://api.flutterwave.com/v3/transactions/%d/refund", input.PaymentIntentID), strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+fg.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := fg.httpClient.Do(req)
	if err != nil {
		fg.logger.Printf("flutterwave refund error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var flutterwaveResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &flutterwaveResponse); err != nil {
		return nil, err
	}

	data, _ := flutterwaveResponse["data"].(map[string]interface{})
	status, _ := data["status"].(string)

	result := &RefundResult{
		Success: status == "successful",
		Status:  status,
		Amount:  input.Amount,
		Message: "Refund processed",
	}

	fg.logger.Printf("flutterwave refund processed: %s", status)
	return result, nil
}

// VerifyWebhookSignature verifies Flutterwave webhook signature
func (fg *FlutterwaveGatewayImpl) VerifyWebhookSignature(signature, timestamp, payload string) (bool, error) {
	// Flutterwave uses SHA256 for webhook signing
	h := sha256.New()
	h.Write([]byte(payload + fg.secretKey))
	expected := hex.EncodeToString(h.Sum(nil))

	return signature == expected, nil
}

// GetPaymentStatus retrieves payment status from Flutterwave
func (fg *FlutterwaveGatewayImpl) GetPaymentStatus(ctx context.Context, providerPaymentID string) (*PaymentStatus, error) {
	// GET https://api.flutterwave.com/v3/transactions/{id}/verify

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.flutterwave.com/v3/transactions/%s/verify", providerPaymentID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+fg.secretKey)

	resp, err := fg.httpClient.Do(req)
	if err != nil {
		fg.logger.Printf("flutterwave status check error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var flutterwaveResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &flutterwaveResponse); err != nil {
		return nil, err
	}

	data, _ := flutterwaveResponse["data"].(map[string]interface{})
	status, _ := data["status"].(string)
	amount, _ := data["amount"].(float64)
	currency, _ := data["currency"].(string)

	paymentStatus := &PaymentStatus{
		ProviderPaymentID: providerPaymentID,
		Status:            status,
		Amount:            amount,
		Currency:          currency,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	return paymentStatus, nil
}
