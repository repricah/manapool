package manapool

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// GetSellerInventory retrieves the seller's inventory with pagination support.
//
// The inventory includes all products (singles and sealed) with their current
// prices, quantities, and product details. Results are paginated.
//
// Example:
//
//	opts := manapool.InventoryOptions{
//	    Limit:  500,  // max 500
//	    Offset: 0,    // start at beginning
//	}
//	resp, err := client.GetSellerInventory(ctx, opts)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Total items: %d, Returned: %d\n",
//	    resp.Pagination.Total, resp.Pagination.Returned)
//	for _, item := range resp.Inventory {
//	    fmt.Printf("  %s: $%.2f (qty: %d)\n",
//	        item.Product.Single.Name, item.PriceDollars(), item.Quantity)
//	}
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - opts: Pagination options (limit and offset)
//
// Returns:
//   - *InventoryResponse: The inventory items and pagination metadata
//   - error: Any error that occurred during the request
func (c *Client) GetSellerInventory(ctx context.Context, opts InventoryOptions) (*InventoryResponse, error) {
	// Validate options
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	c.logger.Debugf("Getting seller inventory: limit=%d, offset=%d", opts.Limit, opts.Offset)

	// Build query parameters
	params := url.Values{}
	params.Add("limit", strconv.Itoa(opts.Limit))
	params.Add("offset", strconv.Itoa(opts.Offset))

	resp, err := c.doRequest(ctx, "GET", "/seller/inventory", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller inventory: %w", err)
	}

	var inventoryResp InventoryResponse
	if err := c.decodeResponse(resp, &inventoryResp); err != nil {
		return nil, fmt.Errorf("failed to decode seller inventory: %w", err)
	}

	c.logger.Debugf("Retrieved %d inventory items (total: %d)",
		inventoryResp.Pagination.Returned, inventoryResp.Pagination.Total)

	return &inventoryResp, nil
}

// GetInventoryByTCGPlayerID retrieves a specific inventory item by its TCGPlayer SKU.
//
// This is useful when you need to look up a specific card by its TCGPlayer ID
// to check its current Manapool price and quantity.
//
// Example:
//
//	item, err := client.GetInventoryByTCGPlayerID(ctx, "4549403")
//	if err != nil {
//	    var apiErr *manapool.APIError
//	    if errors.As(err, &apiErr) && apiErr.IsNotFound() {
//	        fmt.Println("Item not found in inventory")
//	        return
//	    }
//	    log.Fatal(err)
//	}
//	fmt.Printf("Found: %s - $%.2f (qty: %d)\n",
//	    item.Product.Single.Name, item.PriceDollars(), item.Quantity)
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - tcgplayerID: The TCGPlayer SKU to look up
//
// Returns:
//   - *InventoryItem: The inventory item
//   - error: Any error that occurred during the request (404 if not found)
func (c *Client) GetInventoryByTCGPlayerID(ctx context.Context, tcgplayerID string) (*InventoryItem, error) {
	if tcgplayerID == "" {
		return nil, NewValidationError("tcgplayerID", "tcgplayerID cannot be empty")
	}

	c.logger.Debugf("Getting inventory by TCGPlayer ID: %s", tcgplayerID)

	endpoint := fmt.Sprintf("/seller/inventory/tcgsku/%s", tcgplayerID)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory by TCGPlayer ID: %w", err)
	}

	var item InventoryItem
	if err := c.decodeResponse(resp, &item); err != nil {
		return nil, fmt.Errorf("failed to decode inventory item: %w", err)
	}

	c.logger.Debugf("Retrieved inventory item: %s (TCG SKU: %d)",
		item.Product.Single.Name, item.Product.TCGPlayerSKU)

	return &item, nil
}

// IterateInventory is a helper function that automatically handles pagination
// and calls the provided callback for each inventory item.
//
// This is useful when you need to process all inventory items without manually
// managing pagination. The iteration continues until all items are processed
// or an error occurs.
//
// Example:
//
//	err := manapool.IterateInventory(ctx, client, func(item *manapool.InventoryItem) error {
//	    fmt.Printf("%s: $%.2f (qty: %d)\n",
//	        item.Product.Single.Name, item.PriceDollars(), item.Quantity)
//	    return nil
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - client: The Manapool API client
//   - callback: Function called for each inventory item
//
// Returns:
//   - error: Any error that occurred during iteration
func IterateInventory(ctx context.Context, client APIClient, callback func(*InventoryItem) error) error {
	offset := 0
	limit := 500

	for {
		opts := InventoryOptions{
			Limit:  limit,
			Offset: offset,
		}

		resp, err := client.GetSellerInventory(ctx, opts)
		if err != nil {
			return fmt.Errorf("failed to get inventory at offset %d: %w", offset, err)
		}

		// Process items
		for i := range resp.Inventory {
			if err := callback(&resp.Inventory[i]); err != nil {
				return fmt.Errorf("callback error at offset %d: %w", offset, err)
			}
		}

		// Check if we're done
		if resp.Pagination.Returned == 0 || offset+resp.Pagination.Returned >= resp.Pagination.Total {
			break
		}

		offset += resp.Pagination.Returned
	}

	return nil
}
