package payload

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ainsleyclark/go-payloadcms"
	payloadfakes "github.com/ainsleyclark/go-payloadcms/fakes"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestSettingsMiddleware(t *testing.T) {
	t.Parallel()

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
				store.Set(context.TODO(), settingsCacheKey, settings, cache.Options{})
			},
			want: settings,
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
			want: settings,
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

			app.Plug(SettingsMiddleware(payload, store))
			app.Get(test.url, func(c *webkit.Context) error {
				assert.Equal(t, test.want, c.Get(SettingsContextKey))
				return nil
			})
			app.ServeHTTP(rr, req)
		})
	}
}
