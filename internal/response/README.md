# Response Package

The `internal/response` package provides structured HTTP response handling for mux-based applications with support for extensible encoding, hooks, and error handling.

## Features

- **Extensible Encoding**: Pluggable encoder interface (JSON by default)
- **Hook System**: Pre/post-processing hooks for logging, metrics, etc.
- **Error Handling**: Centralized error response handling
- **Convenience Constructors**: Multiple pre-configured responder types
- **Full Test Coverage**: Comprehensive test suite with examples

## Quick Start

```go
import "github.com/julianstephens/feature-flag-service/internal/response"

// Basic usage
responder := response.New()

// In your handler
func handler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{"message": "Hello, World!"}
    responder.Write(w, r, data)
}
```

## API Reference

### Types

- `Responder` - Main struct for handling responses
- `Encoder` - Interface for response encoding
- `BeforeFunc`, `AfterFunc`, `OnErrorFunc` - Hook function types

### Constructors

- `New()` - Basic responder with JSON encoding
- `NewWithLogging()` - Responder with request/response logging
- `NewCustom()` - Fully customizable responder

### Methods

- `Write(w, r, data)` - Write successful response
- `Error(w, r, err)` - Handle error response

## Files

- `response.go` - Core types and interfaces
- `json.go` - JSON encoder implementation
- `hooks.go` - Default hook implementations
- `responder.go` - Convenience constructors
- `doc.go` - Package documentation
- `response_test.go` - Test suite
- `example_test.go` - Usage examples
- `integration_example.md` - Integration guide

## Testing

```bash
go test ./internal/response -v
```

## Examples

See `example_test.go` for runnable examples and `integration_example.md` for integration patterns.