package payload

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestAdapter_Robots(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		ctx  func(c *webkit.Context)
		env  string
		want string
	}{
		"Nil robots Production": {
			ctx:  func(_ *webkit.Context) {},
			env:  env.Production,
			want: "User-agent: *\nDisallow:",
		},
		"Nil robots Dev": {
			ctx:  func(_ *webkit.Context) {},
			env:  env.Development,
			want: "User-agent: *\nDisallow: /",
		},
		"Digital Ocean": {
			ctx: func(c *webkit.Context) {
				u, err := url.Parse("https://s-clark-cms-ejioo.ondigitalocean.app/")
				require.NoError(t, err)
				c.Request.URL = u
			},
			env:  env.Production,
			want: "User-agent: *\nDisallow: /",
		},
		"Settings": {
			ctx: func(c *webkit.Context) {
				c.Set(SettingsContextKey, &Settings{
					Robots: ptr.StringPtr("Custom"),
				})
			},
			env:  env.Development,
			want: "Custom",
		},
		"Settings No robots Set": {
			ctx: func(c *webkit.Context) {
				c.Set(SettingsContextKey, &Settings{})
			},
			env:  env.Development,
			want: "User-agent: *\nDisallow: /",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			app := webkit.New()
			app.Plug(func(next webkit.Handler) webkit.Handler {
				return func(c *webkit.Context) error {
					test.ctx(c)
					return next(c)
				}
			})
			app.Get("/robots.txt", robots(test.env))
			rr := httptest.NewRecorder()
			app.ServeHTTP(rr, httptest.NewRequest("GET", "/robots.txt", nil))

			assert.Equal(t, test.want, rr.Body.String())
			assert.Equal(t, 200, rr.Code)
		})
	}
}
