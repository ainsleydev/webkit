package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	webkitctx "github.com/ainsleydev/webkit/pkg/context"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestURL(t *testing.T) {
	tt := map[string]struct {
		mock func(r *http.Request)
		path string
		want string
	}{
		"From Environment Homepage": {
			mock: func(_ *http.Request) {
				t.Setenv("APP_URL", "http://example.com")
			},
			path: "/",
			want: "http://example.com",
		},
		"From Environment Path": {
			mock: func(_ *http.Request) {
				t.Setenv("APP_URL", "http://example.com")
			},
			path: "/test",
			want: "http://example.com/test",
		},
		"From Request Homepage": {
			mock: func(r *http.Request) {
				r.Host = "example.com"
			},
			path: "/",
			want: "http://example.com",
		},
		"From Request Path": {
			mock: func(r *http.Request) {
				r.Host = "example.com"
			},
			path: "/test",
			want: "http://example.com/test",
		},
		"Localhost": {
			mock: func(r *http.Request) {
				r.Host = ""
			},
			path: "/test",
			want: "http://localhost/test",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			defer func() {
				require.NoError(t, os.Unsetenv("APP_URL"))
			}()

			app := webkit.New()
			app.Get(test.path, func(c *webkit.Context) error {
				got, ok := webkitctx.URL(c.Context())
				assert.True(t, ok)
				assert.Equal(t, test.want, got)
				return nil
			}, URL)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, test.path, nil)
			test.mock(req)

			app.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code)
		})
	}
}
