package manapool

import (
	"context"
	"fmt"
)

// GetSellerAccount retrieves the authenticated seller's account information.
//
// This endpoint returns account details including username, email, verification status,
// and whether singles/sealed products are live on the marketplace.
//
// Example:
//
//	account, err := client.GetSellerAccount(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Account: %s (%s)\n", account.Username, account.Email)
//	fmt.Printf("Singles Live: %v, Sealed Live: %v\n",
//	    account.SinglesLive, account.SealedLive)
//
// Returns:
//   - *Account: The account information
//   - error: Any error that occurred during the request
func (c *Client) GetSellerAccount(ctx context.Context) (*Account, error) {
	c.logger.Debugf("Getting seller account")

	resp, err := c.doRequest(ctx, "GET", "/account", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get seller account: %w", err)
	}

	var account Account
	if err := c.decodeResponse(resp, &account); err != nil {
		return nil, fmt.Errorf("failed to decode seller account: %w", err)
	}

	c.logger.Debugf("Retrieved seller account: %s (%s)", account.Username, account.Email)

	return &account, nil
}
