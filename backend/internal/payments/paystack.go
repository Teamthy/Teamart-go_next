package payments

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// PaystackGatewayImpl implements PaymentGatewayProvider for Paystack
type PaystackGatewayImpl struct {
	secretKey  string
	publicKey  string
	httpClient *http.Client
	logger     interface{ Printf(string, ...interface{}) }
}

// NewPaystackGatewayImpl creates a new Paystack gateway implementation
func NewPaystackGatewayImpl(secretKey, publicKey string, logger interface{ Printf(string, ...interface{}) }) *PaystackGatewayImpl {
	return &PaystackGatewayImpl{
		secretKey: secretKey,
		publicKey: publicKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// CreatePaymentIntent creates a Paystack transaction
func (pg *PaystackGatewayImpl) CreatePaymentIntent(ctx context.Context, input *CreatePaymentIntentInput) (*PaymentIntentResult, error) {
	if input.Amount <= 0 {
		return nil, fmt.Errorf("invalid amount: %.2f", input.Amount)
	}

	// Convert amount to Kobo (1 Naira = 100 Kobo)
	amountKobo := int64(input.Amount * 100)

	// Prepare request to Paystack API
	// POST https://api.paystack.co/transaction/initialize

	payload := map[string]interface{}{
		"amount":    amountKobo,
		"currency":  strings.ToUpper(input.Currency),
		"reference": fmt.Sprintf("order_%d_%d", input.OrderID, time.Now().UnixNano()),
		"metadata": map[string]interface{}{
			"order_id": input.OrderID,
			"user_id":  input.UserID,
		},
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.paystack.co/transaction/initialize", strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+pg.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := pg.httpClient.Do(req)
	if err != nil {
		pg.logger.Printf("paystack API error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("paystack error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var paystackResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &paystackResponse); err != nil {
		return nil, err
	}

	data, _ := paystackResponse["data"].(map[string]interface{})
	reference, _ := data["reference"].(string)
	authURL, _ := data["authorization_url"].(string)

	result := &PaymentIntentResult{
		Provider:         "paystack",
		ProviderIntentID: reference,
		AuthURL:          &authURL,
		Status:           "pending",
	}

	pg.logger.Printf("paystack transaction initialized: %s", reference)
	return result, nil
}

// ProcessPayment verifies and processes a Paystack payment
func (pg *PaystackGatewayImpl) ProcessPayment(ctx context.Context, input *ProcessPaymentInput) (*PaymentResult, error) {
	// GET https://api.paystack.co/transaction/verify/{reference}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.paystack.co/transaction/verify/%s", input.ProviderToken), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+pg.secretKey)

	resp, err := pg.httpClient.Do(req)
	if err != nil {
		pg.logger.Printf("paystack verification error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var paystackResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &paystackResponse); err != nil {
		return nil, err
	}

	data, _ := paystackResponse["data"].(map[string]interface{})
	status, _ := data["status"].(string)
	success := status == "success"

	result := &PaymentResult{
		Success: success,
		Status:  status,
		Message: "Payment verified",
	}

	if !success {
		result.Message = "Payment verification failed"
	}

	pg.logger.Printf("paystack payment verified: %s", status)
	return result, nil
}

// RefundPayment refunds a Paystack payment
func (pg *PaystackGatewayImpl) RefundPayment(ctx context.Context, input *RefundPaymentInput) (*RefundResult, error) {
	// POST https://api.paystack.co/refund

	amountKobo := int64(input.Amount * 100)

	payload := map[string]interface{}{
		"transaction": input.PaymentIntentID,
		"amount":      amountKobo,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.paystack.co/refund", strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+pg.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := pg.httpClient.Do(req)
	if err != nil {
		pg.logger.Printf("paystack refund error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var paystackResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &paystackResponse); err != nil {
		return nil, err
	}

	data, _ := paystackResponse["data"].(map[string]interface{})
	refundID, _ := data["refund_id"].(string)
	status, _ := data["status"].(string)

	result := &RefundResult{
		Success:          status == "success",
		ProviderRefundID: refundID,
		Status:           status,
		Amount:           input.Amount,
		Message:          "Refund processed",
	}

	pg.logger.Printf("paystack refund processed: %s", refundID)
	return result, nil
}

// VerifyWebhookSignature verifies Paystack webhook signature
func (pg *PaystackGatewayImpl) VerifyWebhookSignature(signature, timestamp, payload string) (bool, error) {
	// Paystack uses SHA512 for webhook signing
	h := sha512.New()
	h.Write([]byte(payload + pg.secretKey))
	expected := hex.EncodeToString(h.Sum(nil))

	return signature == expected, nil
}

// GetPaymentStatus retrieves payment status from Paystack
func (pg *PaystackGatewayImpl) GetPaymentStatus(ctx context.Context, providerPaymentID string) (*PaymentStatus, error) {
	// GET https://api.paystack.co/transaction/verify/{reference}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.paystack.co/transaction/verify/%s", providerPaymentID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+pg.secretKey)

	resp, err := pg.httpClient.Do(req)
	if err != nil {
		pg.logger.Printf("paystack status check error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var paystackResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &paystackResponse); err != nil {
		return nil, err
	}

	data, _ := paystackResponse["data"].(map[string]interface{})
	status, _ := data["status"].(string)
	amount, _ := data["amount"].(float64)
	currency, _ := data["currency"].(string)

	paymentStatus := &PaymentStatus{
		ProviderPaymentID: providerPaymentID,
		Status:            status,
		Amount:            amount / 100, // Convert from Kobo to Naira
		Currency:          currency,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	return paymentStatus, nil
}
