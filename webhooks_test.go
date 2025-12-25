package manapool

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_WebhooksEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/webhooks":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"webhooks":[{"id":"wh","topic":"order_created","callback_url":"https://example.com"}]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/webhooks/wh":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"wh","topic":"order_created","callback_url":"https://example.com"}`))
		case r.Method == http.MethodPut && r.URL.Path == "/webhooks/register":
			var payload WebhookRegisterRequest
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode payload: %v", err)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"wh","topic":"order_created","callback_url":"https://example.com"}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/webhooks/wh":
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com", WithBaseURL(server.URL+"/"))
	ctx := context.Background()

	t.Run("GetWebhooks", func(t *testing.T) {
		webhooks, err := client.GetWebhooks(ctx, "")
		if err != nil {
			t.Fatalf("GetWebhooks error: %v", err)
		}
		if len(webhooks.Webhooks) != 1 {
			t.Fatalf("webhooks count = %d, want 1", len(webhooks.Webhooks))
		}
	})

	t.Run("GetWebhook", func(t *testing.T) {
		webhook, err := client.GetWebhook(ctx, "wh")
		if err != nil {
			t.Fatalf("GetWebhook error: %v", err)
		}
		if webhook.ID != "wh" {
			t.Fatalf("webhook id = %s, want wh", webhook.ID)
		}
	})

	t.Run("RegisterWebhook", func(t *testing.T) {
		_, err := client.RegisterWebhook(ctx, WebhookRegisterRequest{Topic: "order_created", CallbackURL: "https://example.com"})
		if err != nil {
			t.Fatalf("RegisterWebhook error: %v", err)
		}
	})

	t.Run("DeleteWebhook", func(t *testing.T) {
		if err := client.DeleteWebhook(ctx, "wh"); err != nil {
			t.Fatalf("DeleteWebhook error: %v", err)
		}
	})
}
