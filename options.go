package manapool

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client.
// Use this to configure timeouts, transport, TLS settings, etc.
//
// Example:
//
//	customClient := &http.Client{
//	    Timeout: 30 * time.Second,
//	    Transport: &http.Transport{
//	        MaxIdleConns: 10,
//	    },
//	}
//	client := manapool.NewClient(token, email,
//	    manapool.WithHTTPClient(customClient),
//	)
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets a custom base URL for the API.
// This is useful for testing against a mock server or staging environment.
//
// The default base URL is https://manapool.com/api/v1/
//
// Example:
//
//	client := manapool.NewClient(token, email,
//	    manapool.WithBaseURL("https://staging.manapool.com/api/v1/"),
//	)
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithRateLimit configures rate limiting for API requests.
// The rate is specified as requests per second, and burst is the maximum
// number of requests that can be made in a single burst.
//
// Default: 10 requests per second with a burst of 1.
//
// Example:
//
//	client := manapool.NewClient(token, email,
//	    manapool.WithRateLimit(5, 2), // 5 req/sec, burst of 2
//	)
func WithRateLimit(requestsPerSecond float64, burst int) ClientOption {
	return func(c *Client) {
		c.rateLimiter = rate.NewLimiter(rate.Limit(requestsPerSecond), burst)
	}
}

// WithTimeout sets the HTTP client timeout.
// This is a convenience method that wraps WithHTTPClient.
//
// Default: 30 seconds.
//
// Example:
//
//	client := manapool.NewClient(token, email,
//	    manapool.WithTimeout(60 * time.Second),
//	)
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}

// WithRetry configures automatic retry behavior for failed requests.
// maxRetries specifies the maximum number of retry attempts.
// initialBackoff specifies the initial backoff duration (doubled on each retry).
//
// Default: 3 retries with 1 second initial backoff.
//
// Example:
//
//	client := manapool.NewClient(token, email,
//	    manapool.WithRetry(5, 2*time.Second), // 5 retries, 2s initial backoff
//	)
func WithRetry(maxRetries int, initialBackoff time.Duration) ClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
		c.initialBackoff = initialBackoff
	}
}

// WithUserAgent sets a custom User-Agent header for API requests.
//
// Default: "manapool-go/<version>"
//
// Example:
//
//	client := manapool.NewClient(token, email,
//	    manapool.WithUserAgent("my-app/1.0"),
//	)
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithLogger sets a custom logger for the client.
// The logger must implement the Logger interface.
//
// Example:
//
//	type myLogger struct{}
//	func (l *myLogger) Debugf(format string, args ...interface{}) {
//	    log.Printf("[DEBUG] "+format, args...)
//	}
//	func (l *myLogger) Errorf(format string, args ...interface{}) {
//	    log.Printf("[ERROR] "+format, args...)
//	}
//
//	client := manapool.NewClient(token, email,
//	    manapool.WithLogger(&myLogger{}),
//	)
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}
