package webkit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/url", nil)
	ctx := NewContext(httptest.NewRecorder(), req)
	assert.NotEmpty(t, ctx.Response)
	assert.Equal(t, req, ctx.Request)
}

func TestContext_GetSet(t *testing.T) {
	c := NewContext(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/url", nil))
	c.Set("test", "value")
	got := c.Get("test")
	assert.Equal(t, got, "value")
}

func TestContext_Context(t *testing.T) {
	c := NewContext(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/url", nil))
	got := c.Context()
	assert.NotNil(t, got)
	assert.IsType(t, got, context.Background())
}

func TestContext_Param(t *testing.T) {
	app := New()
	app.Get("/users/{id}", func(c *Context) error {
		got := c.Param("id")
		assert.Equal(t, "123", got)
		return nil
	})
	app.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/users/123", nil))
}

func TestContext_Render(t *testing.T) {
	t.Run("Redirects", func(t *testing.T) {
		tt := map[int]string{
			http.StatusMovedPermanently:  "/newlocation",
			http.StatusFound:             "/foundlocation",
			http.StatusSeeOther:          "/seelocation",
			http.StatusTemporaryRedirect: "/temporaryredirectlocation",
			http.StatusPermanentRedirect: "/permanentredirectlocation",
		}

		for code, location := range tt {
			app := New()

			app.Get("/redirect", func(c *Context) error {
				require.NoError(t, c.Redirect(code, location))
				return nil
			})

			rr := httptest.NewRecorder()

			app.ServeHTTP(rr, httptest.NewRequest("GET", "/redirect", nil))

			assert.Equal(t, code, rr.Code)
			assert.Equal(t, location, rr.Header().Get("Location"))
		}
	})

	t.Run("Errors", func(t *testing.T) {
		app := New()
		app.Get("/redirect", func(c *Context) error {
			require.NoError(t, c.Redirect(500, "/location"))
			return nil
		})
		app.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/redirect", nil))
	})
}

func TestContext_NoContext(t *testing.T) {
	app := New()
	app.Get("/nocontext", func(c *Context) error {
		require.NoError(t, c.NoContext(http.StatusOK))
		return nil
	})
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, httptest.NewRequest("GET", "/nocontext", nil))
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestContext_Get(t *testing.T) {

}

func TestContext_JSON(t *testing.T) {
	app := New()

	app.Get("/json", func(c *Context) error {
		require.NoError(t, c.JSON(http.StatusOK, map[string]any{"test": 1}))
		return nil
	})

	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/json", nil))

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"test": 1}`, strings.TrimSpace(rr.Body.String()))
}
