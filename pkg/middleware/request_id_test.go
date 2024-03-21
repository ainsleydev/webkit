package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestRequestID(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		app := webkit.New()
		rr := httptest.NewRecorder()
		app.Plug(RequestID)
		app.Get("/", func(ctx *webkit.Context) error {
			assert.NotEmpty(t, ctx.Get(RequestIDContextKey))
			return ctx.String(http.StatusOK, "test")
		})
		app.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	})

	t.Run("Existing", func(t *testing.T) {
		app := webkit.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		req.Header.Set(RequestIDHeader, "123")
		app.Plug(RequestID)
		app.Get("/", func(ctx *webkit.Context) error {
			assert.NotEmpty(t, ctx.Get(RequestIDContextKey))
			assert.Equal(t, ctx.Request.Header.Get(RequestIDHeader), "123")
			assert.Equal(t, ctx.Request.Header.Get(RequestIDHeader), ctx.Get(RequestIDContextKey))
			return ctx.String(http.StatusOK, "test")
		})
		app.ServeHTTP(rr, req)
	})
}
