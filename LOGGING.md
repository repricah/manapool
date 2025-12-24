# Output and Logging Strategy

## Design Philosophy

This library follows **Go library best practices** for output and logging:

### ✅ What the Library Does

1. **Returns Errors** - All operations return errors for the caller to handle
2. **Optional Logging** - Provides a `Logger` interface that users can
   implement
3. **Silent by Default** - Uses a no-op logger that discards all messages
4. **No Direct Output** - Never uses `fmt.Printf()`, `log.Printf()`, or writes
   to stdout/stderr

### ❌ What the Library Does NOT Do

- ❌ Write directly to stdout/stderr
- ❌ Use the standard `log` package
- ❌ Make assumptions about how errors should be displayed
- ❌ Force logging configuration on users

## Rationale

### Why No Direct Output?

**Libraries should be quiet by default.** Here's why:

1. **Flexibility** - Users may run in environments without stdout
   (services, lambdas, etc.)
2. **Control** - Users should decide what gets logged and where
3. **Testing** - Direct output makes testing harder and pollutes test output
4. **Composability** - Libraries that output directly can't be used in
   contexts where output matters
5. **12-Factor App** - Logs should go to stdout, but the *application*
   controls this, not libraries

### Example of the Problem

```go
// ❌ BAD - Library writes directly to stdout
func (c *Client) GetAccount() (*Account, error) {
    fmt.Println("Fetching account...")  // Don't do this in libraries!
    // ...
}

// When a user calls this, they can't control the output:
account, _ := client.GetAccount()
// Output: "Fetching account..." <-- forced on the user
```

## The Correct Approach

### 1. Return Errors

```go
// ✅ GOOD - Return errors, let caller handle them
func (c *Client) GetSellerAccount(ctx context.Context) (*Account, error) {
    resp, err := c.doRequest(ctx, "GET", "/account", nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get seller account: %w", err)
    }
    // ...
}
```

**User decides how to handle:**

```go
account, err := client.GetSellerAccount(ctx)
if err != nil {
    // User chooses: log it, display it, ignore it, wrap it, etc.
    log.Printf("ERROR: %v", err)
    // or: return fmt.Errorf("getting account: %w", err)
    // or: sentry.CaptureException(err)
}
```

### 2. Optional Logger Interface

```go
// ✅ GOOD - Provide optional logging via interface
type Logger interface {
    Debugf(format string, args ...interface{})
    Errorf(format string, args ...interface{})
}

// Default: silent no-op logger
type noopLogger struct{}

func (l *noopLogger) Debugf(format string, args ...interface{}) {}
func (l *noopLogger) Errorf(format string, args ...interface{}) {}
```

**Users can enable logging when they want:**

```go
// Option 1: Use standard log package
type stdLogger struct{}

func (l *stdLogger) Debugf(format string, args ...interface{}) {
    log.Printf("[DEBUG] "+format, args...)
}

func (l *stdLogger) Errorf(format string, args ...interface{}) {
    log.Printf("[ERROR] "+format, args...)
}

client := manapool.NewClient(token, email,
    manapool.WithLogger(&stdLogger{}),
)

// Option 2: Use structured logging (e.g., zap, logrus)
type zapLogger struct {
    logger *zap.SugaredLogger
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
    l.logger.Debugf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
    l.logger.Errorf(format, args...)
}

client := manapool.NewClient(token, email,
    manapool.WithLogger(&zapLogger{logger: zapSugar}),
)

// Option 3: Custom logging (e.g., to metrics/monitoring)
type metricsLogger struct {
    metrics *MetricsClient
}

func (l *metricsLogger) Debugf(format string, args ...interface{}) {
    // Do nothing for debug
}

func (l *metricsLogger) Errorf(format string, args ...interface{}) {
    l.metrics.IncrementCounter("manapool.errors")
}
```

## Common Patterns

### Application-Level Logging

```go
func main() {
    // Application controls logging setup
    log.SetOutput(os.Stdout)
    log.SetPrefix("[myapp] ")

    client := manapool.NewClient(token, email,
        manapool.WithLogger(&myLogger{}),
    )

    account, err := client.GetSellerAccount(ctx)
    if err != nil {
        // Application decides how to handle errors
        log.Printf("Failed to get account: %v", err)
        os.Exit(1)
    }

    // Application controls output format
    fmt.Printf("Account: %s (%s)\n", account.Username, account.Email)
}
```

### Production Service Logging

```go
func (s *Service) SyncInventory(ctx context.Context) error {
    // Structured logging with context
    logger := s.logger.WithField("operation", "sync_inventory")

    client := manapool.NewClient(token, email,
        manapool.WithLogger(newStructuredLogger(logger)),
    )

    inventory, err := client.GetSellerInventory(ctx, opts)
    if err != nil {
        // Log with structured context
        logger.WithError(err).Error("Failed to fetch inventory")

        // Also send to error tracking
        sentry.CaptureException(err)

        return fmt.Errorf("sync inventory: %w", err)
    }

    logger.WithField("count", len(inventory.Inventory)).Info("Synced inventory")
    return nil
}
```

### Testing

```go
func TestSyncInventory(t *testing.T) {
    // No logging output during tests (unless you want it)
    client := manapool.NewClient("test-token", "test@example.com")

    // Or: capture logs during tests
    testLogger := &testLogger{}
    client := manapool.NewClient("test-token", "test@example.com",
        manapool.WithLogger(testLogger),
    )

    // Test behavior, check logs if needed
    if len(testLogger.errors) > 0 {
        t.Errorf("unexpected errors: %v", testLogger.errors)
    }
}
```

## Examples in Documentation

### README Examples

Documentation examples (like in README.md) **can use fmt.Printf** because
they're demonstrating usage:

```go
// ✅ OK in README examples - showing output to users
account, err := client.GetSellerAccount(ctx)
if err != nil {
    log.Fatal(err)  // OK in examples
}
fmt.Printf("Account: %s\n", account.Username)  // OK in examples
```

### GoDoc Examples

```go
func ExampleClient_GetSellerAccount() {
    client := manapool.NewClient("token", "email")
    account, err := client.GetSellerAccount(context.Background())
    if err != nil {
        log.Fatal(err)  // OK in examples
    }
    fmt.Println(account.Username)  // OK in examples
    // Output: testuser
}
```

## Verification

To verify the library has no direct output:

```bash
# Check for fmt.Print* in library code (excluding tests and comments)
grep -rn "fmt\.Print" --include="*.go" --exclude="*_test.go" . | grep -v "//"

# Check for log.Print* in library code
grep -rn "log\.Print" --include="*.go" --exclude="*_test.go" . | grep -v "//"

# Should return no results (except in comments)
```

**Current Status:** ✅ Clean - No direct output in library code

## Summary

| What | Where | Why |
|------|-------|-----|
| **Error Returns** | All library functions | Caller handles errors |
| **Logger Interface** | Optional via `WithLogger()` | User-controlled logging |
| **No-op Logger** | Default | Library is silent by default |
| **fmt.Printf in Examples** | README, GoDoc | OK for documentation/examples |
| **Direct Output** | ❌ Never in library | Violates library design principles |

## References

- [Go Proverbs](https://go-proverbs.github.io/): "Clear is better than clever"
- [Effective Go](https://golang.org/doc/effective_go.html): Error handling
  section
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments):
  Logging section
- [12-Factor App](https://12factor.net/logs): Logs as event streams
