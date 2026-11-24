package log

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	webkitctx "github.com/ainsleydev/webkit/pkg/context"
	"github.com/ainsleydev/webkit/pkg/env"
)

func newTestJSONHandler(opts *JSONHandlerOptions) (slog.Handler, *bytes.Buffer, *bytes.Buffer) {
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	return NewJSONHandler(stdout, stderr, opts), stdout, stderr
}

// decodeLine decodes a single JSON log line and fails if the buffer
// holds anything other than exactly one line.
func decodeLine(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	require.Len(t, lines, 1, "expected exactly one log line, got: %q", buf.String())
	var out map[string]any
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &out))
	return out
}

// countKey counts occurrences of "key": in raw JSON, which catches
// duplicate keys that json.Unmarshal silently collapses.
func countKey(s, key string) int {
	return strings.Count(s, `"`+key+`":`)
}

func TestJSONHandler_Shape(t *testing.T) {
	t.Parallel()

	h, stdout, _ := newTestJSONHandler(nil)
	logger := slog.New(h).With("service", "api")

	ctx := webkitctx.WithRequestID(context.Background(), "req-123")
	logger.InfoContext(ctx, "hello",
		slog.String("url", "/does-not-exist"),
		slog.Int("status", 404),
		slog.Group("user", slog.Int("id", 1)),
	)

	raw := stdout.String()
	got := decodeLine(t, stdout)

	assert.Equal(t, "INFO", got["level"])
	assert.Equal(t, "hello", got["msg"])
	assert.Equal(t, "/does-not-exist", got["url"])
	assert.EqualValues(t, 404, got["status"])
	assert.Equal(t, "api", got["service"], "With attrs should be top level")
	assert.Equal(t, "req-123", got["request_id"])
	assert.Equal(t, map[string]any{"id": float64(1)}, got["user"])
	assert.NotContains(t, got, "attr", "attrs should not be repeated in an attr group")

	for _, key := range []string{"url", "status", "service", "request_id", "user"} {
		assert.Equal(t, 1, countKey(raw, key), "%s should appear exactly once in %s", key, raw)
	}
}

func TestJSONHandler_RequestID(t *testing.T) {
	t.Parallel()

	ctx := webkitctx.WithRequestID(context.Background(), "from-ctx")

	tt := map[string]struct {
		log  func(l *slog.Logger)
		want string
	}{
		"From context": {
			log:  func(l *slog.Logger) { l.InfoContext(ctx, "msg") },
			want: "from-ctx",
		},
		"Record attr wins": {
			log:  func(l *slog.Logger) { l.InfoContext(ctx, "msg", "request_id", "from-attr") },
			want: "from-attr",
		},
		"With attr wins": {
			log:  func(l *slog.Logger) { l.With("request_id", "from-with").InfoContext(ctx, "msg") },
			want: "from-with",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			h, stdout, _ := newTestJSONHandler(nil)
			test.log(slog.New(h))
			raw := stdout.String()
			assert.Equal(t, 1, countKey(raw, "request_id"), raw)
			assert.Equal(t, test.want, decodeLine(t, stdout)["request_id"])
		})
	}

	t.Run("No request ID", func(t *testing.T) {
		t.Parallel()
		h, stdout, _ := newTestJSONHandler(nil)
		slog.New(h).Info("msg")
		assert.NotContains(t, decodeLine(t, stdout), "request_id")
	})
}

func TestJSONHandler_StreamRouting(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		stderrLevel slog.Leveler
		level       slog.Level
		wantStderr  bool
		wantLevel   string
	}{
		"Debug to stdout":           {level: slog.LevelDebug, wantStderr: false, wantLevel: "DEBUG"},
		"Info to stdout":            {level: slog.LevelInfo, wantStderr: false, wantLevel: "INFO"},
		"Warn to stdout by default": {level: slog.LevelWarn, wantStderr: false, wantLevel: "WARN"},
		"Error to stderr":           {level: slog.LevelError, wantStderr: true, wantLevel: "ERROR"},
		"Fatal to stderr":           {level: LevelFatal, wantStderr: true, wantLevel: "FATAL"},
		"Warn to stderr when configured": {
			stderrLevel: slog.LevelWarn,
			level:       slog.LevelWarn,
			wantStderr:  true,
			wantLevel:   "WARN",
		},
		"Info to stdout when Warn configured": {
			stderrLevel: slog.LevelWarn,
			level:       slog.LevelInfo,
			wantStderr:  false,
			wantLevel:   "INFO",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			h, stdout, stderr := newTestJSONHandler(&JSONHandlerOptions{
				HandlerOptions: slog.HandlerOptions{Level: slog.LevelDebug},
				StderrLevel:    test.stderrLevel,
			})
			slog.New(h).With("k", "v").Log(context.Background(), test.level, "msg")

			written, empty := stdout, stderr
			if test.wantStderr {
				written, empty = stderr, stdout
			}
			assert.Empty(t, empty.String())
			got := decodeLine(t, written)
			assert.Equal(t, test.wantLevel, got["level"])
			assert.Equal(t, "v", got["k"], "With attrs should reach both streams")
		})
	}
}

func TestJSONHandler_Groups(t *testing.T) {
	t.Parallel()

	h, _, stderr := newTestJSONHandler(nil)
	slog.New(h).WithGroup("req").With("id", 1).Error("failed", "error", errors.New("boom"))

	got := decodeLine(t, stderr)
	assert.Equal(t, map[string]any{"id": float64(1), "error": "boom"}, got["req"])
}

func TestJSONHandler_Enabled(t *testing.T) {
	t.Parallel()

	h, _, _ := newTestJSONHandler(&JSONHandlerOptions{
		HandlerOptions: slog.HandlerOptions{Level: slog.LevelWarn},
	})
	assert.False(t, h.Enabled(context.Background(), slog.LevelInfo))
	assert.True(t, h.Enabled(context.Background(), slog.LevelWarn))
}

func TestJSONHandler_ReplaceAttrChained(t *testing.T) {
	t.Parallel()

	h, _, stderr := newTestJSONHandler(&JSONHandlerOptions{
		HandlerOptions: slog.HandlerOptions{
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if a.Key == "password" {
					return slog.String("password", "***")
				}
				return a
			},
		},
	})
	slog.New(h).Log(context.Background(), LevelFatal, "msg", "password", "secret")

	got := decodeLine(t, stderr)
	assert.Equal(t, "FATAL", got["level"])
	assert.Equal(t, "***", got["password"])
}

func TestResolveLogHandler_Production(t *testing.T) {
	t.Setenv(env.AppEnvironmentKey, env.Production.String())

	_, ok := NewHandler("test").(*jsonHandler)
	assert.True(t, ok, "production should use the JSON handler")
}
