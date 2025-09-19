package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResponder_Write(t *testing.T) {
	responder := New()

	// Test data
	data := map[string]string{"message": "test"}

	// Create test request and response recorder
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Call Write
	responder.Write(w, req, data)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, "test") {
		t.Errorf("Expected response body to contain 'test', got %s", body)
	}
}

func TestResponder_Error(t *testing.T) {
	responder := New()

	// Create test request and response recorder
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Call Error
	err := errors.New("test error")
	responder.Error(w, req, err)

	// Verify error response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestJSONEncoder(t *testing.T) {
	encoder := NewJSONEncoder()

	// Test data
	data := map[string]string{"key": "value"}

	// Create response recorder
	w := httptest.NewRecorder()

	// Encode
	err := encoder.Encode(w, data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Verify JSON content
	body := w.Body.String()
	if !strings.Contains(body, "key") || !strings.Contains(body, "value") {
		t.Errorf("Expected JSON content, got %s", body)
	}
}

func TestJSONEncoderWithIndent(t *testing.T) {
	encoder := NewJSONEncoderWithIndent("  ")

	// Test data
	data := map[string]string{"key": "value"}

	// Create response recorder
	w := httptest.NewRecorder()

	// Encode
	err := encoder.Encode(w, data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify indented JSON
	body := w.Body.String()
	if !strings.Contains(body, "  ") {
		t.Errorf("Expected indented JSON, got %s", body)
	}
}

func TestHooks(t *testing.T) {
	beforeCalled := false
	afterCalled := false
	errorCalled := false

	responder := &Responder{
		Encoder: NewJSONEncoder(),
		Before: func(w http.ResponseWriter, r *http.Request, data any) {
			beforeCalled = true
		},
		After: func(w http.ResponseWriter, r *http.Request, data any) {
			afterCalled = true
		},
		OnError: func(w http.ResponseWriter, r *http.Request, err error) {
			errorCalled = true
		},
	}

	// Test successful write
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	responder.Write(w, req, map[string]string{"test": "data"})

	if !beforeCalled {
		t.Error("Before hook was not called")
	}
	if !afterCalled {
		t.Error("After hook was not called")
	}
	if errorCalled {
		t.Error("OnError hook should not have been called")
	}

	// Reset flags and test error
	beforeCalled = false
	afterCalled = false
	errorCalled = false

	w = httptest.NewRecorder()
	responder.Error(w, req, errors.New("test error"))

	if !errorCalled {
		t.Error("OnError hook should have been called during Error method")
	}
}

func TestNewConstructors(t *testing.T) {
	// Test New()
	r1 := New()
	if r1.Encoder == nil {
		t.Error("New() should provide a default encoder")
	}
	if r1.Before == nil || r1.After == nil || r1.OnError == nil {
		t.Error("New() should provide default hooks")
	}

	// Test NewWithLogging()
	r2 := NewWithLogging()
	if r2.Encoder == nil {
		t.Error("NewWithLogging() should provide a default encoder")
	}
	if r2.Before == nil || r2.After == nil || r2.OnError == nil {
		t.Error("NewWithLogging() should provide logging hooks")
	}

	// Test NewCustom()
	r3 := NewCustom(nil, nil, nil, nil)
	if r3.Encoder == nil {
		t.Error("NewCustom() should provide defaults when nil is passed")
	}
}

// Test convenience methods for status codes
func TestResponder_StatusMethods(t *testing.T) {
	responder := New()
	data := map[string]string{"message": "test"}

	tests := []struct {
		name           string
		method         func(w http.ResponseWriter, req *http.Request, data any)
		expectedStatus int
	}{
		{"OK", responder.OK, http.StatusOK},
		{"Created", responder.Created, http.StatusCreated},
		{"BadRequest", responder.BadRequest, http.StatusBadRequest},
		{"Unauthorized", responder.Unauthorized, http.StatusUnauthorized},
		{"NotFound", responder.NotFound, http.StatusNotFound},
		{"InternalServerError", responder.InternalServerError, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			tt.method(w, req, data)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			body := w.Body.String()
			if !strings.Contains(body, "test") {
				t.Errorf("Expected response body to contain 'test', got %s", body)
			}
		})
	}
}
