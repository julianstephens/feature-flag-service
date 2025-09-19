// Package response provides structured HTTP response handling for mux-based applications.
// It includes support for extensible encoding, pre/post-processing hooks, and error handling.
package response

import (
	"net/http"
)

// Encoder defines the interface for encoding response data to an HTTP response writer.
type Encoder interface {
	// Encode writes the given value to the response writer.
	// It should set appropriate headers and write the encoded data.
	Encode(w http.ResponseWriter, v any) error
}

// BeforeFunc is called before encoding the response.
// It can be used for logging, metrics, or modifying the response writer.
type BeforeFunc func(w http.ResponseWriter, r *http.Request, data any)

// AfterFunc is called after successfully encoding the response.
// It can be used for cleanup, logging, or metrics collection.
type AfterFunc func(w http.ResponseWriter, r *http.Request, data any)

// OnErrorFunc is called when an error occurs during response processing.
// It should handle the error appropriately, typically by writing an error response.
type OnErrorFunc func(w http.ResponseWriter, r *http.Request, err error)

// Responder provides structured HTTP response handling with extensible hooks.
type Responder struct {
	Encoder Encoder     // Encoder for response data
	Before  BeforeFunc  // Hook called before encoding
	After   AfterFunc   // Hook called after successful encoding
	OnError OnErrorFunc // Hook called on encoding errors
}

// Write encodes and writes a successful response using the configured encoder and hooks.
// It calls Before hook, encodes the data, then calls After hook on success or OnError on failure.
func (r *Responder) Write(w http.ResponseWriter, req *http.Request, data any) {
	if r.Before != nil {
		r.Before(w, req, data)
	}

	if err := r.Encoder.Encode(w, data); err != nil {
		if r.OnError != nil {
			r.OnError(w, req, err)
		}
		return
	}

	if r.After != nil {
		r.After(w, req, data)
	}
}

// Error handles error responses by calling the OnError hook.
// If no OnError hook is configured, it writes a basic 500 Internal Server Error response.
func (r *Responder) Error(w http.ResponseWriter, req *http.Request, err error) {
	if r.OnError != nil {
		r.OnError(w, req, err)
		return
	}

	// Default error handling if no OnError hook is provided
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}