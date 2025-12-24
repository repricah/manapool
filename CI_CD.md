# CI/CD Strategy for Manapool Go Client

## Current Setup (Monorepo)

### Test Execution

While the library is in the `pkg/manapool/` directory of the tcg-repricer monorepo, tests are executed as part of the main project's CI pipeline:

```bash
# Main project's CI runs this, which includes pkg/manapool tests
go test ./...
```

**GitHub Actions Workflow:** removed while CI migrates to a Docker-in-Docker runner (see `docs/technical/docker-in-docker-actions.md`).

Use `make test` / `make ci-local` locally to continue running:
- ✅ All tests including `pkg/manapool`
- ✅ Coverage checks
- ✅ GolangCI-Lint
- ✅ Multi-version testing can be restored in a future workflow once the runner is online

### Benefits of Monorepo Testing

1. **No duplication** - Tests run once for entire codebase
2. **Integration testing** - Can test manapool client usage in tcg-repricer
3. **Shared CI configuration** - One workflow to maintain
4. **Faster feedback** - Single pipeline for all changes

### Running Tests Locally

```bash
# Run all tests (including manapool)
go test ./...

# Run only manapool tests
go test ./pkg/manapool/...

# Run with coverage
go test -cover ./pkg/manapool/...

# Run with verbose output
go test -v ./pkg/manapool/...
```

---

## After Repository Extraction

When the library is extracted to `github.com/repricah/manapool-tools`, it will need its own dedicated CI/CD pipeline.

### Proposed GitHub Actions Workflow

Location: `.github/workflows/ci.yml` (in new repository)

```yaml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.23', '1.24', '1.25']

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}

      - name: Download dependencies
        run: go mod download

      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Check coverage
        run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $coverage%"
          # Fail if coverage drops below 95%
          if (( $(echo "$coverage < 95.0" | bc -l) )); then
            echo "❌ Coverage is below 95%"
            exit 1
          fi

      - name: Upload coverage to Codecov
        if: matrix.go-version == '1.25'
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out

  lint:
    name: Lint
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest

  build:
    name: Build
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Build
        run: go build -v ./...
```

### Additional Workflows

**Release Workflow** (`.github/workflows/release.yml`):

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    name: Create Release
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Run tests
        run: go test -v ./...

      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          generate_release_notes: true
```

### Branch Protection Rules

After extraction, configure the following for the `main` branch:

- ✅ Require pull request reviews (1 reviewer minimum)
- ✅ Require status checks to pass (all CI jobs)
- ✅ Require branches to be up to date
- ✅ Require conversation resolution
- ✅ Do not allow force pushes
- ✅ Do not allow deletions

---

## Test Coverage Requirements

### Current Coverage: 96.5%

**Minimum Requirements:**
- Overall coverage: ≥ 95%
- Critical paths (client, inventory, account): ≥ 98%
- Error handling: 100%
- Public API: 100%

**Coverage Tracking:**

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage by function
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# View in browser
open coverage.html
```

### CI Coverage Enforcement

The CI pipeline should:
1. **Measure coverage** on every PR
2. **Fail the build** if coverage drops below 95%
3. **Report coverage** as a PR comment
4. **Upload to Codecov** for historical tracking

---

## Integration Testing

### Mock Server Testing

The current test suite uses `httptest.NewServer()` to mock the Manapool API:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Mock API response
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"username": "test"}`))
}))
defer server.Close()

client := manapool.NewClient(token, email,
    manapool.WithBaseURL(server.URL+"/"),
)
```

**Benefits:**
- ✅ Fast (no network calls)
- ✅ Reliable (no API downtime)
- ✅ Comprehensive (can test error cases)
- ✅ No credentials needed

### Optional: Live API Testing

For integration testing against the real API (optional):

```bash
# Run with live API credentials (use with caution)
export MANAPOOL_API_TOKEN="your-token"
export MANAPOOL_API_EMAIL="your-email"
export MANAPOOL_RUN_INTEGRATION_TESTS="true"

go test -v -tags=integration ./...
```

**Note:** Live API tests should:
- Be opt-in (require environment variable)
- Use separate build tag (`//go:build integration`)
- Not run in CI by default
- Clean up any created resources

---

## Dependency Management

### Go Modules

The library uses Go modules with minimal dependencies:

