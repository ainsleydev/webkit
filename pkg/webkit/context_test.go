package webkit

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
				err := c.Redirect(code, location)
				assert.NoError(t, err)
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
			err := c.Redirect(500, "/location")
			assert.Error(t, err)
			return nil
		})
		app.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/redirect", nil))
	})
}

func TestContext_JSON(t *testing.T) {
	app := New()

	app.Get("/json", func(c *Context) error {
		err := c.JSON(http.StatusOK, map[string]any{"test": 1})
		assert.NoError(t, err)
		return nil
	})

	rr := httptest.NewRecorder()

	app.ServeHTTP(rr, httptest.NewRequest("GET", "/json", nil))

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"test": 1}`, strings.TrimSpace(rr.Body.String()))
}