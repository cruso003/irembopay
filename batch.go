package irembopay

import (
	"context"
	"fmt"
	"net/http"
)

// BatchService handles operations on IremboPay batch invoices
type BatchService struct {
	client *Client
	config *Config
}

// NewBatchService creates a new batch service
func NewBatchService(client *Client, config *Config) *BatchService {
	return &BatchService{
		client: client,
		config: config,
	}
}

// Create creates a new batch invoice
func (s *BatchService) Create(ctx context.Context, req *BatchInvoiceRequest) (*Invoice, error) {
	var invoice Invoice
	apiReq := Request{
		Method: http.MethodPost,
		Path:   "/payments/invoices/batch",
		Body:   req,
	}

	err := s.client.DoRequest(ctx, apiReq, &invoice)
	if err != nil {
		return nil, fmt.Errorf("failed to create batch invoice: %w", err)
	}

	return &invoice, nil
}

// CreateWithIdempotency creates a new batch invoice with an idempotency key
func (s *BatchService) CreateWithIdempotency(ctx context.Context, req *BatchInvoiceRequest, idempotencyKey string) (*Invoice, error) {
	req.IdempotencyKey = idempotencyKey
	return s.Create(ctx, req)
}
