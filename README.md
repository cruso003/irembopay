# IremboPay Go Client

A simple, clean, and easy-to-use Go client for integrating with the IremboPay API. This package supports both sandbox and production environments.

## Features

- Support for both Sandbox and Production environments
- Invoice creation, retrieval, and updating
- Batch invoice creation
- Mobile money payment initiation
- Webhook signature verification and notification parsing
- Comprehensive error handling

## Installation

```bash
go get github.com/cruso003/irembopay
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/cruso003/irembopay"
)

func main() {
    // Create a sandbox client
    client, err := irembopay.NewSandboxClient("your-secret-key")
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Create an invoice
    invoice, err := client.Invoice.Create(context.Background(), &irembopay.InvoiceRequest{
        TransactionID:           "TST-12345",
        PaymentAccountIdentifier: "TST-RWF",
        PaymentItems: []irembopay.PaymentItem{
            {
                Code:       "PI-3e5fe23f2d",
                Quantity:   1,
                UnitAmount: 2000,
            },
        },
        ExpiryAt:    irembopay.FormatTime(time.Now().Add(24 * time.Hour)),
        Description: "Test invoice",
        Customer: &irembopay.Customer{
            Email:       "test@example.com",
            PhoneNumber: "0780000001",
            Name:        "Test User",
        },
        Language: "EN",
    })
    if err != nil {
        log.Fatalf("Failed to create invoice: %v", err)
    }

    fmt.Printf("Created invoice: %s\n", invoice.InvoiceNumber)
    fmt.Printf("Payment link: %s\n", invoice.PaymentLinkUrl)
}
```

## Configuration

You can configure the client for sandbox or production:

```go
// Sandbox client
client, err := irembopay.NewSandboxClient("your-sandbox-secret-key")

// Production client
client, err := irembopay.NewProductionClient("your-production-secret-key")

// With custom options
client, err := irembopay.NewSandboxClient(
    "your-secret-key",
    irembopay.WithAPIVersion("2"),
    irembopay.WithHost("custom-sandbox.irembopay.com"),
)
```

## Usage Examples

### Creating an Invoice

```go
invoice, err := client.Invoice.Create(ctx, &irembopay.InvoiceRequest{
    TransactionID:           "TST-12345",
    PaymentAccountIdentifier: "TST-RWF",
    PaymentItems: []irembopay.PaymentItem{
        {
            Code:       "PI-3e5fe23f2d",
            Quantity:   1,
            UnitAmount: 2000,
        },
    },
    ExpiryAt:    irembopay.FormatTime(time.Now().Add(24 * time.Hour)),
    Description: "Test invoice",
    Customer: &irembopay.Customer{
        Email:       "test@example.com",
        PhoneNumber: "0780000001",
        Name:        "Test User",
    },
    Language: "EN",
})
```

### Getting an Invoice

```go
invoice, err := client.Invoice.Get(ctx, "880419623157")
```

### Updating an Invoice

```go
updated, err := client.Invoice.Update(ctx, "880419623157", &irembopay.UpdateInvoiceRequest{
    ExpiryAt: irembopay.FormatTime(time.Now().Add(48 * time.Hour)),
    PaymentItems: []irembopay.PaymentItem{
        {
            Code:       "PI-3e5fe23f2d",
            Quantity:   2,
            UnitAmount: 1500,
        },
    },
})
```

### Creating a Batch Invoice

```go
batchInvoice, err := client.Batch.Create(ctx, &irembopay.BatchInvoiceRequest{
    TransactionID:  "TST-BATCH-123",
    InvoiceNumbers: []string{"880419623157", "880419623158"},
    Description:    "Batch invoice",
})
```

### Initiating a Mobile Money Payment

```go
payment, err := client.Payment.InitiateMomoPayment(ctx, &irembopay.MomoPaymentRequest{
    AccountIdentifier: "0780000001",
    PaymentProvider:   "MTN",
    InvoiceNumber:     "880419623157",
})
```

### Handling Webhooks

```go
http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
    // Read the request body
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // Get the signature from the header
    signature := r.Header.Get("irembopay-signature")
    if signature == "" {
        http.Error(w, "Missing signature header", http.StatusBadRequest)
        return
    }

    // Verify and process the notification
    notification, err := client.Payment.HandleWebhook(signature, string(body))
    if err != nil {
        http.Error(w, "Failed to process webhook", http.StatusBadRequest)
        return
    }

    // Process the notification
    fmt.Printf("Received payment for invoice: %s\n", notification.InvoiceNumber)
    fmt.Printf("Amount: %.2f %s\n", notification.Amount, notification.Currency)
})
```

## Error Handling

The package provides specific error types for better error handling:

```go
invoice, err := client.Invoice.Get(ctx, "nonexistent-invoice")
if err != nil {
    if irembopay.IsNotFoundError(err) {
        fmt.Println("Invoice not found!")
    } else if irembopay.IsBadRequestError(err) {
        fmt.Println("Bad request!")
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
