package middleware

import (
	"github.com/google/uuid"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

const (
	// RequestIDHeader is the name of the HTTP Header which contains the request id.
	// Exported so that it can be changed by developers
	RequestIDHeader = "X-Request-ID"
	// RequestIDContextKey is the key used to define a unique
	// identifier within echo context.
	RequestIDContextKey = "request_id"
)

// RequestID assigns a unique identifier to the contact under RequestIDContextKey.
// The ID is also sent back in the payload to the calling client.
func RequestID(next webkit.Handler) webkit.Handler {
	return func(ctx *webkit.Context) error {
		requestID := ctx.Request.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx.Set("request_id", requestID)
		return next(ctx)
	}
}
