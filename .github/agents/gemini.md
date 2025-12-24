# Manapool Go Client - Gemini Agent Instructions

## Quick Reference

For comprehensive project documentation, see the main instructions file:
**→ `.github/copilot-instructions.md`**

This file contains Gemini-specific guidance. All general coding standards, architecture decisions, and project conventions are defined in the copilot instructions file above.

## Gemini-Specific Context

### Project Type
This is a **Go library** (not an application) for the Manapool API, providing type-safe Magic: The Gathering inventory management.

### Key Project Facts
- **Module**: `github.com/repricah/manapool`
- **Go Version**: 1.24.7
- **Test Coverage**: 96.5%
- **Zero external dependencies** (except `golang.org/x/time`)
- **Production-ready** with rate limiting, retries, and comprehensive error handling

## Build & Test Commands

```bash
# Run all tests
go test -v ./...

# Run tests with race detector
go test -race -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Format code (automatically done by go tooling)
gofmt -w .

# Check module
go mod verify
go mod tidy
```

## Architecture Patterns

### Client Pattern
- Single `Client` struct handles all API operations
- Functional options pattern for configuration (`WithTimeout`, `WithRateLimit`, etc.)
- Thread-safe - safe for concurrent use by multiple goroutines

### Error Handling
- Custom error types: `APIError`, `ValidationError`, `NetworkError`
- Error type helper methods: `IsNotFound()`, `IsUnauthorized()`, etc.
- Always wrap errors with context

### Testing Strategy
- Table-driven tests for multiple scenarios
- Use `httptest` for mocking HTTP responses
- Test both success and error paths thoroughly
- Subtests with `t.Run()` for organization

## Code Style Guidelines

### Do This ✅
```go
// Accept context as first parameter
func (c *Client) GetData(ctx context.Context, opts Options) (*Result, error) {
    // Validate inputs
    if err := opts.Validate(); err != nil {
        return nil, err
    }
    
    // Use rate limiter
    if err := c.rateLimiter.Wait(ctx); err != nil {
        return nil, err
    }
    
    // Make request
    result, err := c.doRequest(ctx, "GET", "/endpoint", opts)
    if err != nil {
        return nil, fmt.Errorf("failed to get data: %w", err)
    }
    
    return result, nil
}
```

### Don't Do This ❌
```go
// DON'T: Missing context parameter
func (c *Client) GetData(opts Options) (*Result, error) { }

// DON'T: Storing context in struct
type Client struct {
    ctx context.Context // WRONG
}

// DON'T: Panicking in library code
func Something() {
    panic("error") // WRONG - return error instead
}

// DON'T: Not wrapping errors
return err // WRONG - use fmt.Errorf("context: %w", err)
```

## Preferred Libraries & Tools

### Standard Library First
- Use `net/http` for HTTP operations
- Use `encoding/json` for JSON
- Use `context` for cancellation
- Use `time` for time operations

### Approved External Dependencies
- `golang.org/x/time/rate` - For rate limiting

### Testing
- `testing` (standard library)
- `httptest` (standard library)

## Common Tasks

### Adding a New Endpoint

1. Define request/response types in `types.go`
2. Add validation methods if needed
3. Implement client method in appropriate file
4. Add comprehensive tests
5. Update README.md with examples
6. Maintain >95% test coverage

### Handling API Errors

```go
resp, err := client.DoSomething(ctx, opts)
if err != nil {
    var apiErr *manapool.APIError
    if errors.As(err, &apiErr) {
        switch {
        case apiErr.IsNotFound():
            // Handle 404
        case apiErr.IsUnauthorized():
            // Handle 401
        case apiErr.IsRateLimited():
            // Handle 429
        default:
            // Handle other errors
        }
    }
    return err
}
```

## Things to Avoid

- ❌ Adding unnecessary external dependencies
- ❌ Breaking backwards compatibility (this is pre-v1.0.0, but still be careful)
- ❌ Reducing test coverage below 95%
- ❌ Making real API calls in tests
- ❌ Storing credentials in code
- ❌ Using `panic()` in library code
- ❌ Ignoring context cancellation

## Beads Framework Integration Notes

When suggesting code for Beads integration:

### HTTP Client Usage in Beads
- Beads has native HTTP support - show how to use this library from Beads code
- Emphasize proper error handling patterns
- Show context usage for timeouts and cancellation

### Type Safety
- Beads benefits from Go's type safety - highlight type-safe API patterns
- Show how to handle optional fields and nullable values
- Demonstrate validation before API calls

### Rate Limiting
- Show how to configure rate limiting appropriately
- Explain burst vs. sustained rate limits
- Demonstrate backoff strategies

### Example Beads Integration Pattern

```go
// Go library usage that Beads apps can follow
client := manapool.NewClient(
    getEnvVar("MANAPOOL_TOKEN"),
    getEnvVar("MANAPOOL_EMAIL"),
    manapool.WithRateLimit(10, 1), // 10 req/sec, burst 1
    manapool.WithTimeout(30*time.Second),
)

// Use context for cancellation in Beads UI
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

inventory, err := client.GetSellerInventory(ctx, opts)
if err != nil {
    // Show error in Beads UI
    displayError(err)
    return
}
// Update Beads UI with inventory
updateUI(inventory)
```

## Version Control

### Commit Message Format
```
<type>: <description>

Examples:
feat: add order management endpoints
fix: correct rate limiter initialization
docs: update README with new examples
test: add tests for error handling
refactor: simplify validation logic
chore: update dependencies
```

## Resources

- Main Instructions: `.github/copilot-instructions.md`
- Project README: `README.md`
- CI/CD Docs: `CI_CD.md`
- Logging Docs: `LOGGING.md`
- Go Documentation: https://pkg.go.dev/github.com/repricah/manapool
- Beads Framework: http://www.beadsproject.net/
