package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	webkitctx "github.com/ainsleydev/webkit/pkg/context"
)

// fakeTransport records events sent to Sentry and how often it was flushed.
type fakeTransport struct {
	mu      sync.Mutex
	events  []*sentry.Event
	flushes int
}

func (t *fakeTransport) Configure(sentry.ClientOptions) {}
func (t *fakeTransport) Close()                         {}

func (t *fakeTransport) SendEvent(e *sentry.Event) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = append(t.events, e)
}

func (t *fakeTransport) Flush(time.Duration) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.flushes++
	return true
}

func (t *fakeTransport) FlushWithContext(context.Context) bool { return t.Flush(0) }

func (t *fakeTransport) Events() []*sentry.Event {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]*sentry.Event(nil), t.events...)
}

func newTestSentryClient(t *testing.T) (*sentry.Client, *fakeTransport) {
	t.Helper()
	tr := &fakeTransport{}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://public@example.com/1",
		Transport: tr,
	})
	require.NoError(t, err)
	return client, tr
}

func newTestSentryLogger(t *testing.T, opts *SentryHandlerOptions) (*slog.Logger, *bytes.Buffer, *fakeTransport) {
	t.Helper()
	client, tr := newTestSentryClient(t)
	buf := &bytes.Buffer{}
	next := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	return slog.New(NewSentryHandler(next, client, opts)), buf, tr
}

func TestSentryHandler_Levels(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		opts      *SentryHandlerOptions
		level     slog.Level
		wantEvent bool
		wantLevel sentry.Level
	}{
		"Info not captured":  {level: slog.LevelInfo},
		"Warn not captured":  {level: slog.LevelWarn},
		"Error captured":     {level: slog.LevelError, wantEvent: true, wantLevel: sentry.LevelError},
		"Fatal captured":     {level: LevelFatal, wantEvent: true, wantLevel: sentry.LevelFatal},
		"Above fatal":        {level: LevelFatal + 4, wantEvent: true, wantLevel: sentry.LevelFatal},
		"Fatal only - error": {opts: &SentryHandlerOptions{Level: LevelFatal}, level: slog.LevelError},
		"Warn configured": {
			opts:      &SentryHandlerOptions{Level: slog.LevelWarn},
			level:     slog.LevelWarn,
			wantEvent: true,
			wantLevel: sentry.LevelWarning,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			logger, buf, tr := newTestSentryLogger(t, test.opts)
			logger.Log(context.Background(), test.level, "msg")

			assert.NotEmpty(t, buf.String(), "every record should reach the next handler")

			events := tr.Events()
			if !test.wantEvent {
				assert.Empty(t, events)
				return
			}
			require.Len(t, events, 1)
			assert.Equal(t, test.wantLevel, events[0].Level)
			assert.Equal(t, "msg", events[0].Message)
		})
	}
}

func TestSentryHandler_Event(t *testing.T) {
	t.Parallel()

	logger, _, tr := newTestSentryLogger(t, nil)
	ctx := webkitctx.WithRequestID(context.Background(), "req-1")

	cause := errors.New("connection refused")
	logger.
		With("service", "api").
		WithGroup("db").
		With("host", "localhost").
		ErrorContext(ctx, "query failed",
			"error", fmt.Errorf("query: %w", cause),
			"other_err", errors.New("second"),
			slog.Group("query", slog.String("table", "users")),
			"rows", 0,
		)

	events := tr.Events()
	require.Len(t, events, 1)
	e := events[0]

	assert.Equal(t, "query failed", e.Message)
	assert.Equal(t, sentry.LevelError, e.Level)
	assert.Equal(t, "slog", e.Logger)
	assert.Equal(t, "req-1", e.Tags["request_id"])
	assert.False(t, e.Timestamp.IsZero())

	require.NotEmpty(t, e.Exception)
	values := make([]string, 0, len(e.Exception))
	for _, ex := range e.Exception {
		values = append(values, ex.Value)
	}
	assert.Contains(t, values, "query: connection refused")

	assert.Equal(t, map[string]any{
		"service":        "api",
		"db.host":        "localhost",
		"db.other_err":   "second",
		"db.query.table": "users",
		"db.rows":        int64(0),
	}, e.Extra)
}

func TestSentryHandler_NoError(t *testing.T) {
	t.Parallel()

	logger, _, tr := newTestSentryLogger(t, nil)
	logger.Error("plain", "k", "v")

	events := tr.Events()
	require.Len(t, events, 1)
	assert.Empty(t, events[0].Exception)
	assert.Equal(t, map[string]any{"k": "v"}, events[0].Extra)
}

func TestSentryHandler_Flush(t *testing.T) {
	t.Parallel()

	t.Run("Default flushes", func(t *testing.T) {
		t.Parallel()
		logger, _, tr := newTestSentryLogger(t, nil)
		logger.Error("msg")
		assert.Equal(t, 1, tr.flushes)
	})

	t.Run("Negative disables", func(t *testing.T) {
		t.Parallel()
		logger, _, tr := newTestSentryLogger(t, &SentryHandlerOptions{FlushTimeout: -1})
		logger.Error("msg")
		assert.Len(t, tr.Events(), 1)
		assert.Equal(t, 0, tr.flushes)
	})
}

func TestSentryHandler_HubFromContext(t *testing.T) {
	t.Parallel()

	client, tr := newTestSentryClient(t)
	hub := sentry.NewHub(client, sentry.NewScope())
	hub.Scope().SetTag("route", "/users")
	ctx := sentry.SetHubOnContext(context.Background(), hub)

	buf := &bytes.Buffer{}
	logger := slog.New(NewSentryHandler(slog.NewJSONHandler(buf, nil), nil, nil))
	logger.ErrorContext(ctx, "msg")

	events := tr.Events()
	require.Len(t, events, 1)
	assert.Equal(t, "/users", events[0].Tags["route"], "scope from the context hub should be applied")
}

func TestSentryHandler_Enabled(t *testing.T) {
	t.Parallel()

	client, _ := newTestSentryClient(t)
	next := slog.NewJSONHandler(&bytes.Buffer{}, &slog.HandlerOptions{Level: LevelFatal})
	h := NewSentryHandler(next, client, nil)

	assert.False(t, h.Enabled(context.Background(), slog.LevelWarn))
	assert.True(t, h.Enabled(context.Background(), slog.LevelError), "Sentry wants Error even if next doesn't")
}

func TestSentryLevel(t *testing.T) {
	t.Parallel()

	tt := map[slog.Level]sentry.Level{
		slog.LevelDebug - 4: sentry.LevelDebug,
		slog.LevelDebug:     sentry.LevelDebug,
		slog.LevelInfo:      sentry.LevelInfo,
		slog.LevelWarn:      sentry.LevelWarning,
		slog.LevelError:     sentry.LevelError,
		LevelFatal:          sentry.LevelFatal,
	}
	for in, want := range tt {
		assert.Equal(t, want, sentryLevel(in), in.String())
	}
}
