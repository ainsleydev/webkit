package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/log"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestStatusLevel(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		status int
		want   slog.Level
	}{
		"Unset":        {status: 0, want: slog.LevelWarn},
		"OK":           {status: http.StatusOK, want: slog.LevelInfo},
		"Redirect":     {status: http.StatusFound, want: slog.LevelInfo},
		"Not found":    {status: http.StatusNotFound, want: slog.LevelWarn},
		"Bad request":  {status: http.StatusBadRequest, want: slog.LevelWarn},
		"Server error": {status: http.StatusInternalServerError, want: slog.LevelError},
		"Bad gateway":  {status: http.StatusBadGateway, want: slog.LevelError},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, statusLevel(test.status))
		})
	}
}

func TestLatencyMillis(t *testing.T) {
	t.Parallel()

	assert.InDelta(t, 10.444, latencyMillis(10444*time.Microsecond), 0.0001)
	assert.InDelta(t, 424.496, latencyMillis(424496973*time.Nanosecond), 0.0001)
	assert.Equal(t, 0.0, latencyMillis(0))
}

func TestLogger_Production(t *testing.T) {
	t.Setenv(env.AppEnvironmentKey, env.Production.String())

	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	prev := slog.Default()
	slog.SetDefault(slog.New(log.NewJSONHandler(stdout, stderr, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })

	app := webkit.New()
	app.Plug(RequestID, Logger)
	app.Get("/missing", func(ctx *webkit.Context) error {
		return ctx.String(http.StatusNotFound, "Not Found")
	})

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	req.Header.Set(RequestIDHeader, "req-42")
	app.ServeHTTP(httptest.NewRecorder(), req)

	assert.Empty(t, stderr.String(), "4xx should not be written to stderr")

	line := strings.TrimSpace(stdout.String())
	var got map[string]any
	require.NoError(t, json.Unmarshal([]byte(line), &got), line)

	assert.Equal(t, "WARN", got["level"])
	assert.EqualValues(t, http.StatusNotFound, got["status"])
	assert.Equal(t, "/missing", got["url"])
	assert.Equal(t, "req-42", got["request_id"])
	assert.IsType(t, float64(0), got["latency_ms"])
	assert.NotContains(t, got, "latency", "raw nanosecond latency should be gone")
	assert.NotContains(t, got, "attr")
	assert.NotContains(t, got, "error", "error should be omitted when there is none")
	assert.NotContains(t, got, "cache", "cache should be omitted when the header is empty")
	assert.Equal(t, 1, strings.Count(line, `"status":`), line)
}
