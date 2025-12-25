package manapool

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_OrderEndpoints(t *testing.T) {
	orderSummaryJSON := `{"id":"f4b3b9c5-250e-4c3a-815d-6e41577f28e3","created_at":"2024-04-01T05:44:13.336106Z","label":"1234","total_cents":1100,"shipping_method":"first_class","latest_fulfillment_status":"processing"}`
	orderDetailsJSON := `{"order":{"id":"f4b3b9c5-250e-4c3a-815d-6e41577f28e3","created_at":"2024-04-01T05:44:13.336106Z","label":"1234","total_cents":1100,"shipping_method":"first_class","latest_fulfillment_status":"processing","buyer_id":"buyer","shipping_address":{"name":"John","line1":"123","city":"City","state":"CA","postal_code":"12345","country":"US"},"payment":{"subtotal_cents":1000,"shipping_cents":100,"total_cents":1100,"fee_cents":50,"net_cents":1050},"fulfillments":[],"items":[]}}`
	fulfillmentJSON := `{"fulfillment":{"status":"shipped","tracking_company":"ups","tracking_number":"1Z","tracking_url":"http://example","in_transit_at":"2024-04-01T05:44:13.336106Z","estimated_delivery_at":"2024-04-02T05:44:13.336106Z","delivered_at":null}}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/orders":
			if r.URL.Query().Get("label") != "1234" {
				t.Fatalf("label = %q, want 1234", r.URL.Query().Get("label"))
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"orders":[` + orderSummaryJSON + `]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/orders/abc":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(orderDetailsJSON))
		case r.Method == http.MethodPut && r.URL.Path == "/orders/abc/fulfillment":
			var payload OrderFulfillmentRequest
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode fulfillment payload: %v", err)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(fulfillmentJSON))
		case r.Method == http.MethodGet && r.URL.Path == "/seller/orders":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"orders":[` + orderSummaryJSON + `]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/seller/orders/abc":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(orderDetailsJSON))
		case r.Method == http.MethodPut && r.URL.Path == "/seller/orders/abc/fulfillment":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(fulfillmentJSON))
		case r.Method == http.MethodGet && r.URL.Path == "/seller/orders/abc/reports":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"reports":[{"report_id":"rep","order_id":"ord","order_reported_issues":{"comment":null,"created_at":"2024-04-01T05:44:13.336106Z","proposed_remediation_method":null,"reporter_role":"buyer","is_nondelivery_report":false,"rescinded":false,"items":[{"order_item_id":"item","quantity":1}],"remediations":[{"remediation_expense_cents":null,"comment":null,"created_at":"2024-04-01T05:44:13.336106Z"}],"charges":[{"seller_charge_cents":null,"payout_id":null}]}}]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com", WithBaseURL(server.URL+"/"))
	ctx := context.Background()
	label := "1234"
	since := Timestamp{Time: time.Date(2024, 4, 1, 5, 44, 13, 0, time.UTC)}
	unfulfilled := true
	t.Run("GetOrders", func(t *testing.T) {
		orders, err := client.GetOrders(ctx, OrdersOptions{Since: &since, Label: label, IsUnfulfilled: &unfulfilled})
		if err != nil {
			t.Fatalf("GetOrders error: %v", err)
		}
		if len(orders.Orders) != 1 {
			t.Fatalf("orders count = %d, want 1", len(orders.Orders))
		}
	})

	t.Run("GetOrder", func(t *testing.T) {
		order, err := client.GetOrder(ctx, "abc")
		if err != nil {
			t.Fatalf("GetOrder error: %v", err)
		}
		if order.Order.Label != "1234" {
			t.Fatalf("order label = %q, want 1234", order.Order.Label)
		}
	})

	t.Run("UpdateOrderFulfillment", func(t *testing.T) {
		status := "shipped"
		fulfillment, err := client.UpdateOrderFulfillment(ctx, "abc", OrderFulfillmentRequest{Status: &status})
		if err != nil {
			t.Fatalf("UpdateOrderFulfillment error: %v", err)
		}
		if fulfillment.Fulfillment.Status == nil || *fulfillment.Fulfillment.Status != "shipped" {
			t.Fatalf("fulfillment status mismatch")
		}
	})

	t.Run("GetSellerOrders", func(t *testing.T) {
		_, err := client.GetSellerOrders(ctx, OrdersOptions{})
		if err != nil {
			t.Fatalf("GetSellerOrders error: %v", err)
		}
	})

	t.Run("GetSellerOrder", func(t *testing.T) {
		_, err := client.GetSellerOrder(ctx, "abc")
		if err != nil {
			t.Fatalf("GetSellerOrder error: %v", err)
		}
	})

	t.Run("UpdateSellerOrderFulfillment", func(t *testing.T) {
		status := "shipped"
		_, err := client.UpdateSellerOrderFulfillment(ctx, "abc", OrderFulfillmentRequest{Status: &status})
		if err != nil {
			t.Fatalf("UpdateSellerOrderFulfillment error: %v", err)
		}
	})

	t.Run("GetSellerOrderReports", func(t *testing.T) {
		reports, err := client.GetSellerOrderReports(ctx, "abc")
		if err != nil {
			t.Fatalf("GetSellerOrderReports error: %v", err)
		}
		if len(reports.Reports) != 1 {
			t.Fatalf("reports count = %d, want 1", len(reports.Reports))
		}
	})
}
