package log

import (
	"context"
	"log/slog"
	"time"

	"github.com/getsentry/sentry-go"

	webkitctx "github.com/ainsleydev/webkit/pkg/context"
)

// SentryHandlerOptions configures a SentryHandler.
type SentryHandlerOptions struct {
	// Level is the minimum level reported to Sentry.
	// Defaults to slog.LevelError.
	Level slog.Leveler

	// FlushTimeout is how long Handle waits for Sentry to send each
	// event. Serverless platforms can freeze the process once the
	// response is written, so unflushed events may be lost.
	// Defaults to 2 seconds. A negative value turns flushing off.
	FlushTimeout time.Duration
}

const defaultSentryFlushTimeout = 2 * time.Second

// SentryHandler is a slog.Handler that passes every record to the next
// handler and also reports records at or above a level to Sentry.
//
// The first error-valued attribute becomes the event's exception, with
// its stack trace when the error carries one. All other attributes,
// including those added through With/WithAttrs, go in event.Extra, with
// group names joined by dots. A request ID in the context is added as
// the request_id tag.
//
// It is the slog replacement for the deprecated logrus-based SentryHook.
type SentryHandler struct {
	next         slog.Handler
	client       *sentry.Client
	level        slog.Leveler
	flushTimeout time.Duration
	attrs        []slog.Attr // Attrs from WithAttrs, already qualified by groups.
	groups       []string
}

// NewSentryHandler wraps next so that records at or above opts.Level are
// also sent to Sentry.
//
// The event is sent with client. If client is nil, the client of the hub
// in the record's context is used, falling back to sentry.CurrentHub().
// When the context holds a hub, its scope (tags, user, request) is
// applied to the event.
func NewSentryHandler(next slog.Handler, client *sentry.Client, opts *SentryHandlerOptions) *SentryHandler {
	if opts == nil {
		opts = &SentryHandlerOptions{}
	}
	level := opts.Level
	if level == nil {
		level = slog.LevelError
	}
	flush := opts.FlushTimeout
	if flush == 0 {
		flush = defaultSentryFlushTimeout
	}
	return &SentryHandler{
		next:         next,
		client:       client,
		level:        level,
		flushTimeout: flush,
	}
}

// Enabled reports whether either the next handler or Sentry wants the level.
func (h *SentryHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level.Level() || h.next.Enabled(ctx, level)
}

// WithAttrs returns a new handler with additional attributes included in every log.
func (h *SentryHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	c := *h
	c.next = h.next.WithAttrs(attrs)
	c.attrs = make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	c.attrs = append(c.attrs, h.attrs...)
	for _, a := range attrs {
		c.attrs = append(c.attrs, qualify(h.groups, a))
	}
	return &c
}

// WithGroup returns a new handler that qualifies subsequent attributes with the group name.
func (h *SentryHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	c := *h
	c.next = h.next.WithGroup(name)
	c.groups = append(append([]string{}, h.groups...), name)
	return &c
}

// Handle passes the record to the next handler and, when the record's
// level is high enough, reports it to Sentry.
func (h *SentryHandler) Handle(ctx context.Context, r slog.Record) error {
	var err error
	if h.next.Enabled(ctx, r.Level) {
		err = h.next.Handle(ctx, r.Clone())
	}
	if r.Level >= h.level.Level() {
		h.capture(ctx, r)
	}
	return err
}

func (h *SentryHandler) capture(ctx context.Context, r slog.Record) {
	client, scope := h.resolveClient(ctx)
	if client == nil {
		return
	}

	event := h.buildEvent(ctx, r, client.Options().MaxErrorDepth)
	client.CaptureEvent(event, nil, scope)

	if h.flushTimeout > 0 {
		client.Flush(h.flushTimeout)
	}
}

// resolveClient picks the client to send with and the scope to apply.
func (h *SentryHandler) resolveClient(ctx context.Context) (*sentry.Client, sentry.EventModifier) {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil && h.client == nil {
		hub = sentry.CurrentHub()
	}

	client := h.client
	var scope sentry.EventModifier = sentryModifier
	if hub != nil {
		if client == nil {
			client = hub.Client()
		}
		if s := hub.Scope(); s != nil {
			scope = s
		}
	}
	return client, scope
}

func (h *SentryHandler) buildEvent(ctx context.Context, r slog.Record, maxErrorDepth int) *sentry.Event {
	event := sentry.NewEvent()
	event.Message = r.Message
	event.Level = sentryLevel(r.Level)
	event.Logger = "slog"
	event.Timestamp = r.Time
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	if reqID, ok := webkitctx.RequestID(ctx); ok {
		event.Tags[requestIDKey] = reqID
	}

	var captured error
	add := func(a slog.Attr) {
		flattenAttr(a, "", func(key string, v slog.Value) {
			if e, ok := v.Any().(error); ok {
				if captured == nil {
					captured = e
					return
				}
				event.Extra[key] = e.Error()
				return
			}
			event.Extra[key] = v.Any()
		})
	}

	for _, a := range h.attrs {
		add(a)
	}
	r.Attrs(func(a slog.Attr) bool {
		add(qualify(h.groups, a))
		return true
	})

	if captured != nil {
		event.SetException(captured, maxErrorDepth)
	}
	return event
}

// sentryLevel maps a slog level to the closest Sentry level.
func sentryLevel(l slog.Level) sentry.Level {
	switch {
	case l >= LevelFatal:
		return sentry.LevelFatal
	case l >= slog.LevelError:
		return sentry.LevelError
	case l >= slog.LevelWarn:
		return sentry.LevelWarning
	case l >= slog.LevelInfo:
		return sentry.LevelInfo
	default:
		return sentry.LevelDebug
	}
}

// qualify nests a inside the given groups, outermost first.
func qualify(groups []string, a slog.Attr) slog.Attr {
	for i := len(groups) - 1; i >= 0; i-- {
		a = slog.Attr{Key: groups[i], Value: slog.GroupValue(a)}
	}
	return a
}

// flattenAttr resolves a and calls fn for every leaf value, with group
// keys joined by dots. Empty attrs are skipped, and the attrs of a group
// with an empty key are inlined, as slog handlers do.
func flattenAttr(a slog.Attr, prefix string, fn func(key string, v slog.Value)) {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return
	}
	key := a.Key
	if prefix != "" && key != "" {
		key = prefix + "." + key
	} else if key == "" {
		key = prefix
	}
	if a.Value.Kind() == slog.KindGroup {
		for _, ga := range a.Value.Group() {
			flattenAttr(ga, key, fn)
		}
		return
	}
	if key == "" {
		return
	}
	fn(key, a.Value)
}

// SentryEventIdentityModifier passes Sentry events through without modification.
type SentryEventIdentityModifier struct{}

// ApplyToEvent returns the event unchanged.
func (m *SentryEventIdentityModifier) ApplyToEvent(event *sentry.Event, _ *sentry.EventHint, _ *sentry.Client) *sentry.Event {
	return event
}

var sentryModifier = &SentryEventIdentityModifier{}
