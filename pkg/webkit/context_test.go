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
			require.NoError(t, c.Redirect(301, "/location"))
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

func TestContext_String(t *testing.T) {
	app := New()
	app.Get("/string", func(c *Context) error {
		require.NoError(t, c.String(http.StatusOK, "test"))
		return nil
	})
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, httptest.NewRequest("GET", "/string", nil))
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "test", rr.Body.String())
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

func TestContext_IsTLS(t *testing.T) {
	app := New()
	app.Get("/tls", func(ctx *Context) error {
		assert.False(t, ctx.IsTLS())
		return nil
	})
	app.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/tls", nil))
}

func TestContext_IsWebSocket(t *testing.T) {
	app := New()
	req := httptest.NewRequest("GET", "/websocket", nil)
	req.Header.Set("Upgrade", "websocket")
	app.Get("/websocket", func(ctx *Context) error {
		assert.True(t, ctx.IsWebSocket())
		return nil
	})
	app.ServeHTTP(httptest.NewRecorder(), req)
}

func TestContext_Scheme(t *testing.T) {
	tt := map[string]struct {
		req  func() *http.Request
		want string
	}{
		"HTTP request": {
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			},
			want: "http",
		},
		"HTTPS request": {
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "https://example.com", nil)
			},
			want: "https",
		},
		"Request with X-Forwarded-Proto header (https)": {
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
				req.Header.Set("X-Forwarded-Proto", "https")
				return req
			},
			want: "https",
		},
		//"Request with X-Forwarded-Proto header (http)": {
		//	req: func() *http.Request {
		//		req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
		//		req.Header.Set("X-Forwarded-Proto", "http")
		//		return req
		//	},
		//	want: "http",
		//},
		"Request with X-Forwarded-Protocol header": {
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
				req.Header.Set("X-Forwarded-Protocol", "https")
				return req
			},
			want: "https",
		},
		"Request with X-Forwarded-Ssl header": {
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
				req.Header.Set("X-Forwarded-Ssl", "on")
				return req
			},
			want: "https",
		},
		"Request with X-Url-Scheme header": {
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
				req.Header.Set("X-Url-Scheme", "https")
				return req
			},
			want: "https",
		},
		"Request without any X-Forwarded headers": {
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			},
			want: "http",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			app := New()
			app.Get("/", func(ctx *Context) error {
				assert.Equal(t, test.want, ctx.Scheme())
				return nil
			})
			app.ServeHTTP(httptest.NewRecorder(), test.req())
		})
	}
}
