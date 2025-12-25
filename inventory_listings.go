package manapool

import (
	"context"
	"fmt"
	"net/url"
)

// GetInventoryListings retrieves inventory listings by ID.
func (c *Client) GetInventoryListings(ctx context.Context, ids []string) (*InventoryListingsResponse, error) {
	if len(ids) == 0 {
		return &InventoryListingsResponse{}, nil
	}

	params := url.Values{}
	for _, id := range ids {
		if id != "" {
			params.Add("id", id)
		}
	}

	resp, err := c.doRequest(ctx, "GET", "/inventory/listings", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory listings: %w", err)
	}

	var listings InventoryListingsResponse
	if err := c.decodeResponse(resp, &listings); err != nil {
		return nil, fmt.Errorf("failed to decode inventory listings: %w", err)
	}

	return &listings, nil
}

// GetInventoryListing retrieves a single inventory listing by ID.
func (c *Client) GetInventoryListing(ctx context.Context, id string) (*InventoryItemResponse, error) {
	if id == "" {
		return nil, NewValidationError("id", "id cannot be empty")
	}

	endpoint := fmt.Sprintf("/inventory/listings/%s", id)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory listing: %w", err)
	}

	var listing InventoryItemResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode inventory listing: %w", err)
	}

	return &listing, nil
}

// GetInventoryBySKU retrieves an inventory item by TCGPlayer SKU.
func (c *Client) GetInventoryBySKU(ctx context.Context, sku int) (*InventoryListingResponse, error) {
	if sku <= 0 {
		return nil, NewValidationError("sku", "sku must be positive")
	}

	endpoint := fmt.Sprintf("/inventory/tcgsku/%d", sku)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory by sku: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode inventory by sku: %w", err)
	}

	return &listing, nil
}

// UpdateInventoryBySKU updates an inventory item by TCGPlayer SKU.
func (c *Client) UpdateInventoryBySKU(ctx context.Context, sku int, update InventoryUpdateRequest) (*InventoryListingResponse, error) {
	if sku <= 0 {
		return nil, NewValidationError("sku", "sku must be positive")
	}

	endpoint := fmt.Sprintf("/inventory/tcgsku/%d", sku)
	resp, err := c.doJSONRequest(ctx, "PUT", endpoint, nil, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update inventory by sku: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode inventory update: %w", err)
	}

	return &listing, nil
}

// DeleteInventoryBySKU deletes an inventory item by TCGPlayer SKU.
func (c *Client) DeleteInventoryBySKU(ctx context.Context, sku int) (*InventoryListingResponse, error) {
	if sku <= 0 {
		return nil, NewValidationError("sku", "sku must be positive")
	}

	endpoint := fmt.Sprintf("/inventory/tcgsku/%d", sku)
	resp, err := c.doRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to delete inventory by sku: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode inventory delete: %w", err)
	}

	return &listing, nil
}

// CreateInventoryBulk updates inventory in bulk by SKU.
func (c *Client) CreateInventoryBulk(ctx context.Context, items []InventoryBulkItemBySKU) (*InventoryItemsResponse, error) {
	if len(items) == 0 {
		return nil, NewValidationError("items", "items cannot be empty")
	}

	resp, err := c.doJSONRequest(ctx, "POST", "/seller/inventory", nil, items)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory bulk: %w", err)
	}

	var listing InventoryItemsResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode inventory bulk response: %w", err)
	}

	return &listing, nil
}

// CreateInventoryBulkBySKU updates inventory in bulk by TCGPlayer SKU.
func (c *Client) CreateInventoryBulkBySKU(ctx context.Context, items []InventoryBulkItemBySKU) (*InventoryItemsResponse, error) {
	if len(items) == 0 {
		return nil, NewValidationError("items", "items cannot be empty")
	}

	resp, err := c.doJSONRequest(ctx, "POST", "/seller/inventory/tcgsku", nil, items)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory bulk by sku: %w", err)
	}

	var listing InventoryItemsResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode inventory bulk by sku response: %w", err)
	}

	return &listing, nil
}

// GetSellerInventoryBySKU retrieves a seller inventory item by SKU.
func (c *Client) GetSellerInventoryBySKU(ctx context.Context, sku int) (*InventoryListingResponse, error) {
	if sku <= 0 {
		return nil, NewValidationError("sku", "sku must be positive")
	}

	endpoint := fmt.Sprintf("/seller/inventory/tcgsku/%d", sku)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller inventory by sku: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode seller inventory by sku: %w", err)
	}

	return &listing, nil
}

