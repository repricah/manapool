// Package manapool provides a client library for the Manapool API.
//
// The Manapool API allows sellers to manage their inventory, view account information,
// and interact with their Manapool store programmatically.
//
// # Basic Usage
//
//	client := manapool.NewClient("your-api-token", "your-email@example.com")
//	account, err := client.GetSellerAccount(context.Background())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Account: %s (%s)\n", account.Username, account.Email)
//
// # Authentication
//
// The Manapool API uses token-based authentication with two required headers:
//   - X-ManaPool-Access-Token: Your API access token
//   - X-ManaPool-Email: Your account email address
//
// # Rate Limiting
//
// The client includes built-in rate limiting to avoid overwhelming the API.
// Configure rate limits using WithRateLimit option:
//
//	client := manapool.NewClient(token, email,
//	    manapool.WithRateLimit(10, 1), // 10 requests per second, burst of 1
//	)
//
// # Error Handling
//
// All methods return errors that can be inspected for API-specific details:
//
//	inventory, err := client.GetSellerInventory(ctx, opts)
//	if err != nil {
//	    if apiErr, ok := err.(*manapool.APIError); ok {
//	        fmt.Printf("API Error: %d - %s\n", apiErr.StatusCode, apiErr.Message)
//	    }
//	    return err
//	}
package manapool

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Account represents a Manapool seller account.
type Account struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Verified       bool   `json:"verified"`
	SinglesLive    bool   `json:"singles_live"`
	SealedLive     bool   `json:"sealed_live"`
	PayoutsEnabled bool   `json:"payouts_enabled"`
}

// InventoryResponse represents a paginated response from the inventory API.
type InventoryResponse struct {
	Inventory  []InventoryItem `json:"inventory"`
	Pagination Pagination      `json:"pagination"`
}

// InventoryItem represents a single inventory item in the Manapool system.
type InventoryItem struct {
	ID            string    `json:"id"`
	ProductType   string    `json:"product_type"`
	ProductID     string    `json:"product_id"`
	Product       Product   `json:"product"`
	PriceCents    int       `json:"price_cents"`
	Quantity      int       `json:"quantity"`
	EffectiveAsOf Timestamp `json:"effective_as_of"`
}

// Product represents a product in the Manapool inventory.
type Product struct {
	Type         string `json:"type"`
	ID           string `json:"id"`
	TCGPlayerSKU int    `json:"tcgplayer_sku"`
	Single       Single `json:"single"`
	Sealed       Sealed `json:"sealed"`
}

// Single represents a single card product.
type Single struct {
	ScryfallID  string `json:"scryfall_id"`
	MTGJsonID   string `json:"mtgjson_id"`
	Name        string `json:"name"`
	Set         string `json:"set"`
	Number      string `json:"number"`
	LanguageID  string `json:"language_id"`
	ConditionID string `json:"condition_id"`
	FinishID    string `json:"finish_id"`
}

// Sealed represents a sealed product (booster boxes, etc.).
type Sealed struct {
	MTGJsonID  string `json:"mtgjson_id"`
	Name       string `json:"name"`
	Set        string `json:"set"`
	LanguageID string `json:"language_id"`
}

// Pagination contains pagination metadata for API responses.
type Pagination struct {
	Total    int `json:"total"`
	Returned int `json:"returned"`
	Offset   int `json:"offset"`
	Limit    int `json:"limit"`
}

// Timestamp is a custom time type that handles Manapool's timestamp format.
// The Manapool API returns timestamps in multiple formats:
//   - RFC3339Nano: "2025-08-05T20:38:54.549229Z"
//   - No-colon offset: "2025-08-05T20:38:54.549229+0000"
type Timestamp struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler for Timestamp.
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`) // strip quotes

	// Try standard RFC3339Nano first
	if parsed, err := time.Parse(time.RFC3339Nano, s); err == nil {
		t.Time = parsed
		return nil
	}

	// Fallback: no-colon offset like +0000
	if parsed, err := time.Parse("2006-01-02T15:04:05.999999-0700", s); err == nil {
		t.Time = parsed
		return nil
	}

	return fmt.Errorf("cannot parse timestamp: %q", s)
}

// MarshalJSON implements json.Marshaler for Timestamp.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.Format(time.RFC3339Nano))
}

// InventoryOptions contains options for querying seller inventory.
type InventoryOptions struct {
	// Limit specifies the maximum number of items to return (default: 500, max: 500)
	Limit int

	// Offset specifies the starting position in the result set (default: 0)
	Offset int
}

// Validate validates the inventory options and sets defaults.
func (o *InventoryOptions) Validate() error {
	if o.Limit < 0 {
		return fmt.Errorf("limit must be non-negative, got %d", o.Limit)
	}
	if o.Limit > 500 {
		return fmt.Errorf("limit must not exceed 500, got %d", o.Limit)
	}
	if o.Limit == 0 {
		o.Limit = 500 // default
	}

	if o.Offset < 0 {
		return fmt.Errorf("offset must be non-negative, got %d", o.Offset)
	}

	return nil
}

// ConditionName returns the standardized condition name for a single card.
// It combines the condition ID and finish ID into a human-readable string.
//
// Condition IDs:
//   - NM: Near Mint
//   - LP: Lightly Played
//   - MP: Moderately Played
//   - HP: Heavily Played
//   - DMG: Damaged
//
// Finish IDs:
//   - NF: Non-foil (no suffix)
//   - FO: Foil (adds " Foil" suffix)
//   - EF: Etched Foil (adds " Foil" suffix)
func (s Single) ConditionName() string {
	var condition string

	switch s.ConditionID {
	case "NM":
		condition = "Near Mint"
	case "LP":
		condition = "Lightly Played"
	case "MP":
		condition = "Moderately Played"
	case "HP":
		condition = "Heavily Played"
	case "DMG":
		condition = "Damaged"
	default:
		condition = "Unknown"
	}

	// Add foil suffix for foil finishes
	switch s.FinishID {
	case "FO", "EF":
		condition += " Foil"
	}

	return condition
}

// PriceDollars returns the price in dollars (converts from cents).
func (i InventoryItem) PriceDollars() float64 {
	return float64(i.PriceCents) / 100.0
}
