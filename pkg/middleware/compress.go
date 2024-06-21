package middleware

import (
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

// Gzip returns a middleware that compresses HTTP responses using gzip compression.
// It wraps the provided handler, adding gzip compression to responses
// based on the specified configuration.
func Gzip(next webkit.Handler) webkit.Handler {
	return webkit.WrapMiddleware(next, middleware.Compress(5))
}
