package manapool

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-token", "test@example.com")

	if client.authToken != "test-token" {
		t.Errorf("authToken = %q, want %q", client.authToken, "test-token")
	}
	if client.email != "test@example.com" {
		t.Errorf("email = %q, want %q", client.email, "test@example.com")
	}
	if client.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %q, want %q", client.baseURL, DefaultBaseURL)
	}
	if client.httpClient.Timeout != DefaultTimeout {
		t.Errorf("httpClient.Timeout = %v, want %v", client.httpClient.Timeout, DefaultTimeout)
	}
	if client.maxRetries != DefaultMaxRetries {
		t.Errorf("maxRetries = %d, want %d", client.maxRetries, DefaultMaxRetries)
	}
	if client.initialBackoff != DefaultInitialBackoff {
		t.Errorf("initialBackoff = %v, want %v", client.initialBackoff, DefaultInitialBackoff)
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	customClient := &http.Client{Timeout: 60 * time.Second}
	customLogger := &testLogger{}

	client := NewClient("token", "email",
		WithHTTPClient(customClient),
		WithBaseURL("https://custom.api.com/"),
		WithRateLimit(5, 2),
		WithRetry(5, 2*time.Second),
		WithUserAgent("custom-agent"),
		WithLogger(customLogger),
	)

	if client.httpClient != customClient {
		t.Error("WithHTTPClient did not set custom client")
	}
	if client.baseURL != "https://custom.api.com/" {
		t.Errorf("WithBaseURL did not set custom URL, got %q", client.baseURL)
	}
	if client.maxRetries != 5 {
		t.Errorf("WithRetry did not set maxRetries, got %d", client.maxRetries)
	}
	if client.initialBackoff != 2*time.Second {
		t.Errorf("WithRetry did not set initialBackoff, got %v", client.initialBackoff)
	}
	if client.userAgent != "custom-agent" {
		t.Errorf("WithUserAgent did not set custom agent, got %q", client.userAgent)
	}
	if client.logger != customLogger {
		t.Error("WithLogger did not set custom logger")
	}
}

func TestClient_doRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if got := r.Header.Get("X-ManaPool-Access-Token"); got != "test-token" {
			t.Errorf("X-ManaPool-Access-Token = %q, want %q", got, "test-token")
		}
		if got := r.Header.Get("X-ManaPool-Email"); got != "test@example.com" {
			t.Errorf("X-ManaPool-Email = %q, want %q", got, "test@example.com")
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Errorf("Accept = %q, want %q", got, "application/json")
		}
		if !strings.Contains(r.Header.Get("User-Agent"), "manapool-go") {
			t.Errorf("User-Agent = %q, want to contain 'manapool-go'", r.Header.Get("User-Agent"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx := context.Background()
	resp, err := client.doRequest(ctx, "GET", "/test", nil)
	if err != nil {
		t.Fatalf("doRequest() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestClient_doRequest_WithQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query params
		if got := r.URL.Query().Get("limit"); got != "100" {
			t.Errorf("query param 'limit' = %q, want %q", got, "100")
		}
		if got := r.URL.Query().Get("offset"); got != "0" {
			t.Errorf("query param 'offset' = %q, want %q", got, "0")
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := NewClient("token", "email",
		WithBaseURL(server.URL+"/"),
	)

	ctx := context.Background()
	params := map[string][]string{
		"limit":  {"100"},
		"offset": {"0"},
	}
	resp, err := client.doRequest(ctx, "GET", "/test", params)
	if err != nil {
		t.Fatalf("doRequest() error = %v", err)
	}
	defer resp.Body.Close()
}

func TestClient_doRequest_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("token", "email",
		WithBaseURL(server.URL+"/"),
		WithRateLimit(1000, 1), // High rate limit to avoid waiting
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.doRequest(ctx, "GET", "/test", nil)
	if err == nil {
		t.Error("doRequest() expected error for cancelled context")
	}

	var netErr *NetworkError
	if !errors.As(err, &netErr) {
		t.Errorf("expected NetworkError, got %T", err)
	}
}

func TestClient_doRequest_ServerError_Retry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status": "ok"}`))
		}
	}))
	defer server.Close()

	client := NewClient("token", "email",
		WithBaseURL(server.URL+"/"),
		WithRetry(3, 10*time.Millisecond), // Short backoff for testing
	)

	ctx := context.Background()
	resp, err := client.doRequest(ctx, "GET", "/test", nil)
	if err != nil {
		t.Fatalf("doRequest() error = %v", err)
	}
	defer resp.Body.Close()

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestClient_decodeResponse_Success(t *testing.T) {
	responseBody := `{"username": "testuser", "email": "test@example.com"}`
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(responseBody)),
	}

	client := NewClient("token", "email")

	var account Account
	err := client.decodeResponse(resp, &account)
	if err != nil {
		t.Fatalf("decodeResponse() error = %v", err)
	}

	if account.Username != "testuser" {
		t.Errorf("Username = %q, want %q", account.Username, "testuser")
	}
	if account.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", account.Email, "test@example.com")
	}
}

