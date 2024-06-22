package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestMinify(t *testing.T) {
	content := `<html>
<body>
    <h1>Welcome</h1>
    <p>This is a sample paragraph.</p>
</body>
</html>`

	app := webkit.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	app.Get("/", func(c *webkit.Context) error {
		return c.HTML(http.StatusOK, content)
	}, Minify)

	app.ServeHTTP(rr, req)

	got := rr.Body.String()
	want := `<html><body><h1>Welcome</h1><p>This is a sample paragraph.</p></body></html>`
	assert.Equal(t, want, got)
}