```
require (
    golang.org/x/time v0.5.0
)
```

**Dependency Updates:**
- Monitor security advisories
- Update dependencies monthly
- Run `go mod tidy` after updates
- Test thoroughly after updates

### Dependabot Configuration

`.github/dependabot.yml`:

```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5
```

---

## Pre-commit Checks

### Recommended Pre-commit Hook

`.git/hooks/pre-commit`:

```bash
#!/bin/bash
set -e

echo "Running pre-commit checks..."

# Format check
echo "Checking formatting..."
if [ -n "$(gofmt -l .)" ]; then
    echo "❌ Code is not formatted. Run: gofmt -w ."
    exit 1
fi

# Run tests
echo "Running tests..."
go test ./...

# Run linter (if installed)
if command -v golangci-lint &> /dev/null; then
    echo "Running linter..."
    golangci-lint run
fi

echo "✅ All pre-commit checks passed!"
```

---

## Performance Benchmarks

### Benchmark Tests

Add benchmark tests for critical paths:

```go
func BenchmarkGetSellerInventory(b *testing.B) {
    server := setupMockServer()
    defer server.Close()

    client := manapool.NewClient("token", "email",
        manapool.WithBaseURL(server.URL+"/"),
    )

    ctx := context.Background()
    opts := manapool.InventoryOptions{Limit: 500}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = client.GetSellerInventory(ctx, opts)
    }
}
```

**Run benchmarks:**

```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkGetSellerInventory -benchmem ./...

# Compare benchmarks
go test -bench=. -benchmem -count=5 ./... > old.txt
# ... make changes ...
go test -bench=. -benchmem -count=5 ./... > new.txt
benchstat old.txt new.txt
```

---

## Security Scanning

### Recommended Tools

1. **gosec** - Security scanner for Go code
   ```bash
   gosec ./...
   ```

2. **nancy** - Dependency vulnerability scanner
   ```bash
   go list -json -m all | nancy sleuth
   ```

3. **govulncheck** - Official Go vulnerability checker
   ```bash
   govulncheck ./...
   ```

### GitHub Security Features

Enable in repository settings:
- ✅ Dependabot alerts
- ✅ Code scanning (CodeQL)
- ✅ Secret scanning
- ✅ Security policy (SECURITY.md)

---

## Documentation Generation

### GoDoc

The library is documented using GoDoc comments. Preview locally:

```bash
# Install godoc (if not already installed)
go install golang.org/x/tools/cmd/godoc@latest

# Start local documentation server
godoc -http=:6060

# Open in browser
open http://localhost:6060/pkg/github.com/repricah/manapool-tools/
```

### Example Code

All public functions include example code in tests:

```go
func ExampleClient_GetSellerAccount() {
    client := manapool.NewClient("token", "email")
    account, err := client.GetSellerAccount(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(account.Username)
}
```

---

## Monitoring & Observability

### Recommended Practices

When using the library in production:

1. **Enable logging** with custom logger
2. **Monitor rate limits** via error types
3. **Track retry attempts** for debugging
4. **Measure latency** of API calls
5. **Alert on errors** above threshold

### Example Logging Setup

```go
type productionLogger struct {
    logger *log.Logger
}

func (l *productionLogger) Debugf(format string, args ...interface{}) {
    l.logger.Printf("[DEBUG] "+format, args...)
}

func (l *productionLogger) Errorf(format string, args ...interface{}) {
    l.logger.Printf("[ERROR] "+format, args...)
    // Send to error tracking service (e.g., Sentry)
}

client := manapool.NewClient(token, email,
    manapool.WithLogger(&productionLogger{logger: log.Default()}),
)
```

---

## Summary

### Current State (Monorepo)
- ✅ Tests run as part of main project CI
- ✅ No separate workflow needed
- ✅ 96.5% coverage maintained
- ✅ Comprehensive test suite

### After Extraction
- 🔜 Dedicated GitHub Actions workflows
- 🔜 Independent release cycle
- 🔜 Separate coverage tracking
- 🔜 Branch protection rules
- 🔜 Automated releases

### Next Steps
1. Continue development in monorepo
2. When ready for v1.0.0, extract to separate repo
3. Setup CI/CD workflows using templates above
4. Configure branch protection and security
5. Publish first stable release
