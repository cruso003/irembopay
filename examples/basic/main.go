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

	// Create an invoice
	ctx := context.Background()
	invoice, err := client.Invoice.Create(ctx, &irembopay.InvoiceRequest{
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
	})
	if err != nil {
		log.Fatalf("Failed to create invoice: %v", err)
	}

	fmt.Printf("Created invoice: %s\n", invoice.InvoiceNumber)
	fmt.Printf("Payment link: %s\n", invoice.PaymentLinkUrl)

	// Get the invoice details
	fetchedInvoice, err := client.Invoice.Get(ctx, invoice.InvoiceNumber)
	if err != nil {
		log.Fatalf("Failed to get invoice: %v", err)
	}

	fmt.Printf("Fetched invoice: %s\n", fetchedInvoice.InvoiceNumber)
	fmt.Printf("Status: %s\n", fetchedInvoice.PaymentStatus)

	// Initiate a mobile money payment
	momoPayment, err := client.Payment.InitiateMomoPayment(ctx, &irembopay.MomoPaymentRequest{
		AccountIdentifier: "0780000001",
		PaymentProvider:   "MTN",
		InvoiceNumber:     invoice.InvoiceNumber,
	})
	if err != nil {
		log.Fatalf("Failed to initiate payment: %v", err)
	}

	fmt.Printf("Initiated payment: %s\n", momoPayment.ReferenceID)
}
