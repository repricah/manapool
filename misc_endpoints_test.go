package manapool

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_MiscEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/deck":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"valid":true,"buy_url":"https://manapool.com/add-deck","details":{"commander_count":1,"total_card_count":100,"all_cards_legal":true,"illegal_cards":[],"cards_not_found":[],"valid_quantities":true,"quantity_violations":[],"valid_color_identity":true,"color_identity_violations":[],"valid_partnership":true,"partner_violations":[]}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/card_info":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"cards":[{"name":"Lightning Bolt","set_code":"LEA","set_name":"Alpha","card_number":"1","rarity":"common","from_price_cents":100,"quantity_available":5,"release_date":"1993-08-05","legal_formats":["legacy"],"flavor_name":null,"layout":null,"is_token":false,"promo_types":null,"finishes":null,"text":null,"color_identity":null,"edhrecSaltiness":null,"power":null,"defense":null,"mana_cost":null,"mana_value":null}],"not_found":[]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/job-apply":
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				t.Fatalf("parse multipart: %v", err)
			}
			if r.FormValue("first_name") != "John" {
				t.Fatalf("first_name = %q, want John", r.FormValue("first_name"))
			}
			file, _, err := r.FormFile("application")
			if err != nil {
				t.Fatalf("missing application file: %v", err)
			}
			_, _ = io.ReadAll(file)
			_ = file.Close()
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"success":true,"message":"Application submitted successfully"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test@example.com", WithBaseURL(server.URL+"/"))
	ctx := context.Background()

	t.Run("CreateDeck", func(t *testing.T) {
		deck, err := client.CreateDeck(ctx, DeckCreateRequest{CommanderNames: []string{"Atraxa"}, OtherCards: []OtherCard{{Name: "Lightning Bolt", Quantity: 4}}})
		if err != nil {
			t.Fatalf("CreateDeck error: %v", err)
		}
		if !deck.Valid {
			t.Fatalf("expected deck to be valid")
		}
	})

	t.Run("GetCardInfo", func(t *testing.T) {
		info, err := client.GetCardInfo(ctx, CardInfoRequest{CardNames: []string{"Lightning Bolt"}})
		if err != nil {
			t.Fatalf("GetCardInfo error: %v", err)
		}
		if len(info.Cards) != 1 {
			t.Fatalf("card info count = %d, want 1", len(info.Cards))
		}
	})

	t.Run("SubmitJobApplication", func(t *testing.T) {
		application := []byte("fake zip")
		job, err := client.SubmitJobApplication(ctx, JobApplicationRequest{FirstName: "John", LastName: "Doe", Email: "john@example.com", Application: application, ApplicationFilename: "app.zip"})
		if err != nil {
			t.Fatalf("SubmitJobApplication error: %v", err)
		}
		if !job.Success {
			t.Fatalf("job application success = false, want true")
		}
	})
}
