package manapool

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	// DefaultBaseURL is the default base URL for the Manapool API.
	DefaultBaseURL = "https://manapool.com/api/v1/"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second

	// DefaultRateLimit is the default rate limit (requests per second).
	DefaultRateLimit = 10.0

	// DefaultRateBurst is the default rate limit burst.
	DefaultRateBurst = 1

	// DefaultMaxRetries is the default maximum number of retry attempts.
	DefaultMaxRetries = 3

	// DefaultInitialBackoff is the default initial backoff duration for retries.
	DefaultInitialBackoff = 1 * time.Second

	// Version is the library version.
	Version = "0.2.0"
)

// APIClient defines the interface for interacting with the Manapool API.
// This interface allows for easy mocking and testing.
type APIClient interface {
	// GetSellerAccount retrieves the seller account information.
	GetSellerAccount(ctx context.Context) (*Account, error)

	// GetSellerInventory retrieves the seller's inventory with pagination.
	GetSellerInventory(ctx context.Context, opts InventoryOptions) (*InventoryResponse, error)

	// GetInventoryByTCGPlayerID retrieves a specific inventory item by TCGPlayer SKU.
	GetInventoryByTCGPlayerID(ctx context.Context, tcgplayerID string) (*InventoryItem, error)
}

// Client is the Manapool API client.
// It implements the APIClient interface.
type Client struct {
	// httpClient is the HTTP client used for making requests
	httpClient *http.Client

	// baseURL is the base URL for the API
	baseURL string

	// authToken is the API authentication token
	authToken string

	// email is the account email address
	email string

	// rateLimiter limits the rate of API requests
	rateLimiter *rate.Limiter

	// maxRetries is the maximum number of retry attempts
	maxRetries int

	// initialBackoff is the initial backoff duration for retries
	initialBackoff time.Duration

	// userAgent is the User-Agent header value
	userAgent string

	// logger is used for debug and error logging
	logger Logger
}

// Logger is an interface for logging.
// Implement this interface to provide custom logging.
type Logger interface {
	// Debugf logs a debug message.
	Debugf(format string, args ...interface{})

	// Errorf logs an error message.
	Errorf(format string, args ...interface{})
}

// noopLogger is a no-op logger that discards all log messages.
type noopLogger struct{}

func (l *noopLogger) Debugf(format string, args ...interface{}) {}
func (l *noopLogger) Errorf(format string, args ...interface{}) {}

// NewClient creates a new Manapool API client.
// The authToken and email parameters are required for authentication.
// Additional options can be passed to configure the client.
//
// Example:
//
//	client := manapool.NewClient("your-token", "your-email@example.com",
//	    manapool.WithTimeout(60 * time.Second),
//	    manapool.WithRateLimit(5, 2),
//	)
func NewClient(authToken, email string, opts ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL:        DefaultBaseURL,
		authToken:      authToken,
		email:          email,
		rateLimiter:    rate.NewLimiter(DefaultRateLimit, DefaultRateBurst),
		maxRetries:     DefaultMaxRetries,
		initialBackoff: DefaultInitialBackoff,
		userAgent:      fmt.Sprintf("manapool-go/%s", Version),
		logger:         &noopLogger{},
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// doRequest executes an HTTP request with rate limiting, retries, and error handling.
func (c *Client) doRequest(ctx context.Context, method, endpoint string, params url.Values) (*http.Response, error) {
	// Wait for rate limiter
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, NewNetworkError("rate limiter error", err)
	}

	// Build URL
	reqURL := c.baseURL + strings.TrimPrefix(endpoint, "/")
	if len(params) > 0 {
		reqURL = reqURL + "?" + params.Encode()
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, reqURL, nil)
	if err != nil {
		return nil, NewNetworkError("failed to create request", err)
	}

	// Add headers
	req.Header.Set("X-ManaPool-Access-Token", c.authToken)
	req.Header.Set("X-ManaPool-Email", c.email)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	// Execute with retries
	var resp *http.Response
	backoff := c.initialBackoff

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		c.logger.Debugf("API request: %s %s (attempt %d/%d)", method, reqURL, attempt+1, c.maxRetries+1)

		resp, err = c.httpClient.Do(req)
		if err != nil {
			c.logger.Errorf("Request failed (attempt %d/%d): %v", attempt+1, c.maxRetries+1, err)

			// Don't retry on context errors
			if ctx.Err() != nil {
				return nil, NewNetworkError("request cancelled", ctx.Err())
			}

			// Retry on network errors
			if attempt < c.maxRetries {
				time.Sleep(backoff)
				backoff *= 2
				continue
			}

			return nil, NewNetworkError("request failed after retries", err)
		}

		// Success or non-retryable error
		if resp.StatusCode < 500 || attempt == c.maxRetries {
			break
		}

		// Server error - retry
		c.logger.Errorf("Server error %d (attempt %d/%d), retrying...", resp.StatusCode, attempt+1, c.maxRetries+1)
		_ = resp.Body.Close()
		time.Sleep(backoff)
		backoff *= 2
	}

	return resp, nil
}

// decodeResponse decodes a JSON response and handles HTTP errors.
func (c *Client) decodeResponse(resp *http.Response, v interface{}) error {
	defer func() {
		_ = resp.Body.Close()
	}()

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewNetworkError("failed to read response body", err)
	}

	c.logger.Debugf("API response: status=%d, body=%s", resp.StatusCode, string(body))

	// Check status code
	if resp.StatusCode != http.StatusOK {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
			Response:   resp,
		}

		// Try to extract a better error message from JSON
		var errorResp struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		if json.Unmarshal(body, &errorResp) == nil {
			if errorResp.Error != "" {
				apiErr.Message = errorResp.Error
			} else if errorResp.Message != "" {
				apiErr.Message = errorResp.Message
			}
		}

		return apiErr
	}

	// Decode JSON
	if v != nil {
		if err := json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
