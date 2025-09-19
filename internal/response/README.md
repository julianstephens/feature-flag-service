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

// Using status code convenience methods
func statusHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{"status": "success"}
    responder.OK(w, r, data)  // 200 OK
}

func createHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{"id": "123", "created": true}
    responder.Created(w, r, data)  // 201 Created
}

func authHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{"error": "access denied"}
    responder.Unauthorized(w, r, data)  // 401 Unauthorized
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

#### Status Code Convenience Methods

- `OK(w, r, data)` - Write response with 200 OK status
- `Created(w, r, data)` - Write response with 201 Created status
- `BadRequest(w, r, data)` - Write response with 400 Bad Request status
- `Unauthorized(w, r, data)` - Write response with 401 Unauthorized status
- `NotFound(w, r, data)` - Write response with 404 Not Found status
- `InternalServerError(w, r, data)` - Write response with 500 Internal Server Error status

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