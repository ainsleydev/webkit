package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	webkitctx "github.com/ainsleydev/webkit/pkg/context"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestRequestID(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		app := webkit.New()

		app.Get("/", func(ctx *webkit.Context) error {
			rid, ok := webkitctx.RequestID(ctx.Request.Context())
			assert.True(t, ok)
			assert.NotEmpty(t, rid)
			return ctx.String(http.StatusOK, "test")
		}, RequestID)

		app.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	})

	t.Run("Existing", func(t *testing.T) {
		app := webkit.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(RequestIDHeader, "123")

		app.Get("/", func(ctx *webkit.Context) error {
			rid, ok := webkitctx.RequestID(ctx.Request.Context())
			assert.True(t, ok)
			assert.Equal(t, ctx.Request.Header.Get(RequestIDHeader), "123")
			assert.Equal(t, ctx.Request.Header.Get(RequestIDHeader), rid)
			return ctx.String(http.StatusOK, "test")
		}, RequestID)

		app.ServeHTTP(httptest.NewRecorder(), req)
	})
}
