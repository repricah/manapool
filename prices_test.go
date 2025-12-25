package manapool

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetPricesEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/prices/singles":
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want GET", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"meta":{"as_of":"2024-04-01T05:44:13.336106Z"},"data":[{"url":"https://manapool.com/card/ice/89/polar-kraken","name":"Polar Kraken","set_code":"ICE","number":"89","multiverse_id":null,"scryfall_id":"aee01e9c-0445-4228-a73a-3e5744844ed3","available_quantity":2,"price_cents":123,"price_cents_lp_plus":null,"price_cents_nm":null,"price_cents_foil":null,"price_cents_lp_plus_foil":null,"price_cents_nm_foil":null,"price_cents_etched":null,"price_cents_lp_plus_etched":null,"price_cents_nm_etched":null}]}`))
		case "/prices/variants":
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want GET", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"meta":{"as_of":"2024-04-01T05:44:13.336106Z"},"data":[{"url":"https://manapool.com/card/ice/89/polar-kraken","product_type":"mtg_single","product_id":"123e4567-e89b-12d3-a456-426614174000","set_code":"ICE","number":"89","name":"Polar Kraken","scryfall_id":"aee01e9c-0445-4228-a73a-3e5744844ed3","tcgplayer_product_id":123,"language_id":"EN","condition_id":"NM","finish_id":"FO","low_price":1999,"available_quantity":5}]}`))
		case "/prices/sealed":
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want GET", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"meta":{"as_of":"2024-04-01T05:44:13.336106Z"},"data":[{"url":"https://manapool.com/sealed/ice/box","product_type":"mtg_sealed","product_id":"123e4567-e89b-12d3-a456-426614174000","set_code":"ICE","name":"Ice Age Booster Box","tcgplayer_product_id":321,"language_id":"EN","low_price":2999,"available_quantity":3}]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com", WithBaseURL(server.URL+"/"))
	ctx := context.Background()

	t.Run("GetSinglesPrices", func(t *testing.T) {
		singles, err := client.GetSinglesPrices(ctx)
		if err != nil {
			t.Fatalf("GetSinglesPrices error: %v", err)
		}
		if len(singles.Data) != 1 {
			t.Fatalf("singles count = %d, want 1", len(singles.Data))
		}
	})

	t.Run("GetVariantPrices", func(t *testing.T) {
		variants, err := client.GetVariantPrices(ctx)
		if err != nil {
			t.Fatalf("GetVariantPrices error: %v", err)
		}
		if len(variants.Data) != 1 {
			t.Fatalf("variants count = %d, want 1", len(variants.Data))
		}
	})

	t.Run("GetSealedPrices", func(t *testing.T) {
		sealed, err := client.GetSealedPrices(ctx)
		if err != nil {
			t.Fatalf("GetSealedPrices error: %v", err)
		}
		if len(sealed.Data) != 1 {
			t.Fatalf("sealed count = %d, want 1", len(sealed.Data))
		}
	})
}
