// examples/basic/main.go
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

	// Generate a unique idempotency key
	orderID := "ORDER-12345"
	idempotencyKey := irembopay.GenerateIdempotencyKey("invoice", orderID)
	fmt.Printf("Using idempotency key: %s\n", idempotencyKey)

	// Create an invoice with idempotency
	ctx := context.Background()
	invoiceReq := &irembopay.InvoiceRequest{
		TransactionID:            "TST-12345",
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
	}

	invoice, err := client.Invoice.CreateWithIdempotency(ctx, invoiceReq, idempotencyKey)
	if err != nil {
		log.Fatalf("Failed to create invoice: %v", err)
	}

	fmt.Printf("Created invoice: %s\n", invoice.InvoiceNumber)
	fmt.Printf("Payment link: %s\n", invoice.PaymentLinkUrl)

	// Try to create the same invoice again (should be idempotent)
	fmt.Println("\nTrying to create the same invoice again...")

	duplicateInvoice, err := client.Invoice.CreateWithIdempotency(ctx, invoiceReq, idempotencyKey)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Received invoice: %s\n", duplicateInvoice.InvoiceNumber)

		// Check if it's the same invoice
		if duplicateInvoice.InvoiceNumber == invoice.InvoiceNumber {
			fmt.Println("Success! Received the same invoice (idempotent request worked)")
		} else {
			fmt.Println("Warning: Received a different invoice")
		}
	}
}
