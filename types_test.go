package manapool

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "RFC3339Nano format",
			input:   `"2025-08-05T20:38:54.549229Z"`,
			want:    time.Date(2025, 8, 5, 20, 38, 54, 549229000, time.UTC),
			wantErr: false,
		},
		{
			name:    "no-colon offset format",
			input:   `"2025-08-05T20:38:54.549229+0000"`,
			want:    time.Date(2025, 8, 5, 20, 38, 54, 549229000, time.UTC),
			wantErr: false,
		},
		{
			name:    "no-colon offset with negative timezone",
			input:   `"2025-08-05T20:38:54.549229-0500"`,
			want:    time.Date(2025, 8, 5, 20, 38, 54, 549229000, time.FixedZone("", -5*3600)),
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   `"not-a-timestamp"`,
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   `""`,
			want:    time.Time{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts Timestamp
			err := json.Unmarshal([]byte(tt.input), &ts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Timestamp.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !ts.Equal(tt.want) {
				t.Errorf("Timestamp.UnmarshalJSON() = %v, want %v", ts, tt.want)
			}
		})
	}
}

func TestTimestamp_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		time    time.Time
		want    string
		wantErr bool
	}{
		{
			name:    "valid time",
			time:    time.Date(2025, 8, 5, 20, 38, 54, 549229000, time.UTC),
			want:    `"2025-08-05T20:38:54.549229Z"`,
			wantErr: false,
		},
		{
			name:    "zero time",
			time:    time.Time{},
			want:    `null`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := Timestamp{Time: tt.time}
			got, err := json.Marshal(ts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Timestamp.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if string(got) != tt.want {
				t.Errorf("Timestamp.MarshalJSON() = %s, want %s", string(got), tt.want)
			}
		})
	}
}

func TestInventoryOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    InventoryOptions
		wantErr bool
		wantMsg string
	}{
		{
			name:    "valid options",
			opts:    InventoryOptions{Limit: 100, Offset: 0},
			wantErr: false,
		},
		{
			name:    "default limit (0 becomes 500)",
			opts:    InventoryOptions{Limit: 0, Offset: 0},
			wantErr: false,
		},
		{
			name:    "max limit (500)",
			opts:    InventoryOptions{Limit: 500, Offset: 0},
			wantErr: false,
		},
		{
			name:    "negative limit",
			opts:    InventoryOptions{Limit: -1, Offset: 0},
			wantErr: true,
			wantMsg: "limit must be non-negative",
		},
		{
			name:    "limit exceeds max",
			opts:    InventoryOptions{Limit: 501, Offset: 0},
			wantErr: true,
			wantMsg: "limit must not exceed 500",
		},
		{
			name:    "negative offset",
			opts:    InventoryOptions{Limit: 100, Offset: -1},
			wantErr: true,
			wantMsg: "offset must be non-negative",
		},
		{
			name:    "large offset",
			opts:    InventoryOptions{Limit: 100, Offset: 10000},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("InventoryOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.wantMsg != "" {
				if !contains(err.Error(), tt.wantMsg) {
					t.Errorf("InventoryOptions.Validate() error = %v, want message containing %q", err, tt.wantMsg)
				}
			}

			// Check that default limit is set
			if !tt.wantErr && tt.opts.Limit == 0 {
				if tt.opts.Limit != 500 {
					t.Errorf("InventoryOptions.Validate() did not set default limit, got %d", tt.opts.Limit)
				}
			}
		})
	}
}

func TestSingle_ConditionName(t *testing.T) {
	tests := []struct {
		name        string
		conditionID string
		finishID    string
		want        string
	}{
		{
			name:        "Near Mint non-foil",
			conditionID: "NM",
			finishID:    "NF",
			want:        "Near Mint",
		},
		{
			name:        "Near Mint foil",
			conditionID: "NM",
			finishID:    "FO",
			want:        "Near Mint Foil",
		},
		{
			name:        "Lightly Played non-foil",
			conditionID: "LP",
			finishID:    "NF",
			want:        "Lightly Played",
		},
		{
			name:        "Lightly Played foil",
			conditionID: "LP",
			finishID:    "FO",
			want:        "Lightly Played Foil",
		},
		{
			name:        "Moderately Played etched foil",
			conditionID: "MP",
			finishID:    "EF",
			want:        "Moderately Played Foil",
		},
		{
			name:        "Heavily Played non-foil",
			conditionID: "HP",
			finishID:    "NF",
			want:        "Heavily Played",
		},
		{
			name:        "Damaged non-foil",
			conditionID: "DMG",
			finishID:    "NF",
			want:        "Damaged",
		},
		{
			name:        "unknown condition",
			conditionID: "UNKNOWN",
			finishID:    "NF",
			want:        "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Single{
				ConditionID: tt.conditionID,
				FinishID:    tt.finishID,
			}

			got := s.ConditionName()
			if got != tt.want {
				t.Errorf("Single.ConditionName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInventoryItem_PriceDollars(t *testing.T) {
	tests := []struct {
		name       string
		priceCents int
		want       float64
	}{
		{
			name:       "zero price",
			priceCents: 0,
			want:       0.0,
		},
		{
			name:       "one dollar",
			priceCents: 100,
			want:       1.0,
		},
		{
			name:       "fractional price",
			priceCents: 499,
			want:       4.99,
		},
		{
			name:       "large price",
			priceCents: 123456,
			want:       1234.56,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := InventoryItem{
				PriceCents: tt.priceCents,
			}

			got := item.PriceDollars()
			if got != tt.want {
				t.Errorf("InventoryItem.PriceDollars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_JSON(t *testing.T) {
	// Test JSON marshaling/unmarshaling
	account := Account{
		Username:       "testuser",
		Email:          "test@example.com",
		Verified:       true,
		SinglesLive:    true,
		SealedLive:     false,
		PayoutsEnabled: true,
	}

	// Marshal
	data, err := json.Marshal(account)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	// Unmarshal
	var decoded Account
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Compare
	if decoded != account {
		t.Errorf("JSON round-trip failed: got %+v, want %+v", decoded, account)
	}
}

func TestInventoryResponse_JSON(t *testing.T) {
	// Test JSON unmarshaling with real API response format
	jsonData := `{
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
					"sealed": {
						"mtgjson_id": "",
						"name": "",
						"set": "",
						"language_id": ""
					}
				}
			}
		],
		"pagination": {
			"total": 1,
			"returned": 1,
			"offset": 0,
			"limit": 500
		}
	}`

	var resp InventoryResponse
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Verify data
	if len(resp.Inventory) != 1 {
		t.Errorf("expected 1 inventory item, got %d", len(resp.Inventory))
	}

	item := resp.Inventory[0]
	if item.ID != "inv123" {
		t.Errorf("item.ID = %q, want %q", item.ID, "inv123")
	}
	if item.PriceCents != 499 {
		t.Errorf("item.PriceCents = %d, want %d", item.PriceCents, 499)
	}
	if item.Quantity != 5 {
		t.Errorf("item.Quantity = %d, want %d", item.Quantity, 5)
	}

	if resp.Pagination.Total != 1 {
		t.Errorf("pagination.Total = %d, want %d", resp.Pagination.Total, 1)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
