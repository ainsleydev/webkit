package httputil

import (
	"bytes"
	"net/http"
)

// ResponseRecorder wraps an http.ResponseWriter to capture the response
// body and status code for inspection.
type ResponseRecorder struct {
	http.ResponseWriter
	Status int
	Body   *bytes.Buffer
}

// NewResponseRecorder creates a new ResponseRecorder with the given ResponseWriter.
// It creates a new buffer for the response body.
func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		Body:           &bytes.Buffer{},
	}
}

// Write writes the response body to the captured buffer and the underlying ResponseWriter.
func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.Body.Write(b)
	return r.ResponseWriter.Write(b)
}

// WriteHeader captures the Status code and writes it to the underlying ResponseWriter.
func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.Status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
