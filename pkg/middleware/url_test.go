package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestURL(t *testing.T) {
	app := webkit.New()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.URL.Host = "example.com"

	app.Get("/", func(ctx *webkit.Context) error {
		url := ctx.Get(URLContextKey)
		assert.Equal(t, "http://example.com/", url)
		return nil
	}, URL)

	app.ServeHTTP(rr, req)
}
