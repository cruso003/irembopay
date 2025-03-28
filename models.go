package irembopay

import (
	"encoding/json"
	"time"
)

// PaymentItem represents an item in the invoice
type PaymentItem struct {
	Code       string  `json:"code"`       // Identifier of the product
	Quantity   int     `json:"quantity"`   // Must be > 0
	UnitAmount float64 `json:"unitAmount"` // Price per unit
}

// Customer represents the customer information
type Customer struct {
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Name        string `json:"name,omitempty"`
}

// InvoiceRequest represents the request to create an invoice
type InvoiceRequest struct {
	TransactionID            string        `json:"transactionId"`            // Unique transaction identifier
	PaymentAccountIdentifier string        `json:"paymentAccountIdentifier"` // Identifier of the payment account
	PaymentItems             []PaymentItem `json:"paymentItems"`             // List of items to be paid
	ExpiryAt                 string        `json:"expiryAt,omitempty"`       // Time when the invoice will expire
	Description              string        `json:"description,omitempty"`    // Description of the invoice
	Customer                 *Customer     `json:"customer,omitempty"`       // Customer information
	Language                 string        `json:"language,omitempty"`       // Language (FR, EN, RW)
}

// BatchInvoiceRequest represents the request to create a batch invoice
type BatchInvoiceRequest struct {
	TransactionID  string   `json:"transactionId"`         // Unique transaction identifier
	InvoiceNumbers []string `json:"invoiceNumbers"`        // List of invoice numbers to include in the batch
	Description    string   `json:"description,omitempty"` // Description of the batch invoice
}

// UpdateInvoiceRequest represents the request to update an invoice
type UpdateInvoiceRequest struct {
	ExpiryAt     string        `json:"expiryAt,omitempty"`     // New expiration date
	PaymentItems []PaymentItem `json:"paymentItems,omitempty"` // Updated payment items
}

// Invoice represents an invoice response from IremboPay
type Invoice struct {
	Amount                   float64       `json:"amount"`                     // Amount of the invoice
	InvoiceNumber            string        `json:"invoiceNumber"`              // Identifier of the invoice
	TransactionID            string        `json:"transactionId"`              // Transaction identifier
	CreatedAt                string        `json:"createdAt"`                  // Creation date
	UpdatedAt                string        `json:"updatedAt,omitempty"`        // Last update date
	ExpiryAt                 string        `json:"expiryAt,omitempty"`         // Expiration date
	PaidAt                   string        `json:"paidAt,omitempty"`           // Payment date
	PaymentAccountIdentifier string        `json:"paymentAccountIdentifier"`   // Payment account identifier
	PaymentItems             []PaymentItem `json:"paymentItems"`               // List of payment items
	Description              string        `json:"description,omitempty"`      // Description of the invoice
	Type                     string        `json:"type"`                       // SINGLE or BATCH
	PaymentStatus            string        `json:"paymentStatus"`              // NEW or PAID
	PaymentReference         string        `json:"paymentReference,omitempty"` // Reference provided after payment
	PaymentMethod            string        `json:"paymentMethod,omitempty"`    // MTN_MOMO, AIRTEL_MONEY, etc.
	Currency                 string        `json:"currency"`                   // RWF, EUR, GBP, USD
	Customer                 *Customer     `json:"customer,omitempty"`         // Customer information
	Language                 string        `json:"language,omitempty"`         // Language
	BatchNumber              string        `json:"batchNumber,omitempty"`      // Batch invoice number
	ChildInvoices            []string      `json:"childInvoices,omitempty"`    // Invoices in the batch
	PaymentLinkUrl           string        `json:"paymentLinkUrl"`             // Checkout URL
}

// MomoPaymentRequest represents a request to initiate a mobile money payment
type MomoPaymentRequest struct {
	AccountIdentifier    string `json:"accountIdentifier"`              // Phone number
	PaymentProvider      string `json:"paymentProvider"`                // MTN or AIRTEL
	InvoiceNumber        string `json:"invoiceNumber"`                  // Invoice number
	TransactionReference string `json:"transactionReference,omitempty"` // Optional reference
}

// MomoPaymentResponse represents the response to a mobile money payment request
type MomoPaymentResponse struct {
	AccountIdentifier string  `json:"accountIdentifier"` // Phone number
	PaymentProvider   string  `json:"paymentProvider"`   // MTN or AIRTEL
	InvoiceNumber     string  `json:"invoiceNumber"`     // Invoice number
	Amount            float64 `json:"amount"`            // Amount
	ReferenceID       string  `json:"referenceId"`       // IremboPay reference
}

// Response is the standard API response from IremboPay
type Response struct {
	Message string          `json:"message"` // Response message
	Success bool            `json:"success"` // Whether the request was successful
	Data    json.RawMessage `json:"data"`    // Response data
}

// PaymentNotification represents a payment notification from IremboPay
type PaymentNotification struct {
	InvoiceNumber     string  `json:"invoiceNumber"`     // Invoice number
	TransactionID     string  `json:"transactionId"`     // Transaction ID
	PaymentStatus     string  `json:"paymentStatus"`     // Payment status
	PaymentReference  string  `json:"paymentReference"`  // Payment reference
	Amount            float64 `json:"amount"`            // Amount
	Currency          string  `json:"currency"`          // Currency
	PaymentMethod     string  `json:"paymentMethod"`     // Payment method
	PaidAt            string  `json:"paidAt"`            // Payment date
	PaymentAccountID  string  `json:"paymentAccountId"`  // Payment account ID
	PaymentMerchantID string  `json:"paymentMerchantId"` // Merchant ID
}

// FormatTime formats a time.Time for IremboPay API (RFC3339 format)
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseTime parses a time string from IremboPay API
func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
