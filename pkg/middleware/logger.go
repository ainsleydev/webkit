package middleware

import (
	"github.com/ainsleydev/webkit/pkg/logger"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// Logger - TODO
func Logger(next webkit.Handler) webkit.Handler {
	return func(ctx *webkit.Context) error {
		req := ctx.Request
		logger.Info("Request: %s %s", req.Method, req.URL.Path)
		return next(ctx)
	}
}
