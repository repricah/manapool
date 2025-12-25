# Manapool Go Client

> **⚠️ Pre-Release (v0.2.0)**: This library is under active development. The API may change before v1.0.0. Planned features are tracked in GitHub issues GitHub issues.

A Go client library for the [Manapool API](https://manapool.com). This library provides a Go interface for managing your Magic: The Gathering inventory on Manapool.

[![Go Reference](https://pkg.go.dev/badge/github.com/repricah/manapool.svg)](https://pkg.go.dev/github.com/repricah/manapool)
[![Go Report Card](https://goreportcard.com/badge/github.com/repricah/manapool)](https://goreportcard.com/report/github.com/repricah/manapool)

## Features

### Currently Implemented (v0.2.0)

- ✅ **Seller Inventory Endpoints** - Get account, list inventory, lookup by TCG SKU
- ✅ **Type-Safe** - Full Go type definitions for all API models
- ✅ **Automatic Rate Limiting** - Built-in rate limiter to respect API limits
- ✅ **Automatic Retries** - Configurable retry logic with exponential backoff
- ✅ **Context Support** - First-class context support for cancellation and timeouts
- ✅ **Error Handling** - Specific error types with helper methods
- ✅ **Tested** - 96.5% test coverage with integration tests
- ✅ **Production Use** - Used in production for TCG inventory management
- ✅ **Zero Dependencies** - Only depends on `golang.org/x/time/rate`

### Planned Features

Planned features are tracked in GitHub:
- 🔜 **Additional Lookups** () - Scryfall ID, product ID, TCGPlayer ID lookups
- 🔜 **Order Management** () - List and view seller orders
- 🔜 **Order Fulfillment** () - Mark orders as shipped/fulfilled
- 🔜 **Inventory Updates** () - Create, update, delete inventory items
- 🔜 **Webhook Support** () - Register and manage webhooks
- 🔜 **Release Readiness** () - v1.0.0 stabilization and publishing steps
- 🔜 **Repository Extraction** () - Move the client into a standalone repository

## Installation

```bash
go get github.com/repricah/manapool
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/repricah/manapool"
)

func main() {
    // Create client with your API credentials
    client := manapool.NewClient(
        "your-api-token",
        "your-email@example.com",
        manapool.WithTimeout(30*time.Second),
        manapool.WithRateLimit(10, 1), // 10 requests/second, burst of 1
    )

    ctx := context.Background()

    // Get your seller account information
    account, err := client.GetSellerAccount(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Account: %s (%s)\n", account.Username, account.Email)
    fmt.Printf("Singles Live: %v, Sealed Live: %v\n",
        account.SinglesLive, account.SealedLive)

    // Get your inventory with pagination
    opts := manapool.InventoryOptions{
        Limit:  500, // max 500 items per request
        Offset: 0,
    }
    inventory, err := client.GetSellerInventory(ctx, opts)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Total inventory: %d items\n", inventory.Pagination.Total)

    // Print first few items
    for _, item := range inventory.Inventory {
        fmt.Printf("  %s - $%.2f (qty: %d)\n",
            item.Product.Single.Name,
            item.PriceDollars(),
            item.Quantity)
    }
}
```

## Beads Integration

If you're using this client in Beads applications, follow the same patterns you would in Go services: configure credentials via environment variables, propagate context for cancellation, respect rate limits, and surface structured errors to the UI or logs.

```go
client := manapool.NewClient(
    os.Getenv("MANAPOOL_TOKEN"),
    os.Getenv("MANAPOOL_EMAIL"),
    manapool.WithTimeout(30*time.Second),
    manapool.WithRateLimit(10, 1),
)

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

opts := manapool.InventoryOptions{
    Limit:  100,
    Offset: 0,
}
inventory, err := client.GetSellerInventory(ctx, opts)
if err != nil {
    var apiErr *manapool.APIError
    if errors.As(err, &apiErr) {
        // Display a user-friendly message in Beads UI.
        return
    }
    // Handle network/unknown errors.
    return
}
```

Beads resources:
- Official Site: https://www.beadsproject.net/
- Documentation: https://www.beadsproject.net/ref/
- Examples: https://github.com/magicmouse/beads-examples

## Usage Examples

### Authentication

The Manapool API uses token-based authentication with two required headers:

```go
client := manapool.NewClient("your-api-token", "your-email@example.com")
```

### Get Seller Account

```go
account, err := client.GetSellerAccount(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Account: %s\n", account.Username)
fmt.Printf("Verified: %v\n", account.Verified)
fmt.Printf("Payouts Enabled: %v\n", account.PayoutsEnabled)
```

### Get Inventory with Pagination

```go
opts := manapool.InventoryOptions{
    Limit:  500,
    Offset: 0,
}

inventory, err := client.GetSellerInventory(ctx, opts)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total: %d, Returned: %d\n",
    inventory.Pagination.Total,
    inventory.Pagination.Returned)
```

### Iterate All Inventory

Use the helper function to automatically handle pagination:

```go
err := manapool.IterateInventory(ctx, client, func(item *manapool.InventoryItem) error {
    fmt.Printf("%s: $%.2f (TCG SKU: %d)\n",
        item.Product.Single.Name,
        item.PriceDollars(),
        item.Product.TCGPlayerSKU)
    return nil
})
if err != nil {
    log.Fatal(err)
}
```

### Look Up Item by TCGPlayer SKU

```go
item, err := client.GetInventoryByTCGPlayerID(ctx, "4549403")
if err != nil {
    var apiErr *manapool.APIError
    if errors.As(err, &apiErr) && apiErr.IsNotFound() {
        fmt.Println("Item not found in inventory")
        return
    }
    log.Fatal(err)
}

fmt.Printf("Found: %s\n", item.Product.Single.Name)
fmt.Printf("Price: $%.2f\n", item.PriceDollars())
fmt.Printf("Quantity: %d\n", item.Quantity)
fmt.Printf("Condition: %s\n", item.Product.Single.ConditionName())
```

## Configuration Options

### Custom HTTP Client

```go
customClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:    10,
        IdleConnTimeout: 90 * time.Second,
    },
}

client := manapool.NewClient(token, email,
    manapool.WithHTTPClient(customClient),
)
```

### Rate Limiting

```go
client := manapool.NewClient(token, email,
    manapool.WithRateLimit(5, 2), // 5 requests/second, burst of 2
)
```

### Retry Configuration

```go
client := manapool.NewClient(token, email,
    manapool.WithRetry(5, 2*time.Second), // 5 retries, 2s initial backoff
)
```

### Custom Logger

```go
type myLogger struct{}

func (l *myLogger) Debugf(format string, args ...interface{}) {
    log.Printf("[DEBUG] "+format, args...)
}

func (l *myLogger) Errorf(format string, args ...interface{}) {
    log.Printf("[ERROR] "+format, args...)
}

client := manapool.NewClient(token, email,
    manapool.WithLogger(&myLogger{}),
)
```

### All Options Together

```go
client := manapool.NewClient(token, email,
    manapool.WithHTTPClient(customHTTP),
    manapool.WithBaseURL("https://custom.api.com/v1/"),
    manapool.WithRateLimit(10, 1),
    manapool.WithRetry(3, time.Second),
    manapool.WithTimeout(30*time.Second),
    manapool.WithUserAgent("my-app/1.0"),
    manapool.WithLogger(logger),
)
```

## Error Handling

The library provides specific error types for different scenarios:

### API Errors

```go
inventory, err := client.GetSellerInventory(ctx, opts)
if err != nil {
    var apiErr *manapool.APIError
    if errors.As(err, &apiErr) {
        switch {
        case apiErr.IsNotFound():
            fmt.Println("Resource not found")
        case apiErr.IsUnauthorized():
            fmt.Println("Invalid credentials")
        case apiErr.IsForbidden():
            fmt.Println("Access denied")
        case apiErr.IsRateLimited():
            fmt.Println("Rate limit exceeded")
        case apiErr.IsServerError():
            fmt.Println("Server error, retry later")
        default:
            fmt.Printf("API error: %v\n", apiErr)
        }
        return
    }
    log.Fatal(err)
}
```

### Validation Errors

```go
item, err := client.GetInventoryByTCGPlayerID(ctx, "")
if err != nil {
    var valErr *manapool.ValidationError
    if errors.As(err, &valErr) {
        fmt.Printf("Validation error for %s: %s\n", valErr.Field, valErr.Message)
        return
    }
}
```

### Network Errors

```go
inventory, err := client.GetSellerInventory(ctx, opts)
if err != nil {
    var netErr *manapool.NetworkError
    if errors.As(err, &netErr) {
        fmt.Printf("Network error: %v\n", netErr)
        // Maybe retry with exponential backoff
        return
    }
}
```

## Type Definitions

### Account

```go
type Account struct {
    Username       string
    Email          string
    Verified       bool
    SinglesLive    bool
    SealedLive     bool
    PayoutsEnabled bool
}
```

### InventoryItem

```go
type InventoryItem struct {
    ID            string
    ProductType   string
    ProductID     string
    Product       Product
    PriceCents    int      // Price in cents
    Quantity      int
    EffectiveAsOf Timestamp
}

// Helper methods
func (i InventoryItem) PriceDollars() float64 // Convert cents to dollars
```

### Product

```go
type Product struct {
    Type         string
    ID           string
    TCGPlayerSKU int
    Single       Single
    Sealed       Sealed
}
```

### Single (Card)

```go
type Single struct {
    ScryfallID  string
    MTGJsonID   string
    Name        string
    Set         string
    Number      string
    LanguageID  string
    ConditionID string // NM, LP, MP, HP, DMG
    FinishID    string // NF (non-foil), FO (foil), EF (etched foil)
}

// Helper methods
func (s Single) ConditionName() string // Returns "Near Mint", "Near Mint Foil", etc.
```

## Testing

The library includes comprehensive tests with 96.5% coverage:

```bash
go test -v -cover ./...
```

Run tests with race detector:

```bash
go test -race -v ./...
```

Generate coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/new-feature`)
3. Write tests for your changes
4. Ensure tests pass and coverage remains high
5. Commit your changes (`git commit -m 'Add new feature'`)
6. Push to the branch (`git push origin feature/new-feature`)
7. Open a Pull Request

## License

MIT License - see LICENSE file for details

## Support

- 📖 [API Documentation](https://pkg.go.dev/github.com/repricah/manapool)
- 🐛 [Issue Tracker](https://github.com/repricah/manapool/issues)
- 💬 [Discussions](https://github.com/repricah/manapool/discussions)

## Changelog

### v0.2.0 (2025-12-23)

- 🔄 Rename module to `github.com/repricah/manapool`
- ⚖️ Use neutral tone in documentation
- 👤 Corrected authorship to `jblotus`

### v0.2.0 (2025-01-28)

- 🎉 Initial pre-release
- ✅ Seller account endpoint (`GetSellerAccount`)
- ✅ Seller inventory endpoints (`GetSellerInventory`, `GetInventoryByTCGPlayerID`, `IterateInventory`)
- ✅ Test coverage (96.5%)
- ✅ Rate limiting and retries
- ✅ Context support for all operations
- ✅ Structured error handling with helper methods
- ✅ Configurable client options
- ⚠️ API may change before v1.0.0