// UpdateSellerInventoryBySKU updates a seller inventory item by SKU.
func (c *Client) UpdateSellerInventoryBySKU(ctx context.Context, sku int, update InventoryUpdateRequest) (*InventoryListingResponse, error) {
	if sku <= 0 {
		return nil, NewValidationError("sku", "sku must be positive")
	}

	endpoint := fmt.Sprintf("/seller/inventory/tcgsku/%d", sku)
	resp, err := c.doJSONRequest(ctx, "PUT", endpoint, nil, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update seller inventory by sku: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode seller inventory update: %w", err)
	}

	return &listing, nil
}

// DeleteSellerInventoryBySKU deletes a seller inventory item by SKU.
func (c *Client) DeleteSellerInventoryBySKU(ctx context.Context, sku int) (*InventoryListingResponse, error) {
	if sku <= 0 {
		return nil, NewValidationError("sku", "sku must be positive")
	}

	endpoint := fmt.Sprintf("/seller/inventory/tcgsku/%d", sku)
	resp, err := c.doRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to delete seller inventory by sku: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode seller inventory delete: %w", err)
	}

	return &listing, nil
}

// CreateInventoryBulkByProduct upserts inventory in bulk by product.
func (c *Client) CreateInventoryBulkByProduct(ctx context.Context, items []InventoryBulkItemByProduct) (*InventoryItemsResponse, error) {
	if len(items) == 0 {
		return nil, NewValidationError("items", "items cannot be empty")
	}

	resp, err := c.doJSONRequest(ctx, "POST", "/seller/inventory/product", nil, items)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory bulk by product: %w", err)
	}

	var listing InventoryItemsResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode inventory bulk by product response: %w", err)
	}

	return &listing, nil
}

// GetSellerInventoryByProduct retrieves inventory by product ID.
func (c *Client) GetSellerInventoryByProduct(ctx context.Context, productType, productID string) (*InventoryListingResponse, error) {
	if productType == "" || productID == "" {
		return nil, NewValidationError("product", "productType and productID are required")
	}

	endpoint := fmt.Sprintf("/seller/inventory/product/%s/%s", productType, productID)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller inventory by product: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode seller inventory by product: %w", err)
	}

	return &listing, nil
}

// UpdateSellerInventoryByProduct updates inventory by product ID.
func (c *Client) UpdateSellerInventoryByProduct(ctx context.Context, productType, productID string, update InventoryUpdateRequest) (*InventoryListingResponse, error) {
	if productType == "" || productID == "" {
		return nil, NewValidationError("product", "productType and productID are required")
	}

	endpoint := fmt.Sprintf("/seller/inventory/product/%s/%s", productType, productID)
	resp, err := c.doJSONRequest(ctx, "PUT", endpoint, nil, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update seller inventory by product: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode seller inventory update by product: %w", err)
	}

	return &listing, nil
}

// DeleteSellerInventoryByProduct deletes inventory by product ID.
func (c *Client) DeleteSellerInventoryByProduct(ctx context.Context, productType, productID string) (*InventoryListingResponse, error) {
	if productType == "" || productID == "" {
		return nil, NewValidationError("product", "productType and productID are required")
	}

	endpoint := fmt.Sprintf("/seller/inventory/product/%s/%s", productType, productID)
	resp, err := c.doRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to delete seller inventory by product: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode seller inventory delete by product: %w", err)
	}

	return &listing, nil
}

// CreateInventoryBulkByScryfall upserts inventory in bulk by Scryfall ID.
func (c *Client) CreateInventoryBulkByScryfall(ctx context.Context, items []InventoryBulkItemByScryfall) (*InventoryItemsResponse, error) {
	if len(items) == 0 {
		return nil, NewValidationError("items", "items cannot be empty")
	}

	resp, err := c.doJSONRequest(ctx, "POST", "/seller/inventory/scryfall_id", nil, items)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory bulk by scryfall: %w", err)
	}

	var listing InventoryItemsResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode inventory bulk by scryfall response: %w", err)
	}

	return &listing, nil
}

// GetSellerInventoryByScryfall retrieves inventory by Scryfall ID.
func (c *Client) GetSellerInventoryByScryfall(ctx context.Context, scryfallID string, opts InventoryByScryfallOptions) (*InventoryListingResponse, error) {
	if scryfallID == "" {
		return nil, NewValidationError("scryfall_id", "scryfallID cannot be empty")
	}

	params := opts.toParams()

	endpoint := fmt.Sprintf("/seller/inventory/scryfall_id/%s", scryfallID)
	resp, err := c.doRequest(ctx, "GET", endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller inventory by scryfall: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode seller inventory by scryfall: %w", err)
	}

	return &listing, nil
}

