// examples/webhook/main.go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/cruso003/irembopay"
)

func main() {
	// Create a sandbox client
	client, err := irembopay.NewSandboxClient("your-secret-key")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Set up webhook handler
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read request body: %v", err)
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Get the signature from the header
		signature := r.Header.Get("irembopay-signature")
		if signature == "" {
			log.Printf("Missing signature header")
			http.Error(w, "Missing signature header", http.StatusBadRequest)
			return
		}

		// Validate the timestamp (prevent replay attacks)
		valid, err := client.Payment.ValidateWebhookTimestamp(signature, 5*time.Minute)
		if err != nil {
			log.Printf("Failed to validate timestamp: %v", err)
			http.Error(w, "Failed to validate timestamp", http.StatusBadRequest)
			return
		}
		if !valid {
			log.Printf("Timestamp is too old, possible replay attack")
			http.Error(w, "Timestamp is too old", http.StatusBadRequest)
			return
		}

		// Verify the signature
		payload := string(body)
		valid, err = client.Payment.VerifyWebhookSignature(signature, payload)
		if err != nil {
			log.Printf("Failed to verify signature: %v", err)
			http.Error(w, "Failed to verify signature", http.StatusBadRequest)
			return
		}
		if !valid {
			log.Printf("Invalid signature")
			http.Error(w, "Invalid signature", http.StatusBadRequest)
			return
		}

		// Parse the notification
		notification, err := client.Payment.ParseNotification(payload)
		if err != nil {
			log.Printf("Failed to parse notification: %v", err)
			http.Error(w, "Failed to parse notification", http.StatusBadRequest)
			return
		}

		// Process the notification
		log.Printf("Received payment notification:")
		log.Printf("  Invoice Number: %s", notification.InvoiceNumber)
		log.Printf("  Transaction ID: %s", notification.TransactionID)
		log.Printf("  Status: %s", notification.PaymentStatus)
		log.Printf("  Amount: %.2f %s", notification.Amount, notification.Currency)
		log.Printf("  Payment Method: %s", notification.PaymentMethod)
		log.Printf("  Paid At: %s", notification.PaidAt)

		// Respond with success
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	// Start the server
	port := ":8080"
	log.Printf("Starting webhook server on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
