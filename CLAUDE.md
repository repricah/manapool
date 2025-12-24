# Manapool Go Client - Claude Agent Instructions

## Primary Instructions

For comprehensive project documentation and coding standards, see:
**→ `.github/copilot-instructions.md`**

This file serves as the Claude-specific entry point. All detailed coding standards, architecture patterns, testing requirements, and project conventions are maintained in the copilot instructions file linked above.

## Quick Project Context

**Project**: Go client library for the Manapool API (Magic: The Gathering inventory management)  
**Module**: `github.com/repricah/manapool`  
**Language**: Go 1.24.7  
**Status**: Pre-release v0.2.0  
**Test Coverage**: 96.5%

## Essential Commands

```bash
# Run tests
go test -v ./...

# Run tests with race detector
go test -race -v ./...

# Check coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Verify module
go mod verify
go mod tidy
```

## Project Structure

```
manapool/
├── .github/
│   ├── copilot-instructions.md  # ← Main instructions file
│   ├── agents/
│   │   └── gemini.md            # Gemini-specific instructions
│   └── workflows/
│       └── ci.yml               # CI configuration
├── client.go                    # HTTP client implementation
├── account.go                   # Account endpoints
├── inventory.go                 # Inventory endpoints
├── types.go                     # API type definitions
├── errors.go                    # Error types
├── options.go                   # Client options
├── *_test.go                    # Unit tests
├── go.mod                       # Module definition
├── README.md                    # User documentation
├── CI_CD.md                     # CI/CD documentation
├── LOGGING.md                   # Logging documentation
└── CLAUDE.md                    # This file
```

## Core Principles

### 1. Zero Dependencies Philosophy
- Only one external dependency: `golang.org/x/time` for rate limiting
- Justify any new dependency thoroughly
- Prefer standard library solutions

### 2. High Test Coverage
- Maintain **>95% test coverage**
- Test both success and error paths
- Use table-driven tests
- Mock HTTP with `httptest`

### 3. Context-First Design
- Always accept `context.Context` as first parameter for I/O operations
- Never store context in structs
- Respect context cancellation

### 4. Strong Error Handling
- Use custom error types: `APIError`, `ValidationError`, `NetworkError`
- Provide helper methods: `IsNotFound()`, `IsUnauthorized()`, etc.
- Always wrap errors with context: `fmt.Errorf("operation failed: %w", err)`

### 5. Thread Safety
- Client is safe for concurrent use
- Use rate limiting to protect API
- No mutable shared state without synchronization

## Code Style Checklist

Before suggesting code changes, verify:

- ✅ Context is first parameter for I/O operations
- ✅ Errors are wrapped with context
- ✅ No panics (return errors instead)
- ✅ Tests cover both success and error cases
- ✅ godoc comments for all exported symbols
- ✅ Follows Go naming conventions (camelCase/PascalCase)
- ✅ Rate limiter is used for API calls
- ✅ Input validation before API calls
- ✅ No hardcoded credentials or tokens

## Common Patterns

### Client Method Pattern
```go
func (c *Client) DoOperation(ctx context.Context, opts Options) (*Result, error) {
    // 1. Validate inputs
    if err := opts.Validate(); err != nil {
        return nil, &ValidationError{Field: "opts", Message: err.Error()}
    }
    
    // 2. Wait for rate limiter
    if err := c.rateLimiter.Wait(ctx); err != nil {
        return nil, err
    }
    
    // 3. Make request
    var result Result
    err := c.doRequest(ctx, "GET", "/endpoint", opts, &result)
    if err != nil {
        return nil, fmt.Errorf("failed to do operation: %w", err)
    }
    
    return &result, nil
}
```

### Test Pattern
```go
func TestClient_DoOperation(t *testing.T) {
    tests := []struct {
        name       string
        opts       Options
        setupMock  func(*httptest.Server)
        wantErr    bool
        errType    interface{}
    }{
        {
            name: "success",
            opts: Options{Field: "value"},
            setupMock: func(s *httptest.Server) {
                // Configure mock response
            },
            wantErr: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Things to Never Do

- ❌ Make real API calls in tests
- ❌ Store context in structs
- ❌ Use `panic()` in library code
- ❌ Add dependencies without justification
- ❌ Commit credentials or tokens
- ❌ Ignore context cancellation
- ❌ Skip error wrapping
- ❌ Reduce test coverage

## Beads Framework Integration

When this library is used with Beads applications:

### Key Integration Points

1. **HTTP Client Patterns**: Demonstrate proper error handling and retries
2. **Type Safety**: Show strongly-typed API interactions
3. **Rate Limiting**: Configure appropriate rate limits for API usage
4. **Context Propagation**: Use context for cancellation and timeouts

### Beads-Specific Guidance

```go
// Initialize client with environment config
client := manapool.NewClient(
    os.Getenv("MANAPOOL_TOKEN"),
    os.Getenv("MANAPOOL_EMAIL"),
    manapool.WithTimeout(30*time.Second),
    manapool.WithRateLimit(10, 1), // 10 req/sec
)

// Use context for Beads UI operations
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Handle errors for Beads UI
inventory, err := client.GetSellerInventory(ctx, opts)
if err != nil {
    var apiErr *manapool.APIError
    if errors.As(err, &apiErr) {
        // Display user-friendly error in Beads UI
        if apiErr.IsRateLimited() {
            showError("Too many requests. Please wait.")
        } else {
            showError(fmt.Sprintf("API error: %s", apiErr.Message))
        }
        return
    }
    // Handle other errors
    showError("Connection error. Please try again.")
    return
}

// Update Beads UI with data
updateInventoryDisplay(inventory)
```

### Beads Resources
- Official Site: http://www.beadsproject.net/
- Examples: https://github.com/magicmouse/beads-examples
- Documentation: http://www.beadsproject.net/ref/

## Git Workflow

### Commit Format
```
<type>: <short description>

Types: feat, fix, docs, test, refactor, chore
```

### Before Committing
1. Run tests: `go test -v ./...`
2. Check race conditions: `go test -race -v ./...`
3. Verify coverage: `go test -cover ./...`
4. Format code: `gofmt -w .`
5. Tidy modules: `go mod tidy`

## Documentation Files

- **This file** (CLAUDE.md): Claude agent entry point
- **.github/copilot-instructions.md**: Complete coding standards and conventions
- **.github/agents/gemini.md**: Gemini-specific instructions
- **README.md**: User-facing documentation and examples
- **CI_CD.md**: Continuous integration and deployment guide
- **LOGGING.md**: Logging configuration and best practices

## API Documentation

Full API documentation available at:
https://pkg.go.dev/github.com/repricah/manapool

## Support & Resources

- 🐛 Issues: https://github.com/repricah/manapool/issues
- 💬 Discussions: https://github.com/repricah/manapool/discussions
- 📖 Go Docs: https://pkg.go.dev/github.com/repricah/manapool

## Version Notes

**Current**: v0.2.0 (Pre-release)  
- API may change before v1.0.0
- Backwards compatibility is a goal but not guaranteed
- See README.md for changelog and planned features

---

**Remember**: Always refer to `.github/copilot-instructions.md` for comprehensive guidance. This file provides quick context for Claude-specific workflows.
