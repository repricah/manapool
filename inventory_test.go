package manapool

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetSellerInventory_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method and path
		if r.Method != "GET" {
			t.Errorf("Method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/seller/inventory" {
			t.Errorf("Path = %s, want /seller/inventory", r.URL.Path)
		}

		// Verify query parameters
		if got := r.URL.Query().Get("limit"); got != "100" {
			t.Errorf("limit = %q, want %q", got, "100")
		}
		if got := r.URL.Query().Get("offset"); got != "50" {
			t.Errorf("offset = %q, want %q", got, "50")
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"inventory": [
				{
					"id": "inv123",
					"product_type": "single",
					"product_id": "prod456",
					"price_cents": 499,
					"quantity": 5,
					"effective_as_of": "2025-08-05T20:38:54.549229Z",
					"product": {
						"type": "single",
						"id": "prod456",
						"tcgplayer_sku": 123456,
						"single": {
							"scryfall_id": "abc123",
							"mtgjson_id": "def456",
							"name": "Black Lotus",
							"set": "LEA",
							"number": "232",
							"language_id": "EN",
							"condition_id": "NM",
							"finish_id": "NF"
						},
						"sealed": {}
					}
				}
			],
			"pagination": {
				"total": 1000,
				"returned": 1,
				"offset": 50,
				"limit": 100
			}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx := context.Background()
	opts := InventoryOptions{
		Limit:  100,
		Offset: 50,
	}

	resp, err := client.GetSellerInventory(ctx, opts)
	if err != nil {
		t.Fatalf("GetSellerInventory() error = %v", err)
	}

	// Verify response
	if len(resp.Inventory) != 1 {
		t.Errorf("len(Inventory) = %d, want 1", len(resp.Inventory))
	}

	item := resp.Inventory[0]
	if item.ID != "inv123" {
		t.Errorf("ID = %q, want %q", item.ID, "inv123")
	}
	if item.PriceCents != 499 {
		t.Errorf("PriceCents = %d, want 499", item.PriceCents)
	}
	if item.Quantity != 5 {
		t.Errorf("Quantity = %d, want 5", item.Quantity)
	}
	if item.Product.Single.Name != "Black Lotus" {
		t.Errorf("Product.Single.Name = %q, want %q", item.Product.Single.Name, "Black Lotus")
	}

	// Verify pagination
	if resp.Pagination.Total != 1000 {
		t.Errorf("Pagination.Total = %d, want 1000", resp.Pagination.Total)
	}
	if resp.Pagination.Returned != 1 {
		t.Errorf("Pagination.Returned = %d, want 1", resp.Pagination.Returned)
	}
	if resp.Pagination.Offset != 50 {
		t.Errorf("Pagination.Offset = %d, want 50", resp.Pagination.Offset)
	}
	if resp.Pagination.Limit != 100 {
		t.Errorf("Pagination.Limit = %d, want 100", resp.Pagination.Limit)
	}
}

func TestClient_GetSellerInventory_DefaultLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify default limit is set to 500
		if got := r.URL.Query().Get("limit"); got != "500" {
			t.Errorf("limit = %q, want %q (default)", got, "500")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"inventory": [],
			"pagination": {"total": 0, "returned": 0, "offset": 0, "limit": 500}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx := context.Background()
	opts := InventoryOptions{
		Limit:  0, // Should default to 500
		Offset: 0,
	}

	_, err := client.GetSellerInventory(ctx, opts)
	if err != nil {
		t.Fatalf("GetSellerInventory() error = %v", err)
	}
}

func TestClient_GetSellerInventory_ValidationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"inventory": [],
			"pagination": {"total": 0, "returned": 0, "offset": 0, "limit": 500}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	tests := []struct {
		name    string
		opts    InventoryOptions
		wantErr bool
	}{
		{
			name:    "negative limit",
			opts:    InventoryOptions{Limit: -1, Offset: 0},
			wantErr: true,
		},
		{
			name:    "limit exceeds max",
			opts:    InventoryOptions{Limit: 501, Offset: 0},
			wantErr: true,
		},
		{
			name:    "negative offset",
			opts:    InventoryOptions{Limit: 100, Offset: -1},
			wantErr: true,
		},
		{
			name:    "valid options",
			opts:    InventoryOptions{Limit: 100, Offset: 0},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := client.GetSellerInventory(ctx, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSellerInventory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetInventoryByTCGPlayerID_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method and path
		if r.Method != "GET" {
			t.Errorf("Method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/seller/inventory/tcgsku/4549403" {
			t.Errorf("Path = %s, want /seller/inventory/tcgsku/4549403", r.URL.Path)
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": "inv123",
			"product_type": "single",
			"product_id": "prod456",
			"price_cents": 999,
			"quantity": 3,
			"effective_as_of": "2025-08-05T20:38:54.549229Z",
			"product": {
				"type": "single",
				"id": "prod456",
				"tcgplayer_sku": 4549403,
				"single": {
					"scryfall_id": "abc123",
					"mtgjson_id": "def456",
					"name": "Lightning Bolt",
					"set": "LEA",
					"number": "161",
					"language_id": "EN",
					"condition_id": "NM",
					"finish_id": "FO"
				},
				"sealed": {}
			}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx := context.Background()
	item, err := client.GetInventoryByTCGPlayerID(ctx, "4549403")
	if err != nil {
		t.Fatalf("GetInventoryByTCGPlayerID() error = %v", err)
	}

	// Verify item
	if item.ID != "inv123" {
		t.Errorf("ID = %q, want %q", item.ID, "inv123")
	}
	if item.PriceCents != 999 {
		t.Errorf("PriceCents = %d, want 999", item.PriceCents)
	}
	if item.Quantity != 3 {
		t.Errorf("Quantity = %d, want 3", item.Quantity)
	}
	if item.Product.TCGPlayerSKU != 4549403 {
		t.Errorf("Product.TCGPlayerSKU = %d, want 4549403", item.Product.TCGPlayerSKU)
	}
	if item.Product.Single.Name != "Lightning Bolt" {
		t.Errorf("Product.Single.Name = %q, want %q", item.Product.Single.Name, "Lightning Bolt")
	}
}

func TestClient_GetInventoryByTCGPlayerID_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx := context.Background()
	_, err := client.GetInventoryByTCGPlayerID(ctx, "999999")
	if err == nil {
		t.Fatal("GetInventoryByTCGPlayerID() expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if !apiErr.IsNotFound() {
		t.Error("expected IsNotFound() to be true")
	}
}

func TestClient_GetInventoryByTCGPlayerID_EmptyID(t *testing.T) {
	client := NewClient("test-token", "test@example.com")

	ctx := context.Background()
	_, err := client.GetInventoryByTCGPlayerID(ctx, "")
	if err == nil {
		t.Fatal("GetInventoryByTCGPlayerID() expected error for empty ID, got nil")
	}

	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
}

func TestIterateInventory_Success(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		offset := r.URL.Query().Get("offset")

		var response string
		if offset == "0" {
			// First page
			response = `{
				"inventory": [
					{
						"id": "inv1",
						"product_type": "single",
						"product_id": "prod1",
						"price_cents": 100,
						"quantity": 1,
						"effective_as_of": "2025-08-05T20:38:54.549229Z",
						"product": {
							"type": "single",
							"id": "prod1",
							"tcgplayer_sku": 1,
							"single": {"name": "Card 1", "condition_id": "NM", "finish_id": "NF"},
							"sealed": {}
						}
					},
					{
						"id": "inv2",
						"product_type": "single",
						"product_id": "prod2",
						"price_cents": 200,
						"quantity": 2,
						"effective_as_of": "2025-08-05T20:38:54.549229Z",
						"product": {
							"type": "single",
							"id": "prod2",
							"tcgplayer_sku": 2,
							"single": {"name": "Card 2", "condition_id": "NM", "finish_id": "NF"},
							"sealed": {}
						}
					}
				],
				"pagination": {"total": 3, "returned": 2, "offset": 0, "limit": 500}
			}`
		} else {
			// Second page (last item)
			response = `{
				"inventory": [
					{
						"id": "inv3",
						"product_type": "single",
						"product_id": "prod3",
						"price_cents": 300,
						"quantity": 3,
						"effective_as_of": "2025-08-05T20:38:54.549229Z",
						"product": {
							"type": "single",
							"id": "prod3",
							"tcgplayer_sku": 3,
							"single": {"name": "Card 3", "condition_id": "NM", "finish_id": "NF"},
							"sealed": {}
						}
					}
				],
				"pagination": {"total": 3, "returned": 1, "offset": 2, "limit": 500}
			}`
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx := context.Background()
	itemCount := 0
	err := IterateInventory(ctx, client, func(item *InventoryItem) error {
		itemCount++
		return nil
	})

	if err != nil {
		t.Fatalf("IterateInventory() error = %v", err)
	}

	if itemCount != 3 {
		t.Errorf("itemCount = %d, want 3", itemCount)
	}

	if callCount != 2 {
		t.Errorf("API called %d times, want 2", callCount)
	}
}

func TestIterateInventory_CallbackError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"inventory": [
				{
					"id": "inv1",
					"product_type": "single",
					"product_id": "prod1",
					"price_cents": 100,
					"quantity": 1,
					"effective_as_of": "2025-08-05T20:38:54.549229Z",
					"product": {
						"type": "single",
						"id": "prod1",
						"tcgplayer_sku": 1,
						"single": {"name": "Card 1", "condition_id": "NM", "finish_id": "NF"},
						"sealed": {}
					}
				}
			],
			"pagination": {"total": 1, "returned": 1, "offset": 0, "limit": 500}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx := context.Background()
	expectedErr := errors.New("callback error")
	err := IterateInventory(ctx, client, func(item *InventoryItem) error {
		return expectedErr
	})

	if err == nil {
		t.Fatal("IterateInventory() expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("IterateInventory() error = %v, want %v", err, expectedErr)
	}
}

func TestIterateInventory_APIError(t *testing.T) {
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
	err := IterateInventory(ctx, client, func(item *InventoryItem) error {
		return nil
	})

	if err == nil {
		t.Fatal("IterateInventory() expected error, got nil")
	}
}

func TestIterateInventory_EmptyInventory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"inventory": [],
			"pagination": {"total": 0, "returned": 0, "offset": 0, "limit": 500}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx := context.Background()
	callbackCalled := false
	err := IterateInventory(ctx, client, func(item *InventoryItem) error {
		callbackCalled = true
		return nil
	})

	if err != nil {
		t.Fatalf("IterateInventory() error = %v", err)
	}

	if callbackCalled {
		t.Error("callback should not be called for empty inventory")
	}
}

func TestIterateInventory_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should not be reached
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := IterateInventory(ctx, client, func(item *InventoryItem) error {
		return nil
	})

	if err == nil {
		t.Fatal("IterateInventory() expected error for cancelled context, got nil")
	}
}

func TestClient_GetSellerInventory_WithLogger(t *testing.T) {
	logger := &testLogger{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"inventory": [],
			"pagination": {"total": 0, "returned": 0, "offset": 0, "limit": 500}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
		WithLogger(logger),
	)

	ctx := context.Background()
	opts := InventoryOptions{Limit: 100, Offset: 0}
	_, err := client.GetSellerInventory(ctx, opts)
	if err != nil {
		t.Fatalf("GetSellerInventory() error = %v", err)
	}

	// Verify logger was called
	if len(logger.debugMessages) == 0 {
		t.Error("expected debug messages to be logged")
	}
}

func TestInventoryItem_PriceDollars_Helper(t *testing.T) {
	// Test helper method
	item := InventoryItem{PriceCents: 1234}
	if got := item.PriceDollars(); got != 12.34 {
		t.Errorf("PriceDollars() = %v, want 12.34", got)
	}
}

func TestIterateInventory_LargeDataset(t *testing.T) {
	// Test with multiple pages to verify pagination logic
	pageSize := 500
	totalItems := 1250 // Should result in 3 pages

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset := r.URL.Query().Get("offset")
		var offsetInt int
		if _, err := fmt.Sscanf(offset, "%d", &offsetInt); err != nil {
			t.Fatalf("parse offset %q: %v", offset, err)
		}

		remaining := totalItems - offsetInt
		if remaining > pageSize {
			remaining = pageSize
		}

		// Generate mock items
		response := fmt.Sprintf(`{
			"inventory": [%s],
			"pagination": {"total": %d, "returned": %d, "offset": %d, "limit": 500}
		}`, generateMockItems(remaining), totalItems, remaining, offsetInt)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com",
		WithBaseURL(server.URL+"/"),
	)

	ctx := context.Background()
	itemCount := 0
	err := IterateInventory(ctx, client, func(item *InventoryItem) error {
		itemCount++
		return nil
	})

	if err != nil {
		t.Fatalf("IterateInventory() error = %v", err)
	}

	if itemCount != totalItems {
		t.Errorf("itemCount = %d, want %d", itemCount, totalItems)
	}
}

// Helper function to generate mock inventory items for testing
func generateMockItems(count int) string {
	if count == 0 {
		return ""
	}

	items := make([]string, count)
	for i := 0; i < count; i++ {
		items[i] = fmt.Sprintf(`{
			"id": "inv%d",
			"product_type": "single",
			"product_id": "prod%d",
			"price_cents": %d,
			"quantity": 1,
			"effective_as_of": "2025-08-05T20:38:54.549229Z",
			"product": {
				"type": "single",
				"id": "prod%d",
				"tcgplayer_sku": %d,
				"single": {"name": "Card %d", "condition_id": "NM", "finish_id": "NF"},
				"sealed": {}
			}
		}`, i, i, 100*(i+1), i, i, i)
	}

	result := ""
	for i, item := range items {
		if i > 0 {
			result += ","
		}
		result += item
	}
	return result
}
