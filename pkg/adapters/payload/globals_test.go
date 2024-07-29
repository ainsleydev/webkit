package payload

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ainsleyclark/go-payloadcms"
	payloadfakes "github.com/ainsleyclark/go-payloadcms/fakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestGlobalsMiddleware(t *testing.T) {
	t.Parallel()

	GlobalMiddlewareTestHelper(t, func(client *payloadcms.Client, store cache.Store) webkit.Plug {
		return globalsMiddleware[Settings](client, store, "settings")
	})
}

func GlobalMiddlewareTestHelper(t *testing.T, fn func(client *payloadcms.Client, store cache.Store) webkit.Plug) {
	t.Helper()

	err := os.Setenv(env.AppEnvironmentKey, env.Production)
	require.NoError(t, err)

	settings := Settings{
		SiteName: ptr.StringPtr("Site Name"),
	}

	tt := map[string]struct {
		url  string
		mock func(gb *payloadfakes.MockGlobalsService, store cache.Store)
		want any
	}{
		"File Request": {
			url:  "/favicon.ico",
			mock: func(gb *payloadfakes.MockGlobalsService, store cache.Store) {},
			want: nil,
		},
		"From Cache": {
			url: "/want",
			mock: func(gb *payloadfakes.MockGlobalsService, store cache.Store) {
				store.Set(context.TODO(), GlobalsContextKey("settings"), &settings, cache.Options{})
			},
			want: &settings,
		},
		"API Error": {
			url: "/want",
			mock: func(gb *payloadfakes.MockGlobalsService, store cache.Store) {
				gb.GetFunc = func(_ context.Context, _ payloadcms.Global, out any) (payloadcms.Response, error) {
					return payloadcms.Response{}, assert.AnError
				}
			},
			want: nil,
		},
		"From API": {
			url: "/want",
			mock: func(gb *payloadfakes.MockGlobalsService, store cache.Store) {
				gb.GetFunc = func(_ context.Context, _ payloadcms.Global, out any) (payloadcms.Response, error) {
					*out.(*Settings) = settings
					return payloadcms.Response{}, nil
				}
			},
			want: &settings,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			app := webkit.New()
			req := httptest.NewRequest(http.MethodGet, test.url, nil)
			rr := httptest.NewRecorder()

			store := cache.NewInMemory(time.Hour)
			globals := payloadfakes.NewMockGlobalsService()
			payload := &payloadcms.Client{
				Globals: globals,
			}

			test.mock(globals, store)

			app.Plug(fn(payload, store))
			app.Get(test.url, func(c *webkit.Context) error {
				assert.Equal(t, test.want, c.Get(SettingsContextKey))
				return nil
			})
			app.ServeHTTP(rr, req)
		})
	}
}
