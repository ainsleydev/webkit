package middleware

import (
	"github.com/google/uuid"

	webkitctx "github.com/ainsleydev/webkit/pkg/context"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

const (
	// RequestIDHeader is the name of the HTTP Header which contains the request id.
	// Exported so that it can be changed by developers
	RequestIDHeader = "X-Request-ID"
)

// RequestID assigns a unique identifier to the contact under RequestIDContextKey.
// The ID is also sent back in the payload to the calling client.
func RequestID(next webkit.Handler) webkit.Handler {
	return func(c *webkit.Context) error {
		requestID := c.Request.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx := webkitctx.WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)
		return next(c)
	}
}
