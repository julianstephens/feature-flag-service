# Integration Example: Using response package with existing server

This example shows how the new `internal/response` package can be integrated with the existing server code to provide structured HTTP response handling.

## Before (existing code in internal/server/server.go):

```go
router.HandleFunc("/checkhealth", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
})
```

## After (using the response package):

```go
import (
    "github.com/julianstephens/feature-flag-service/internal/response"
)

func StartREST(addr string, flagSvc flag.Service) error {
    // Create a responder with logging for better observability
    responder := response.NewWithLogging()
    
    router := mux.NewRouter()
    
    // Enhanced health check with structured response
    router.HandleFunc("/checkhealth", func(w http.ResponseWriter, r *http.Request) {
        healthData := map[string]interface{}{
            "status": "healthy",
            "service": "feature-flag-service",
            "timestamp": time.Now().UTC(),
        }
        responder.Write(w, r, healthData)
    })
    
    // Flag CRUD endpoints with structured error handling
    router.HandleFunc("/v1/flags", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            flags, err := flagSvc.ListFlags(r.Context())
            if err != nil {
                responder.Error(w, r, err)
                return
            }
            responder.Write(w, r, map[string]interface{}{
                "flags": flags,
                "count": len(flags),
            })
            
        case http.MethodPost:
            // Parse request body, create flag
            // On error: responder.Error(w, r, err)
            // On success: responder.Write(w, r, createdFlag)
        }
    })
    
    // ... rest of server setup
}
```

## Benefits of this integration:

1. **Consistent Response Format**: All endpoints return properly formatted JSON
2. **Automatic Logging**: Request/response logging with the logging responder
3. **Error Handling**: Centralized error response formatting
4. **Extensibility**: Easy to add metrics, authentication, etc. via hooks
5. **Maintainability**: Clean separation of response handling logic

## No Breaking Changes

The new package is purely additive - existing code continues to work unchanged while new endpoints can benefit from the structured response handling.