// UpdateSellerInventoryByScryfall updates inventory by Scryfall ID.
func (c *Client) UpdateSellerInventoryByScryfall(ctx context.Context, scryfallID string, opts InventoryByScryfallOptions, update InventoryUpdateRequest) (*InventoryListingResponse, error) {
	if scryfallID == "" {
		return nil, NewValidationError("scryfall_id", "scryfallID cannot be empty")
	}

	params := opts.toParams()

	endpoint := fmt.Sprintf("/seller/inventory/scryfall_id/%s", scryfallID)
	resp, err := c.doJSONRequest(ctx, "PUT", endpoint, params, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update seller inventory by scryfall: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode seller inventory update by scryfall: %w", err)
	}

	return &listing, nil
}

// DeleteSellerInventoryByScryfall deletes inventory by Scryfall ID.
func (c *Client) DeleteSellerInventoryByScryfall(ctx context.Context, scryfallID string, opts InventoryByScryfallOptions) (*InventoryListingResponse, error) {
	if scryfallID == "" {
		return nil, NewValidationError("scryfall_id", "scryfallID cannot be empty")
	}

	params := opts.toParams()

	endpoint := fmt.Sprintf("/seller/inventory/scryfall_id/%s", scryfallID)
	resp, err := c.doRequest(ctx, "DELETE", endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("failed to delete seller inventory by scryfall: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode seller inventory delete by scryfall: %w", err)
	}

	return &listing, nil
}

// CreateInventoryBulkByTCGPlayerID upserts inventory in bulk by TCGPlayer ID.
func (c *Client) CreateInventoryBulkByTCGPlayerID(ctx context.Context, items []InventoryBulkItemByTCGPlayerID) (*InventoryItemsResponse, error) {
	if len(items) == 0 {
		return nil, NewValidationError("items", "items cannot be empty")
	}

	resp, err := c.doJSONRequest(ctx, "POST", "/seller/inventory/tcgplayer_id", nil, items)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory bulk by tcgplayer: %w", err)
	}

	var listing InventoryItemsResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode inventory bulk by tcgplayer response: %w", err)
	}

	return &listing, nil
}

// GetSellerInventoryByTCGPlayerID retrieves inventory by TCGPlayer ID.
func (c *Client) GetSellerInventoryByTCGPlayerID(ctx context.Context, tcgplayerID int, opts InventoryByTCGPlayerOptions) (*InventoryListingResponse, error) {
	if tcgplayerID <= 0 {
		return nil, NewValidationError("tcgplayer_id", "tcgplayerID must be positive")
	}

	params := opts.toParams()

	endpoint := fmt.Sprintf("/seller/inventory/tcgplayer_id/%d", tcgplayerID)
	resp, err := c.doRequest(ctx, "GET", endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller inventory by tcgplayer: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode seller inventory by tcgplayer: %w", err)
	}

	return &listing, nil
}

// UpdateSellerInventoryByTCGPlayerID updates inventory by TCGPlayer ID.
func (c *Client) UpdateSellerInventoryByTCGPlayerID(ctx context.Context, tcgplayerID int, opts InventoryByTCGPlayerOptions, update InventoryUpdateRequest) (*InventoryListingResponse, error) {
	if tcgplayerID <= 0 {
		return nil, NewValidationError("tcgplayer_id", "tcgplayerID must be positive")
	}

	params := opts.toParams()

	endpoint := fmt.Sprintf("/seller/inventory/tcgplayer_id/%d", tcgplayerID)
	resp, err := c.doJSONRequest(ctx, "PUT", endpoint, params, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update seller inventory by tcgplayer: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode seller inventory update by tcgplayer: %w", err)
	}

	return &listing, nil
}

// DeleteSellerInventoryByTCGPlayerID deletes inventory by TCGPlayer ID.
func (c *Client) DeleteSellerInventoryByTCGPlayerID(ctx context.Context, tcgplayerID int, opts InventoryByTCGPlayerOptions) (*InventoryListingResponse, error) {
	if tcgplayerID <= 0 {
		return nil, NewValidationError("tcgplayer_id", "tcgplayerID must be positive")
	}

	params := opts.toParams()

	endpoint := fmt.Sprintf("/seller/inventory/tcgplayer_id/%d", tcgplayerID)
	resp, err := c.doRequest(ctx, "DELETE", endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("failed to delete seller inventory by tcgplayer: %w", err)
	}

	var listing InventoryListingResponse
	if err := c.decodeResponse(resp, &listing); err != nil {
		return nil, fmt.Errorf("failed to decode seller inventory delete by tcgplayer: %w", err)
	}

	return &listing, nil
}
