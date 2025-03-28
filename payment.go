package irembopay

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// PaymentService handles payment operations
type PaymentService struct {
	client *Client
	config *Config
}

// NewPaymentService creates a new payment service
func NewPaymentService(client *Client, config *Config) *PaymentService {
	return &PaymentService{
		client: client,
		config: config,
	}
}

// InitiateMomoPayment initiates a mobile money payment
func (s *PaymentService) InitiateMomoPayment(ctx context.Context, req *MomoPaymentRequest) (*MomoPaymentResponse, error) {
	var response MomoPaymentResponse
	apiReq := Request{
		Method: http.MethodPost,
		Path:   "/payments/transactions/initiate",
		Body:   req,
	}

	err := s.client.DoRequest(ctx, apiReq, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate mobile money payment: %w", err)
	}

	return &response, nil
}

// VerifyWebhookSignature verifies the signature of a webhook notification
func (s *PaymentService) VerifyWebhookSignature(signature, payload string) (bool, error) {
	// Parse the signature header
	// Format: t=<timestamp>, s=<signature>
	parts := strings.Split(signature, ",")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid signature format")
	}

	var timestamp, sig string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "t=") {
			timestamp = strings.TrimPrefix(part, "t=")
		} else if strings.HasPrefix(part, "s=") {
			sig = strings.TrimPrefix(part, "s=")
		}
	}

	if timestamp == "" || sig == "" {
		return false, fmt.Errorf("missing timestamp or signature")
	}

	// Verify the signature
	// signature = HMAC_SHA256(secretKey, timestamp + "#" + payload)
	mac := hmac.New(sha256.New, []byte(s.config.SecretKey))
	mac.Write([]byte(timestamp + "#" + payload))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	return expectedSig == sig, nil
}

// ParseNotification parses a payment notification from a webhook payload
func (s *PaymentService) ParseNotification(payload string) (*PaymentNotification, error) {
	var notification PaymentNotification
	err := json.Unmarshal([]byte(payload), &notification)
	if err != nil {
		return nil, fmt.Errorf("failed to parse notification: %w", err)
	}

	return &notification, nil
}

// HandleWebhook is a utility function to handle webhook notifications
func (s *PaymentService) HandleWebhook(signature, payload string) (*PaymentNotification, error) {
	// Verify the signature
	valid, err := s.VerifyWebhookSignature(signature, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to verify signature: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid signature")
	}

	// Parse the notification
	notification, err := s.ParseNotification(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse notification: %w", err)
	}

	return notification, nil
}

// ValidateWebhookTimestamp validates that the webhook timestamp is not too old
// to prevent replay attacks
func (s *PaymentService) ValidateWebhookTimestamp(signature string, maxAge time.Duration) (bool, error) {
	// Parse the signature header
	parts := strings.Split(signature, ",")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid signature format")
	}

	var timestamp string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "t=") {
			timestamp = strings.TrimPrefix(part, "t=")
			break
		}
	}

	if timestamp == "" {
		return false, fmt.Errorf("missing timestamp")
	}

	// Convert timestamp to time.Time
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid timestamp format: %w", err)
	}

	webhookTime := time.Unix(ts/1000, 0) // Convert milliseconds to seconds
	now := time.Now()

	// Check if the timestamp is too old
	if now.Sub(webhookTime) > maxAge {
		return false, nil
	}

	return true, nil
}
