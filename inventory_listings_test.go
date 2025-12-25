package manapool

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_InventoryListingEndpoints(t *testing.T) {
	inventoryJSON := `{"id":"inv123","product_type":"mtg_single","product_id":"prod456","price_cents":499,"quantity":5,"effective_as_of":"2025-08-05T20:38:54.549229Z","product":{"type":"mtg_single","id":"prod456","tcgplayer_sku":123456,"single":{"scryfall_id":"abc123","mtgjson_id":"def456","tcgplayer_id":111,"name":"Black Lotus","set":"LEA","number":"232","language_id":"EN","condition_id":"NM","finish_id":"NF"},"sealed":null}}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/inventory/listings":
			if got := r.URL.Query().Get("id"); got != "inv123" {
				t.Fatalf("id query = %q, want inv123", got)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory_items":[` + inventoryJSON + `]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/inventory/listings/inv123":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory_item":` + inventoryJSON + `}`))
		case r.Method == http.MethodGet && r.URL.Path == "/inventory/tcgsku/123456":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		case r.Method == http.MethodPut && r.URL.Path == "/inventory/tcgsku/123456":
			var payload InventoryUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode update payload: %v", err)
			}
			if payload.PriceCents == 0 {
				t.Fatalf("expected price_cents in payload")
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/inventory/tcgsku/123456":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		case r.Method == http.MethodPost && r.URL.Path == "/seller/inventory":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":[` + inventoryJSON + `]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/seller/inventory/tcgsku":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":[` + inventoryJSON + `]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/seller/inventory/tcgsku/123456":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		case r.Method == http.MethodPut && r.URL.Path == "/seller/inventory/tcgsku/123456":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/seller/inventory/tcgsku/123456":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		case r.Method == http.MethodPost && r.URL.Path == "/seller/inventory/product":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":[` + inventoryJSON + `]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/seller/inventory/product/mtg_single/prod456":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		case r.Method == http.MethodPut && r.URL.Path == "/seller/inventory/product/mtg_single/prod456":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/seller/inventory/product/mtg_single/prod456":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		case r.Method == http.MethodPost && r.URL.Path == "/seller/inventory/scryfall_id":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":[` + inventoryJSON + `]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/seller/inventory/scryfall_id/abc123":
			if r.URL.Query().Get("language_id") != "EN" {
				t.Fatalf("expected language_id query")
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		case r.Method == http.MethodPut && r.URL.Path == "/seller/inventory/scryfall_id/abc123":
			if r.URL.Query().Get("finish_id") != "NF" {
				t.Fatalf("expected finish_id query")
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/seller/inventory/scryfall_id/abc123":
			if r.URL.Query().Get("condition_id") != "NM" {
				t.Fatalf("expected condition_id query")
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		case r.Method == http.MethodPost && r.URL.Path == "/seller/inventory/tcgplayer_id":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":[` + inventoryJSON + `]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/seller/inventory/tcgplayer_id/4841":
			if r.URL.Query().Get("language_id") != "EN" {
				t.Fatalf("expected language_id query")
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		case r.Method == http.MethodPut && r.URL.Path == "/seller/inventory/tcgplayer_id/4841":
			if r.URL.Query().Get("finish_id") != "NF" {
				t.Fatalf("expected finish_id query")
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/seller/inventory/tcgplayer_id/4841":
			if r.URL.Query().Get("condition_id") != "NM" {
				t.Fatalf("expected condition_id query")
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"inventory":` + inventoryJSON + `}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com", WithBaseURL(server.URL+"/"))
	ctx := context.Background()

	t.Run("GetInventoryListings", func(t *testing.T) {
		listings, err := client.GetInventoryListings(ctx, []string{"inv123"})
		if err != nil {
			t.Fatalf("GetInventoryListings error: %v", err)
		}
		if len(listings.InventoryItems) != 1 {
			t.Fatalf("inventory items = %d, want 1", len(listings.InventoryItems))
		}

		item, err := client.GetInventoryListing(ctx, "inv123")
		if err != nil {
			t.Fatalf("GetInventoryListing error: %v", err)
		}
		if item.InventoryItem.ID != "inv123" {
			t.Fatalf("inventory id = %s, want inv123", item.InventoryItem.ID)
		}
	})

	t.Run("InventoryBySKU", func(t *testing.T) {
		_, err := client.GetInventoryBySKU(ctx, 123456)
		if err != nil {
			t.Fatalf("GetInventoryBySKU error: %v", err)
		}

		_, err = client.UpdateInventoryBySKU(ctx, 123456, InventoryUpdateRequest{PriceCents: 500, Quantity: 3})
		if err != nil {
			t.Fatalf("UpdateInventoryBySKU error: %v", err)
		}

		_, err = client.DeleteInventoryBySKU(ctx, 123456)
		if err != nil {
			t.Fatalf("DeleteInventoryBySKU error: %v", err)
		}
	})

	t.Run("BulkUpdates", func(t *testing.T) {
		_, err := client.CreateInventoryBulk(ctx, []InventoryBulkItemBySKU{{TCGPlayerSKU: 123456, PriceCents: 100, Quantity: 1}})
		if err != nil {
			t.Fatalf("CreateInventoryBulk error: %v", err)
		}

		_, err = client.CreateInventoryBulkBySKU(ctx, []InventoryBulkItemBySKU{{TCGPlayerSKU: 123456, PriceCents: 100, Quantity: 1}})
		if err != nil {
			t.Fatalf("CreateInventoryBulkBySKU error: %v", err)
		}

		_, err = client.CreateInventoryBulkByProduct(ctx, []InventoryBulkItemByProduct{{ProductType: "mtg_single", ProductID: "prod456", PriceCents: 100, Quantity: 1}})
		if err != nil {
			t.Fatalf("CreateInventoryBulkByProduct error: %v", err)
		}

		_, err = client.CreateInventoryBulkByScryfall(ctx, []InventoryBulkItemByScryfall{{ScryfallID: "abc123", LanguageID: "EN", FinishID: "NF", ConditionID: "NM", PriceCents: 100, Quantity: 1}})
		if err != nil {
			t.Fatalf("CreateInventoryBulkByScryfall error: %v", err)
		}

		finish := "NF"
		condition := "NM"
		_, err = client.CreateInventoryBulkByTCGPlayerID(ctx, []InventoryBulkItemByTCGPlayerID{{TCGPlayerID: 4841, LanguageID: "EN", FinishID: &finish, ConditionID: &condition, PriceCents: 100, Quantity: 1}})
		if err != nil {
			t.Fatalf("CreateInventoryBulkByTCGPlayerID error: %v", err)
		}
	})

	t.Run("SellerInventoryBySKU", func(t *testing.T) {
		_, err := client.GetSellerInventoryBySKU(ctx, 123456)
		if err != nil {
			t.Fatalf("GetSellerInventoryBySKU error: %v", err)
		}

		_, err = client.UpdateSellerInventoryBySKU(ctx, 123456, InventoryUpdateRequest{PriceCents: 200, Quantity: 2})
		if err != nil {
			t.Fatalf("UpdateSellerInventoryBySKU error: %v", err)
		}

		_, err = client.DeleteSellerInventoryBySKU(ctx, 123456)
		if err != nil {
			t.Fatalf("DeleteSellerInventoryBySKU error: %v", err)
		}
	})

	t.Run("SellerInventoryByProduct", func(t *testing.T) {
		_, err := client.GetSellerInventoryByProduct(ctx, "mtg_single", "prod456")
		if err != nil {
			t.Fatalf("GetSellerInventoryByProduct error: %v", err)
		}

		_, err = client.UpdateSellerInventoryByProduct(ctx, "mtg_single", "prod456", InventoryUpdateRequest{PriceCents: 200, Quantity: 1})
		if err != nil {
			t.Fatalf("UpdateSellerInventoryByProduct error: %v", err)
		}

		_, err = client.DeleteSellerInventoryByProduct(ctx, "mtg_single", "prod456")
		if err != nil {
			t.Fatalf("DeleteSellerInventoryByProduct error: %v", err)
		}
	})

	t.Run("SellerInventoryByScryfall", func(t *testing.T) {
		opts := InventoryByScryfallOptions{LanguageID: "EN", FinishID: "NF", ConditionID: "NM"}
		_, err := client.GetSellerInventoryByScryfall(ctx, "abc123", opts)
		if err != nil {
			t.Fatalf("GetSellerInventoryByScryfall error: %v", err)
		}

		_, err = client.UpdateSellerInventoryByScryfall(ctx, "abc123", opts, InventoryUpdateRequest{PriceCents: 200, Quantity: 2})
		if err != nil {
			t.Fatalf("UpdateSellerInventoryByScryfall error: %v", err)
		}

		_, err = client.DeleteSellerInventoryByScryfall(ctx, "abc123", opts)
		if err != nil {
			t.Fatalf("DeleteSellerInventoryByScryfall error: %v", err)
		}
	})

	t.Run("SellerInventoryByTCGPlayerID", func(t *testing.T) {
		optsTCG := InventoryByTCGPlayerOptions{LanguageID: "EN", FinishID: "NF", ConditionID: "NM"}
		_, err := client.GetSellerInventoryByTCGPlayerID(ctx, 4841, optsTCG)
		if err != nil {
			t.Fatalf("GetSellerInventoryByTCGPlayerID error: %v", err)
		}

		_, err = client.UpdateSellerInventoryByTCGPlayerID(ctx, 4841, optsTCG, InventoryUpdateRequest{PriceCents: 200, Quantity: 2})
		if err != nil {
			t.Fatalf("UpdateSellerInventoryByTCGPlayerID error: %v", err)
		}

		_, err = client.DeleteSellerInventoryByTCGPlayerID(ctx, 4841, optsTCG)
		if err != nil {
			t.Fatalf("DeleteSellerInventoryByTCGPlayerID error: %v", err)
		}
	})

	if !strings.Contains(inventoryJSON, "Black Lotus") {
		t.Fatalf("expected inventory json to include card name")
	}
}
