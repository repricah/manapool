# Manapool Go Client - GitHub Copilot Instructions

## Project Overview

This is a Go client library for the Manapool API, which manages Magic: The Gathering inventory. The library provides a type-safe, production-ready interface with automatic rate limiting, retries, and comprehensive error handling.

**Current Version**: v0.2.0 (Pre-release)  
**Repository**: github.com/repricah/manapool  
**Primary Author**: jblotus

## Tech Stack & Dependencies

- **Language**: Go 1.24.7
- **Dependencies**: 
  - `golang.org/x/time` - For rate limiting
- **Zero external dependencies** philosophy - Keep it simple and maintainable

## Project Structure

```
├── client.go          # Main HTTP client implementation
├── account.go         # Seller account endpoints
├── inventory.go       # Inventory management endpoints
├── types.go           # API type definitions
├── errors.go          # Error handling types
├── options.go         # Client configuration options
├── *_test.go          # Comprehensive unit tests (96.5% coverage)
├── go.mod             # Module definition
├── README.md          # Documentation
├── CI_CD.md           # CI/CD documentation
└── LOGGING.md         # Logging documentation
```

## Coding Standards

### General Go Practices

- Follow standard Go conventions and idioms
- Use `gofmt` for code formatting (automatically applied)
- Follow effective Go guidelines
- Keep code simple and readable - favor clarity over cleverness
- Write self-documenting code with clear variable and function names

### Naming Conventions

- Use camelCase for unexported names
- Use PascalCase for exported names
- Use descriptive names for functions and types
- Keep acronyms uppercase (e.g., `API`, `HTTP`, `ID`, `TCG`, `SKU`)

### Error Handling

- Always return errors, never panic in library code
- Use custom error types for domain-specific errors:
  - `APIError` - For HTTP API errors
  - `ValidationError` - For input validation errors
  - `NetworkError` - For network-related errors
- Provide helper methods on error types (e.g., `IsNotFound()`, `IsUnauthorized()`)
- Wrap errors with context using `fmt.Errorf` with `%w` verb

### Context Usage

- **Always** accept `context.Context` as the first parameter for operations that involve I/O
- Use context for cancellation, timeouts, and request-scoped values
- Never store context in structs
- Check for context cancellation in long-running operations

### Testing

- Maintain **high test coverage** (current: 96.5%)
- Write table-driven tests when testing multiple scenarios
- Use meaningful test names: `TestFunction_Scenario_ExpectedResult`
- Test error cases as thoroughly as success cases
- Use `t.Run()` for subtests to organize related test cases
- Mock HTTP responses for unit tests using `httptest`
- Never commit test code that makes real external API calls

### API Design

- Follow REST principles for API endpoint methods
- Use option pattern for configuration (see `options.go`)
- Provide sensible defaults for all options
- Make the zero value useful when possible
- Use functional options for client configuration

### Documentation

- Write godoc-style comments for all exported types, functions, and methods
- Include examples in documentation where helpful
- Keep README.md up to date with feature changes
- Document breaking changes clearly in commit messages

### Performance

- Avoid allocations in hot paths when possible
- Use rate limiting to respect API constraints
- Implement automatic retry with exponential backoff for transient failures
- Reuse HTTP client connections (client is safe for concurrent use)

## Git Workflow

### Commit Messages

- Use conventional commit format:
  - `feat:` for new features
  - `fix:` for bug fixes
  - `docs:` for documentation changes
  - `test:` for test additions/changes
  - `refactor:` for code refactoring
  - `chore:` for maintenance tasks

### Branching

- `main` branch is protected and requires PR review
- Create feature branches: `feature/<description>`
- Create fix branches: `fix/<description>`

### Pull Requests

- Run tests before submitting: `go test -v ./...`
- Run tests with race detector: `go test -race -v ./...`
- Check test coverage: `go test -coverprofile=coverage.out ./...`
- Ensure all tests pass in CI before merging
- Keep PRs focused and atomic

## CI/CD

- CI runs on every push and pull request
- Tests must pass before merge
- See `CI_CD.md` for detailed CI/CD documentation
- GoReleaser is used for release automation

## Dependencies

- **Minimize external dependencies** - Only add new dependencies if absolutely necessary
- Current dependencies are limited to `golang.org/x/time` for rate limiting
- Justify any new dependency additions in PR description

## Logging

- Use the `Logger` interface for all logging needs
- Default is a no-op logger (no output unless configured)
- Support custom logger injection via `WithLogger` option
- See `LOGGING.md` for logging best practices and examples

## Security

- Never commit API tokens or credentials
- Use environment variables for sensitive configuration
- Validate all inputs before making API calls
- Follow secure coding practices for HTTP clients

## Beads Framework Integration

This repository can be used as a reference for building Go libraries that integrate with the Beads programming framework:

### About Beads

Beads is a modern programming language and framework designed to simplify web development by replacing traditional JS/HTML/CSS stacks. Key features relevant to this project:

- **Specification-driven language**: Beads uses a declarative approach
- **Built-in sync engine**: Automatic UI updates and state management
- **Protected arithmetic**: Built-in safety for numerical operations
- **Proportional layouts**: Simplified responsive design
- **AI agent support**: Session memory, multimodal prompting, and autonomous task management

### Beads Integration Guidelines

When integrating this library with Beads applications:

1. **HTTP Client Patterns**: Beads supports HTTP operations natively - use this library's client patterns as reference for proper error handling and retries
2. **Type Safety**: Follow this library's approach to type-safe API interactions
3. **Rate Limiting**: Beads apps should respect API rate limits - use similar rate limiting patterns
4. **Context Propagation**: Maintain context throughout request chains for proper cancellation
5. **Error Handling**: Use structured error types that Beads UI can present meaningfully

### Beads Resources

- Official Site: http://www.beadsproject.net/
- Examples: https://github.com/magicmouse/beads-examples
- Quick Reference: http://www.beadsproject.net/ref/

## API Usage Patterns

### Client Initialization

```go
client := manapool.NewClient(
    "api-token",
    "email@example.com",
    manapool.WithTimeout(30*time.Second),
    manapool.WithRateLimit(10, 1),
    manapool.WithRetry(3, time.Second),
)
```

### Pagination

Always use the provided helper functions for pagination:
```go
err := manapool.IterateInventory(ctx, client, func(item *manapool.InventoryItem) error {
    // Process each item
    return nil
})
```

### Error Handling Pattern

```go
result, err := client.SomeOperation(ctx, params)
if err != nil {
    var apiErr *manapool.APIError
    if errors.As(err, &apiErr) {
        // Handle specific API errors
        if apiErr.IsNotFound() {
            // Handle not found
        }
    }
    return err
}
```

## Future Roadmap

Planned features tracked in GitHub issues:
- Additional lookup methods (Scryfall ID, Product ID)
- Order management endpoints
- Order fulfillment operations
- Inventory update operations (create, update, delete)
- Webhook support
- v1.0.0 stabilization

When implementing new features:
- Maintain backwards compatibility where possible
- Follow existing patterns and conventions
- Add comprehensive tests
- Update documentation
- Consider rate limiting implications
