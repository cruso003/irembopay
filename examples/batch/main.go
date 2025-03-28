// examples/batch/main.go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cruso003/irembopay"
)

func main() {
	// Create a sandbox client
	client, err := irembopay.NewSandboxClient("your-secret-key")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create two individual invoices first
	ctx := context.Background()

	invoice1, err := client.Invoice.Create(ctx, &irembopay.InvoiceRequest{
		TransactionID:            "TST-12346",
		PaymentAccountIdentifier: "TST-RWF",
		PaymentItems: []irembopay.PaymentItem{
			{
				Code:       "PI-3e5fe23f2d",
				Quantity:   1,
				UnitAmount: 1000,
			},
		},
		Description: "Invoice 1",
	})
	if err != nil {
		log.Fatalf("Failed to create invoice 1: %v", err)
	}

	invoice2, err := client.Invoice.Create(ctx, &irembopay.InvoiceRequest{
		TransactionID:            "TST-12347",
		PaymentAccountIdentifier: "TST-RWF",
		PaymentItems: []irembopay.PaymentItem{
			{
				Code:       "PI-3e5fe23f2d",
				Quantity:   2,
				UnitAmount: 1500,
			},
		},
		Description: "Invoice 2",
	})
	if err != nil {
		log.Fatalf("Failed to create invoice 2: %v", err)
	}

	fmt.Printf("Created invoice 1: %s\n", invoice1.InvoiceNumber)
	fmt.Printf("Created invoice 2: %s\n", invoice2.InvoiceNumber)

	// Create a batch invoice
	batchInvoice, err := client.Batch.Create(ctx, &irembopay.BatchInvoiceRequest{
		TransactionID:  "TST-BATCH-123",
		InvoiceNumbers: []string{invoice1.InvoiceNumber, invoice2.InvoiceNumber},
		Description:    "Batch invoice",
	})
	if err != nil {
		log.Fatalf("Failed to create batch invoice: %v", err)
	}

	fmt.Printf("Created batch invoice: %s\n", batchInvoice.InvoiceNumber)
	fmt.Printf("Batch total amount: %.2f %s\n", batchInvoice.Amount, batchInvoice.Currency)
	fmt.Printf("Payment link: %s\n", batchInvoice.PaymentLinkUrl)
}
