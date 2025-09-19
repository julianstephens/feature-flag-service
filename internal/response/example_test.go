package response_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/julianstephens/feature-flag-service/internal/response"
)

// Example demonstrates basic usage of the response package with mux.
func Example() {
	// Create a new responder with default JSON encoding
	responder := response.New()

	// Create a new mux router
	router := mux.NewRouter()

	// Example handler using the responder
	router.HandleFunc("/api/greeting", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"message": "Hello, World!",
			"status":  "success",
		}
		responder.Write(w, r, data)
	})

	// Example error handler using the responder
	router.HandleFunc("/api/error", func(w http.ResponseWriter, r *http.Request) {
		responder.Error(w, r, fmt.Errorf("something went wrong"))
	})

	// Test the handlers
	req := httptest.NewRequest("GET", "/api/greeting", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Printf("Status: %d\n", w.Code)
	fmt.Printf("Content-Type: %s\n", w.Header().Get("Content-Type"))
	// Output:
	// Status: 200
	// Content-Type: application/json
}

// ExampleNewWithLogging demonstrates usage with logging hooks.
func ExampleNewWithLogging() {
	// Create a responder with logging hooks
	responder := response.NewWithLogging()

	// Create a simple handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{"result": "processed"}
		responder.Write(w, r, data)
	}

	// Test the handler
	req := httptest.NewRequest("GET", "/api/process", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	fmt.Printf("Status: %d\n", w.Code)
	// Output:
	// Status: 200
}

// ExampleNewCustom demonstrates usage with custom hooks and encoding.
func ExampleNewCustom() {
	// Create a responder with custom pretty-printed JSON
	responder := response.NewCustom(
		response.NewJSONEncoderWithIndent("  "), // Pretty JSON
		nil,                                     // Use default Before hook
		nil,                                     // Use default After hook
		nil,                                     // Use default OnError hook
	)

	handler := func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"data":   []string{"item1", "item2"},
			"status": "ok",
		}
		responder.Write(w, r, data)
	}

	// Test the handler
	req := httptest.NewRequest("GET", "/api/data", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	fmt.Printf("Status: %d\n", w.Code)
	// Output:
	// Status: 200
}

// Example_statusMethods demonstrates the convenience status methods.
func Example_statusMethods() {
	responder := response.New()

	// Example handlers using status methods
	okHandler := func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{"status": "success"}
		responder.OK(w, r, data)
	}

	createdHandler := func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"id":      "123",
			"created": true,
		}
		responder.Created(w, r, data)
	}

	unauthorizedHandler := func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{"error": "access denied"}
		responder.Unauthorized(w, r, data)
	}

	// Test OK method
	req := httptest.NewRequest("GET", "/success", nil)
	w := httptest.NewRecorder()
	okHandler(w, req)
	fmt.Printf("OK Status: %d\n", w.Code)

	// Test Created method
	req = httptest.NewRequest("POST", "/create", nil)
	w = httptest.NewRecorder()
	createdHandler(w, req)
	fmt.Printf("Created Status: %d\n", w.Code)

	// Test Unauthorized method
	req = httptest.NewRequest("GET", "/protected", nil)
	w = httptest.NewRecorder()
	unauthorizedHandler(w, req)
	fmt.Printf("Unauthorized Status: %d\n", w.Code)

	// Output:
	// OK Status: 200
	// Created Status: 201
	// Unauthorized Status: 401
}
