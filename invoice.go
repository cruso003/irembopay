package irembopay

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// InvoiceService handles operations on IremboPay invoices
type InvoiceService struct {
	client *Client
	config *Config
}

// NewInvoiceService creates a new invoice service
func NewInvoiceService(client *Client, config *Config) *InvoiceService {
	return &InvoiceService{
		client: client,
		config: config,
	}
}

// Create creates a new invoice
func (s *InvoiceService) Create(ctx context.Context, req *InvoiceRequest) (*Invoice, error) {
	var invoice Invoice
	apiReq := Request{
		Method: http.MethodPost,
		Path:   "/payments/invoices",
		Body:   req,
	}

	err := s.client.DoRequest(ctx, apiReq, &invoice)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	return &invoice, nil
}

// CreateWithExpiry creates a new invoice with an expiry time
func (s *InvoiceService) CreateWithExpiry(ctx context.Context, req *InvoiceRequest, expiryDuration time.Duration) (*Invoice, error) {
	// Calculate expiry time
	expiryTime := time.Now().Add(expiryDuration)
	req.ExpiryAt = FormatTime(expiryTime)

	return s.Create(ctx, req)
}

// Get retrieves an invoice by its number or transaction ID
func (s *InvoiceService) Get(ctx context.Context, invoiceReference string) (*Invoice, error) {
	var invoice Invoice
	apiReq := Request{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/payments/invoices/%s", invoiceReference),
	}

	err := s.client.DoRequest(ctx, apiReq, &invoice)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	return &invoice, nil
}

// Update updates an existing invoice
func (s *InvoiceService) Update(ctx context.Context, invoiceNumber string, req *UpdateInvoiceRequest) (*Invoice, error) {
	var invoice Invoice
	apiReq := Request{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/payments/invoices/%s", invoiceNumber),
		Body:   req,
	}

	err := s.client.DoRequest(ctx, apiReq, &invoice)
	if err != nil {
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}

	return &invoice, nil
}

// UpdateExpiryTime updates the expiry time of an invoice
func (s *InvoiceService) UpdateExpiryTime(ctx context.Context, invoiceNumber string, expiryTime time.Time) (*Invoice, error) {
	req := &UpdateInvoiceRequest{
		ExpiryAt: FormatTime(expiryTime),
	}

	return s.Update(ctx, invoiceNumber, req)
}
