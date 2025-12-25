package manapool

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_BuyerEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/buyer/optimizer":
			var req OptimizerRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode optimizer request: %v", err)
			}
			if len(req.Cart) != 1 {
				t.Fatalf("cart items = %d, want 1", len(req.Cart))
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"cart":[{"inventory_id":"aee01e9c-0445-4228-a73a-3e5744844ed3","quantity_selected":1}],"totals":{"subtotal_cents":1000,"shipping_cents":500,"total_cents":1500,"seller_count":1}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/buyer/orders":
			if r.URL.Query().Get("since") == "" {
				t.Fatalf("missing since query param")
			}
			if r.URL.Query().Get("limit") != "2" {
				t.Fatalf("limit = %q, want 2", r.URL.Query().Get("limit"))
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"orders":[{"id":"f4b3b9c5-250e-4c3a-815d-6e41577f28e3","created_at":"2024-04-01T05:44:13.336106Z","subtotal_cents":1000,"tax_cents":100,"shipping_cents":100,"total_cents":1200,"order_number":"1234","order_seller_details":[{"order_number":"1234-5678","seller_id":"f4b3b9c5-250e-4c3a-815d-6e41577f28e3","seller_username":"seller","item_count":1}]}]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/buyer/orders/abc":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"order":{"id":"f4b3b9c5-250e-4c3a-815d-6e41577f28e3","created_at":"2024-04-01T05:44:13.336106Z","subtotal_cents":1000,"tax_cents":100,"shipping_cents":100,"total_cents":1200,"order_number":"1234","order_seller_details":[{"order_number":"1234-5678","seller_id":"f4b3b9c5-250e-4c3a-815d-6e41577f28e3","seller_username":"seller","fulfillments":[{"status":"shipped","tracking_company":"ups","tracking_number":"1Z","tracking_url":"http://example","in_transit_at":"2024-04-01T05:44:13.336106Z"}],"items":[{"price_cents":1000,"quantity":1,"product":{"product_type":"mtg_single","product_id":"f4b3b9c5-250e-4c3a-815d-6e41577f28e3","single":{"scryfall_id":"a","mtgjson_id":"b","name":"Card","set":"ICE","number":"1","language_id":"EN","condition_id":"NM","finish_id":"NF"},"sealed":null}}]}]}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/buyer/orders/pending-orders":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"123e4567-e89b-12d3-a456-426614174000","line_items":[{"inventory_id":"inv","quantity_selected":1}],"status":"pending","totals":{"subtotal_cents":1000,"shipping_cents":500,"tax_cents":100,"total_cents":1600},"order":null}`))
		case r.Method == http.MethodGet && r.URL.Path == "/buyer/orders/pending-orders/123":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"123","line_items":[{"inventory_id":"inv","quantity_selected":1}],"status":"pending","totals":{"subtotal_cents":1000,"shipping_cents":500,"tax_cents":100,"total_cents":1600},"order":null}`))
		case r.Method == http.MethodPut && r.URL.Path == "/buyer/orders/pending-orders/123":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"123","line_items":[{"inventory_id":"inv","quantity_selected":2}],"status":"pending","totals":{"subtotal_cents":2000,"shipping_cents":500,"tax_cents":100,"total_cents":2600},"order":null}`))
		case r.Method == http.MethodPost && r.URL.Path == "/buyer/orders/pending-orders/123/purchase":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"123","line_items":[{"inventory_id":"inv","quantity_selected":2}],"status":"completed","totals":{"subtotal_cents":2000,"shipping_cents":500,"tax_cents":100,"total_cents":2600},"order":{"id":"ord"}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/buyer/credit":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"user_credit_cents":500}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com", WithBaseURL(server.URL+"/"))
	ctx := context.Background()

	t.Run("OptimizeCart", func(t *testing.T) {
		optimized, err := client.OptimizeCart(ctx, OptimizerRequest{Cart: []OptimizerCartItem{{Type: "mtg_single", Name: "Card", QuantityRequested: 1}}})
		if err != nil {
			t.Fatalf("OptimizeCart error: %v", err)
		}
		if optimized.Totals.TotalCents != 1500 {
			t.Fatalf("total cents = %d, want 1500", optimized.Totals.TotalCents)
		}
	})

	t.Run("GetBuyerOrders", func(t *testing.T) {
		since := Timestamp{Time: time.Date(2024, 4, 1, 5, 44, 13, 0, time.UTC)}
		orders, err := client.GetBuyerOrders(ctx, BuyerOrdersOptions{Since: &since, Limit: 2, Offset: 0})
		if err != nil {
			t.Fatalf("GetBuyerOrders error: %v", err)
		}
		if len(orders.Orders) != 1 {
			t.Fatalf("orders count = %d, want 1", len(orders.Orders))
		}
	})

	t.Run("GetBuyerOrder", func(t *testing.T) {
		order, err := client.GetBuyerOrder(ctx, "abc")
		if err != nil {
			t.Fatalf("GetBuyerOrder error: %v", err)
		}
		if order.Order.OrderNumber != "1234" {
			t.Fatalf("order number = %q, want 1234", order.Order.OrderNumber)
		}
	})

	t.Run("PendingOrders", func(t *testing.T) {
		pending, err := client.CreatePendingOrder(ctx, PendingOrderRequest{LineItems: []PendingOrderLineItem{{InventoryID: "inv", QuantitySelected: 1}}})
		if err != nil {
			t.Fatalf("CreatePendingOrder error: %v", err)
		}
		if pending.Status != "pending" {
			t.Fatalf("pending status = %q, want pending", pending.Status)
		}

		pending, err = client.GetPendingOrder(ctx, "123")
		if err != nil {
			t.Fatalf("GetPendingOrder error: %v", err)
		}
		if pending.ID != "123" {
			t.Fatalf("pending id = %q, want 123", pending.ID)
		}

		pending, err = client.UpdatePendingOrder(ctx, "123", PendingOrderRequest{LineItems: []PendingOrderLineItem{{InventoryID: "inv", QuantitySelected: 2}}})
		if err != nil {
			t.Fatalf("UpdatePendingOrder error: %v", err)
		}
		if pending.Totals.SubtotalCents != 2000 {
			t.Fatalf("pending subtotal = %d, want 2000", pending.Totals.SubtotalCents)
		}

		pending, err = client.PurchasePendingOrder(ctx, "123", PurchasePendingOrderRequest{PaymentMethod: "user_credit", BillingAddress: Address{Line1: "line", City: "City", State: "CA", PostalCode: "12345", Country: "US"}, ShippingAddress: Address{Line1: "line", City: "City", State: "CA", PostalCode: "12345", Country: "US"}})
		if err != nil {
			t.Fatalf("PurchasePendingOrder error: %v", err)
		}
		if pending.Status != "completed" {
			t.Fatalf("pending status = %q, want completed", pending.Status)
		}
	})

	t.Run("GetBuyerCredit", func(t *testing.T) {
		credit, err := client.GetBuyerCredit(ctx)
		if err != nil {
			t.Fatalf("GetBuyerCredit error: %v", err)
		}
		if credit.UserCreditCents != 500 {
			t.Fatalf("credit = %d, want 500", credit.UserCreditCents)
		}
	})
}
