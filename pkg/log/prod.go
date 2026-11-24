package log

import (
	"context"
	"io"
	"log/slog"

	webkitctx "github.com/ainsleydev/webkit/pkg/context"
)

// JSONHandlerOptions configures the production JSON handler.
type JSONHandlerOptions struct {
	slog.HandlerOptions

	// StderrLevel is the minimum level written to the stderr writer.
	// Records below it go to the stdout writer.
	//
	// Defaults to slog.LevelError. Platforms such as Vercel label lines
	// by stream (stdout is info, stderr is error), so sending WARN to
	// stderr would show warnings as errors there.
	StderrLevel slog.Leveler
}

// NewJSONHandler returns a handler that writes one JSON object per record.
//
// Attributes, including those added through With/WithAttrs, are written
// once at the top level. A request ID stored in the context is added as
// request_id. Records at or above opts.StderrLevel go to stderr, the rest
// to stdout.
func NewJSONHandler(stdout, stderr io.Writer, opts *JSONHandlerOptions) slog.Handler {
	if opts == nil {
		opts = &JSONHandlerOptions{}
	}
	hopts := opts.HandlerOptions
	hopts.ReplaceAttr = replaceLevelNames(hopts.ReplaceAttr)

	stderrLevel := opts.StderrLevel
	if stderrLevel == nil {
		stderrLevel = slog.LevelError
	}

	return &jsonHandler{
		stdout:      slog.NewJSONHandler(stdout, &hopts),
		stderr:      slog.NewJSONHandler(stderr, &hopts),
		stderrLevel: stderrLevel,
	}
}

// jsonHandler routes each record to one of two JSON handlers by level.
// Both handlers share the same attrs and groups, so the output shape
// doesn't depend on which stream a record lands on.
type jsonHandler struct {
	stdout      slog.Handler
	stderr      slog.Handler
	stderrLevel slog.Leveler
	// grouped is true once WithGroup has been called. The request ID is
	// then nested in the group, so it's no longer checked for duplicates.
	grouped bool
	// hasRequestID is true when request_id was added with WithAttrs.
	hasRequestID bool
}

// Enabled reports whether the handler will log at the given level.
func (h *jsonHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.stdout.Enabled(ctx, level)
}

// WithAttrs returns a new handler with additional attributes included in every log.
func (h *jsonHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	c := *h
	c.stdout = h.stdout.WithAttrs(attrs)
	c.stderr = h.stderr.WithAttrs(attrs)
	if !h.grouped && containsKey(attrs, requestIDKey) {
		c.hasRequestID = true
	}
	return &c
}

// WithGroup returns a new handler that qualifies subsequent attributes with the group name.
func (h *jsonHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	c := *h
	c.stdout = h.stdout.WithGroup(name)
	c.stderr = h.stderr.WithGroup(name)
	c.grouped = true
	return &c
}

// Handle writes the record to stdout or stderr depending on its level.
func (h *jsonHandler) Handle(ctx context.Context, r slog.Record) error {
	if reqID, ok := webkitctx.RequestID(ctx); ok && !h.hasRequestID && !recordHasKey(r, requestIDKey) {
		r = r.Clone()
		r.AddAttrs(slog.String(requestIDKey, reqID))
	}
	if r.Level >= h.stderrLevel.Level() {
		return h.stderr.Handle(ctx, r)
	}
	return h.stdout.Handle(ctx, r)
}

const requestIDKey = string(webkitctx.ContextKeyRequestID)

func containsKey(attrs []slog.Attr, key string) bool {
	for _, a := range attrs {
		if a.Key == key {
			return true
		}
	}
	return false
}

func recordHasKey(r slog.Record, key string) bool {
	found := false
	r.Attrs(func(a slog.Attr) bool {
		found = a.Key == key
		return !found
	})
	return found
}

// replaceLevelNames writes custom levels, such as LevelFatal, by name
// instead of as "ERROR+4", then calls next.
func replaceLevelNames(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if len(groups) == 0 && a.Key == slog.LevelKey {
			if lvl, ok := a.Value.Any().(slog.Level); ok && lvl == LevelFatal {
				a.Value = slog.StringValue(levelFatalName)
			}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}
