package manapool

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetSellerAccount_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method and path
		if r.Method != "GET" {
			t.Errorf("Method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/account" {
			t.Errorf("Path = %s, want /account", r.URL.Path)
		}

		// Verify headers
		if got := r.Header.Get("X-ManaPool-Access-Token"); got != "test-token" {
			t.Errorf("X-ManaPool-Access-Token = %q, want %q", got, "test-token")
		}
		if got := r.Header.Get("X-ManaPool-Email"); got != "test@example.com" {
			t.Errorf("X-ManaPool-Email = %q, want %q", got, "test@example.com")
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"username": "testuser",
			"email": "test@example.com",
			"verified": true,
			"singles_live": true,
			"sealed_live": false,
			"payouts_enabled": true
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx := context.Background()
	account, err := client.GetSellerAccount(ctx)
	if err != nil {
		t.Fatalf("GetSellerAccount() error = %v", err)
	}

	// Verify account data
	if account.Username != "testuser" {
		t.Errorf("Username = %q, want %q", account.Username, "testuser")
	}
	if account.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", account.Email, "test@example.com")
	}
	if !account.Verified {
		t.Error("Verified = false, want true")
	}
	if !account.SinglesLive {
		t.Error("SinglesLive = false, want true")
	}
	if account.SealedLive {
		t.Error("SealedLive = true, want false")
	}
	if !account.PayoutsEnabled {
		t.Error("PayoutsEnabled = false, want true")
	}
}

func TestClient_GetSellerAccount_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "account not found"}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx := context.Background()
	_, err := client.GetSellerAccount(ctx)
	if err == nil {
		t.Fatal("GetSellerAccount() expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if !apiErr.IsNotFound() {
		t.Error("expected IsNotFound() to be true")
	}
}

func TestClient_GetSellerAccount_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "unauthorized"}`))
	}))
	defer server.Close()

	client := NewClient("invalid-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx := context.Background()
	_, err := client.GetSellerAccount(ctx)
	if err == nil {
		t.Fatal("GetSellerAccount() expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if !apiErr.IsUnauthorized() {
		t.Error("expected IsUnauthorized() to be true")
	}
}

func TestClient_GetSellerAccount_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`Internal Server Error`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
		WithRetry(0, 0), // No retries for faster test
	)

	ctx := context.Background()
	_, err := client.GetSellerAccount(ctx)
	if err == nil {
		t.Fatal("GetSellerAccount() expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if !apiErr.IsServerError() {
		t.Error("expected IsServerError() to be true")
	}
}

func TestClient_GetSellerAccount_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid json}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx := context.Background()
	_, err := client.GetSellerAccount(ctx)
	if err == nil {
		t.Fatal("GetSellerAccount() expected error for invalid JSON, got nil")
	}
}

func TestClient_GetSellerAccount_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This should not be reached
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.GetSellerAccount(ctx)
	if err == nil {
		t.Fatal("GetSellerAccount() expected error for cancelled context, got nil")
	}
}

func TestClient_GetSellerAccount_WithLogger(t *testing.T) {
	logger := &testLogger{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"username": "testuser",
			"email": "test@example.com",
			"verified": true,
			"singles_live": true,
			"sealed_live": false,
			"payouts_enabled": true
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
		WithLogger(logger),
	)

	ctx := context.Background()
	_, err := client.GetSellerAccount(ctx)
	if err != nil {
		t.Fatalf("GetSellerAccount() error = %v", err)
	}

	// Verify logger was called
	if len(logger.debugMessages) == 0 {
		t.Error("expected debug messages to be logged")
	}
}
