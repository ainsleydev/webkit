package payload

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestMaintenanceMiddleware(t *testing.T) {
	t.Parallel()

	maintenance := &Maintenance{
		Enabled: true,
		Title:   "Title",
		Content: "Content",
	}

	tt := map[string]struct {
		url    string
		ctx    func() context.Context
		called bool
	}{
		"Skipped": {
			url: "/favicon.ico",
			ctx: func() context.Context {
				return context.TODO()
			},
			called: false,
		},
		"No Settings": {
			url: "/page",
			ctx: func() context.Context {
				return context.TODO()
			},
			called: false,
		},
		"Nil Maintenance": {
			url: "/page",
			ctx: func() context.Context {
				return context.WithValue(context.TODO(), SettingsContextKey, &Settings{}) //nolint
			},
			called: false,
		},
		"Not Enabled": {
			url: "/page",
			ctx: func() context.Context {
				return context.WithValue(context.TODO(), SettingsContextKey, &Settings{
					Maintenance: &Maintenance{
						Enabled: false,
					},
				})
			},
			called: false,
		},
		"Enabled": {
			url: "/page",
			ctx: func() context.Context {
				return context.WithValue(context.TODO(), SettingsContextKey, &Settings{
					Maintenance: maintenance,
				})
			},
			called: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			app := webkit.New()
			req := httptest.NewRequest(http.MethodGet, test.url, nil)
			req = req.WithContext(test.ctx())
			rr := httptest.NewRecorder()

			called := false
			handler := func(_ *webkit.Context, m Maintenance) error {
				called = true
				assert.EqualValues(t, maintenance, &m)
				return nil
			}

			app.Plug(maintenanceMiddleware(handler))
			app.Get(test.url, func(c *webkit.Context) error {
				return c.String(http.StatusOK, "OK")
			})
			app.ServeHTTP(rr, req)

			assert.Equal(t, test.called, called)
		})
	}
}
