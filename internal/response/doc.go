/*
Package response provides structured HTTP response handling for mux-based applications.

This package offers a flexible way to handle HTTP responses with support for:
- Custom encoding strategies (JSON by default)
- Pre-processing hooks (before encoding)
- Post-processing hooks (after encoding)
- Error handling hooks

Basic Usage:

	package main

	import (
		"net/http"
		"github.com/gorilla/mux"
		"github.com/julianstephens/feature-flag-service/internal/response"
	)

	func main() {
		responder := response.New()

		router := mux.NewRouter()
		router.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
			data := map[string]string{"message": "Hello, World!"}
			responder.Write(w, r, data)
		})

		http.ListenAndServe(":8080", router)
	}

Advanced Usage with Custom Hooks:

	func main() {
		// Custom hooks for request/response logging and metrics
		beforeHook := func(w http.ResponseWriter, r *http.Request, data any) {
			log.Printf("Processing request: %s %s", r.Method, r.URL.Path)
		}

		afterHook := func(w http.ResponseWriter, r *http.Request, data any) {
			log.Printf("Response sent for: %s %s", r.Method, r.URL.Path)
		}

		errorHook := func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Error in %s %s: %v", r.Method, r.URL.Path, err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}

		responder := response.NewCustom(
			response.NewJSONEncoderWithIndent("  "), // Pretty JSON
			beforeHook,
			afterHook,
			errorHook,
		)

		// Use responder in handlers...
	}

Error Handling:

	func handler(w http.ResponseWriter, r *http.Request) {
		data, err := processRequest(r)
		if err != nil {
			// Handle errors using the responder's error method
			responder.Error(w, r, err)
			return
		}

		// Write successful response
		responder.Write(w, r, data)
	}

The package is designed to be extensible and can be easily integrated with existing
gorilla/mux applications while providing a consistent and maintainable approach
to HTTP response handling.
*/
package response