func TestClient_decodeResponse_HTTPError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    bool
		checkErr   func(*testing.T, error)
	}{
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
			body:       `{"error": "not found"}`,
			wantErr:    true,
			checkErr: func(t *testing.T, err error) {
				var apiErr *APIError
				if !errors.As(err, &apiErr) {
					t.Errorf("expected APIError, got %T", err)
					return
				}
				if !apiErr.IsNotFound() {
					t.Error("expected IsNotFound() to be true")
				}
				if apiErr.Message != "not found" {
					t.Errorf("Message = %q, want %q", apiErr.Message, "not found")
				}
			},
		},
		{
			name:       "401 Unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       `{"message": "unauthorized"}`,
			wantErr:    true,
			checkErr: func(t *testing.T, err error) {
				var apiErr *APIError
				if !errors.As(err, &apiErr) {
					t.Errorf("expected APIError, got %T", err)
					return
				}
				if !apiErr.IsUnauthorized() {
					t.Error("expected IsUnauthorized() to be true")
				}
			},
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			body:       `Internal Server Error`,
			wantErr:    true,
			checkErr: func(t *testing.T, err error) {
				var apiErr *APIError
				if !errors.As(err, &apiErr) {
					t.Errorf("expected APIError, got %T", err)
					return
				}
				if !apiErr.IsServerError() {
					t.Error("expected IsServerError() to be true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(strings.NewReader(tt.body)),
			}

			client := NewClient("token", "email")

			var result interface{}
			err := client.decodeResponse(resp, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("decodeResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.checkErr != nil {
				tt.checkErr(t, err)
			}
		})
	}
}

func TestClient_decodeResponse_InvalidJSON(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{invalid json}`)),
	}

	client := NewClient("token", "email")

	var result interface{}
	err := client.decodeResponse(resp, &result)
	if err == nil {
		t.Error("decodeResponse() expected error for invalid JSON")
	}
}

func TestClient_decodeResponse_NilTarget(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"status": "ok"}`)),
	}

	client := NewClient("token", "email")

	err := client.decodeResponse(resp, nil)
	if err != nil {
		t.Errorf("decodeResponse() error = %v, expected nil for nil target", err)
	}
}

func TestNoopLogger(t *testing.T) {
	logger := &noopLogger{}
	// These should not panic
	logger.Debugf("test %s", "debug")
	logger.Errorf("test %s", "error")
}

func TestWithTimeout(t *testing.T) {
	client := NewClient("token", "email",
		WithTimeout(60*time.Second),
	)

	if client.httpClient.Timeout != 60*time.Second {
		t.Errorf("WithTimeout did not set timeout, got %v", client.httpClient.Timeout)
	}
}

func TestWithRateLimit(t *testing.T) {
	client := NewClient("token", "email",
		WithRateLimit(5, 2),
	)

	// Test that rate limiter is set
	if client.rateLimiter == nil {
		t.Error("WithRateLimit did not set rate limiter")
	}

	// Test rate limiting behavior
	ctx := context.Background()
	start := time.Now()

	// Make 3 requests (burst is 2, so 3rd should wait)
	for i := 0; i < 3; i++ {
		if err := client.rateLimiter.Wait(ctx); err != nil {
			t.Fatalf("rate limiter error: %v", err)
		}
	}

	elapsed := time.Since(start)
	// With rate of 5 req/sec and burst of 2, the 3rd request should wait ~200ms
	// We'll be lenient and check for >100ms
	if elapsed < 100*time.Millisecond {
		t.Errorf("rate limiter did not throttle, elapsed: %v", elapsed)
	}
}

func TestClient_RateLimiting(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	// Create client with very restrictive rate limit
	client := NewClient("token", "email",
		WithBaseURL(server.URL+"/"),
		WithRateLimit(2, 1), // 2 req/sec, burst 1
	)

	ctx := context.Background()
	start := time.Now()

	// Make 3 requests
	for i := 0; i < 3; i++ {
		resp, err := client.doRequest(ctx, "GET", "/test", nil)
		if err != nil {
			t.Fatalf("doRequest() error = %v", err)
		}
		_ = resp.Body.Close()
	}

	elapsed := time.Since(start)

	// With 2 req/sec rate, 3 requests should take at least 1 second
	// (burst allows 1 immediate, then need to wait 0.5s for each of the next 2)
	if elapsed < 900*time.Millisecond {
		t.Errorf("rate limiting did not work as expected, elapsed: %v", elapsed)
	}

	if requestCount != 3 {
		t.Errorf("expected 3 requests, got %d", requestCount)
	}
}

// testLogger is a test implementation of Logger
type testLogger struct {
	debugMessages []string
	errorMessages []string
}

