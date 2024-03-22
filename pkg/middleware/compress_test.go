package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestGzipHandler(t *testing.T) {
	app := webkit.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	req.Header.Set("Accept-Encoding", "gzip")

	app.Get("/", func(ctx *webkit.Context) error {
		bs := make([]byte, 1000) // Ensure it's over 512.
		require.NoError(t, ctx.String(http.StatusOK, string(bs)))
		return nil
	}, Gzip)

	app.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))
	assert.Equal(t, "Accept-Encoding", rr.Header().Get("Vary"))
}
