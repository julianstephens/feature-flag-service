package response

import (
	"encoding/json"
	"net/http"
)

// JSONEncoder implements the Encoder interface for JSON responses.
type JSONEncoder struct {
	// Indent specifies the JSON indentation string. Empty string means no indentation.
	Indent string
}

// NewJSONEncoder creates a new JSONEncoder with default settings.
func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{}
}

// NewJSONEncoderWithIndent creates a new JSONEncoder with pretty-printing enabled.
func NewJSONEncoderWithIndent(indent string) *JSONEncoder {
	return &JSONEncoder{Indent: indent}
}

// Encode encodes the given value as JSON and writes it to the response writer.
// It sets the Content-Type header to application/json.
func (j *JSONEncoder) Encode(w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json")
	
	encoder := json.NewEncoder(w)
	if j.Indent != "" {
		encoder.SetIndent("", j.Indent)
	}
	
	return encoder.Encode(v)
}