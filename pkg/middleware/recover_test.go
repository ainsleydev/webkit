package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestRecover(t *testing.T) {
	app := webkit.New()
	rr := httptest.NewRecorder()

	app.Get("/", func(_ *webkit.Context) error {
		panic("test")
	}, Recover)

	assert.Panics(t, func() {
		app.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