func (l *testLogger) Debugf(format string, args ...interface{}) {
	l.debugMessages = append(l.debugMessages, format)
}

func (l *testLogger) Errorf(format string, args ...interface{}) {
	l.errorMessages = append(l.errorMessages, format)
}

func TestClient_WithLogger(t *testing.T) {
	logger := &testLogger{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	client := NewClient("token", "email",
		WithBaseURL(server.URL+"/"),
		WithLogger(logger),
	)

	ctx := context.Background()
	resp, err := client.doRequest(ctx, "GET", "/test", nil)
	if err != nil {
		t.Fatalf("doRequest() error = %v", err)
	}
	defer resp.Body.Close()

	var result interface{}
	if err := client.decodeResponse(resp, &result); err != nil {
		t.Fatalf("decodeResponse() error = %v", err)
	}

	// Logger should have been called
	if len(logger.debugMessages) == 0 {
		t.Error("expected debug messages to be logged")
	}
}

func TestClient_doRequest_NetworkError(t *testing.T) {
	// Test network error handling by creating a server that closes immediately
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Close connection immediately to force network error
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Fatal("webserver doesn't support hijacking")
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			t.Fatal(err)
		}
		if err := conn.Close(); err != nil {
			t.Fatalf("close hijacked conn: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient("token", "email",
		WithBaseURL(server.URL+"/"),
		WithRetry(0, time.Millisecond), // No retries
	)

	ctx := context.Background()
	_, err := client.doRequest(ctx, "GET", "/test", nil)
	if err == nil {
		t.Error("expected network error")
		return
	}

	// Should get some kind of error (NetworkError or other error type)
	t.Logf("got expected error: %v", err)
}

func TestClient_Constants(t *testing.T) {
	// Verify constants are set correctly
	if DefaultBaseURL != "https://manapool.com/api/v1/" {
		t.Errorf("DefaultBaseURL = %q, want %q", DefaultBaseURL, "https://manapool.com/api/v1/")
	}
	if DefaultTimeout != 30*time.Second {
		t.Errorf("DefaultTimeout = %v, want %v", DefaultTimeout, 30*time.Second)
	}
	if DefaultRateLimit != 10.0 {
		t.Errorf("DefaultRateLimit = %v, want %v", DefaultRateLimit, 10.0)
	}
	if DefaultRateBurst != 1 {
		t.Errorf("DefaultRateBurst = %v, want %v", DefaultRateBurst, 1)
	}
	if DefaultMaxRetries != 3 {
		t.Errorf("DefaultMaxRetries = %v, want %v", DefaultMaxRetries, 3)
	}
	if DefaultInitialBackoff != 1*time.Second {
		t.Errorf("DefaultInitialBackoff = %v, want %v", DefaultInitialBackoff, 1*time.Second)
	}
	if Version != "0.2.0" {
		t.Errorf("Version = %q, want %q", Version, "0.2.0")
	}
}

func TestClient_RateLimiter_NotNil(t *testing.T) {
	client := NewClient("token", "email")
	if client.rateLimiter == nil {
		t.Error("rateLimiter should not be nil")
	}

	// Verify it's a valid rate limiter
	if _, ok := interface{}(client.rateLimiter).(*rate.Limiter); !ok {
		t.Error("rateLimiter is not a *rate.Limiter")
	}
}

func TestClient_decodeResponse_ReadBodyError(t *testing.T) {
	// Test error reading body
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       &errorReader{},
	}

	client := NewClient("token", "email")

	var result interface{}
	err := client.decodeResponse(resp, &result)
	if err == nil {
		t.Error("decodeResponse() expected error for unreadable body")
	}
}

// errorReader always returns an error when read
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("read error")
}

func (e *errorReader) Close() error {
	return nil
}

func TestWithTimeout_NilHTTPClient(t *testing.T) {
	// Test WithTimeout when httpClient is nil (creates new one)
	client := &Client{} // No httpClient set
	opt := WithTimeout(60 * time.Second)
	opt(client)

	if client.httpClient == nil {
		t.Error("WithTimeout should create httpClient if nil")
	}
	if client.httpClient.Timeout != 60*time.Second {
		t.Errorf("Timeout = %v, want %v", client.httpClient.Timeout, 60*time.Second)
	}
}

func TestClient_doRequest_FailedAfterRetries(t *testing.T) {
	// Test that we return error after all retries are exhausted
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("server error"))
	}))
	defer server.Close()

	client := NewClient("token", "email",
		WithBaseURL(server.URL+"/"),
		WithRetry(2, 10*time.Millisecond),
	)

	ctx := context.Background()
	resp, err := client.doRequest(ctx, "GET", "/test", nil)
	if err != nil {
		t.Fatalf("doRequest should not return error, got: %v", err)
	}

	// Should have made 3 attempts (initial + 2 retries)
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}

	// Response should be the final server error
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want 500", resp.StatusCode)
	}
	_ = resp.Body.Close()
}
