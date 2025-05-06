package log

import (
	"context"
	"log/slog"

	webkitctx "github.com/ainsleydev/webkit/pkg/context"
)

type jsonHandler struct {
	slog.Handler
}

// Handle handles the logging record and formats it.
func (h jsonHandler) Handle(ctx context.Context, r slog.Record) error {
	// Add request ID if present
	if reqID, ok := ctx.Value(webkitctx.ContextKeyRequestID).(string); ok {
		r.AddAttrs(slog.String(string(webkitctx.ContextKeyRequestID), reqID))
	}

	// Collect all non-core attributes
	var attrs []any
	r.Attrs(func(a slog.Attr) bool {
		switch a.Key {
		case slog.TimeKey, slog.LevelKey, slog.MessageKey, string(webkitctx.ContextKeyRequestID):
			// Skip core keys
		default:
			attrs = append(attrs, a)
		}
		return true
	})

	// Add grouped attributes under "attr"
	if len(attrs) > 0 {
		r.AddAttrs(slog.Group("attr", attrs...))
	}

	return h.Handler.Handle(ctx, r)
}
