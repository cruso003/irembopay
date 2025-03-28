package irembopay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client handles HTTP communication with the IremboPay API
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient creates a new IremboPay API client
func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Request represents an HTTP request to the API
type Request struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
	Params  map[string]string
}

// DoRequest performs an HTTP request and decodes the response
func (c *Client) DoRequest(ctx context.Context, req Request, result interface{}) error {
	var bodyReader io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	url := fmt.Sprintf("https://%s%s", c.config.Host, req.Path)
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Set default headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("irembopay-secretKey", c.config.SecretKey)
	httpReq.Header.Set("X-API-Version", c.config.APIVersion)

	// Add request-specific headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Add query parameters
	if len(req.Params) > 0 {
		q := httpReq.URL.Query()
		for key, value := range req.Params {
			q.Add(key, value)
		}
		httpReq.URL.RawQuery = q.Encode()
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	// Check for error status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorResp struct {
			Message string `json:"message"`
			Success bool   `json:"success"`
			Error   string `json:"error"`
		}

		if err := json.Unmarshal(body, &errorResp); err != nil {
			return fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
		}

		errorMessage := errorResp.Message
		if errorResp.Error != "" {
			errorMessage = errorResp.Error
		}

		return NewIremboPayError(resp.StatusCode, errorMessage, string(body))
	}

	// Parse the response
	if result != nil {
		var apiResp Response
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return fmt.Errorf("error parsing API response: %w", err)
		}

		if !apiResp.Success {
			return fmt.Errorf("API request unsuccessful: %s", apiResp.Message)
		}

		// Parse the data field into the provided result
		if err := json.Unmarshal(apiResp.Data, result); err != nil {
			return fmt.Errorf("error parsing response data: %w", err)
		}
	}

	return nil
}
