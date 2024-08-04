package payload

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ainsleyclark/go-payloadcms"
	payloadfakes "github.com/ainsleyclark/go-payloadcms/fakes"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestRedirects(t *testing.T) {
	t.Parallel()

	var (
		fromURL   = "/test"
		redirects = []redirect{
			{From: fromURL, To: "/new", Code: redirectsCode301},
		}
	)

	tt := map[string]struct {
		url        string
		mock       func(cols *payloadfakes.MockCollectionService, store cache.Store)
		wantURL    string
		wantStatus int
	}{
		"Skipped": {
			url:        "/favicon.ico",
			wantStatus: 200,
		},
		"API error returns nil": {
			url: fromURL,
			mock: func(cols *payloadfakes.MockCollectionService, store cache.Store) {
				cols.ListFunc = func(_ context.Context, _ payloadcms.Collection, _ payloadcms.ListParams, _ any) (payloadcms.Response, error) {
					return payloadcms.Response{}, errors.New("error")
				}
			},
			wantStatus: http.StatusOK,
		},
		"Invalid number defaults to 301": {
			url: fromURL,
			mock: func(_ *payloadfakes.MockCollectionService, store cache.Store) {
				store.Set(context.TODO(), redirectCacheKey, []redirect{
					{From: fromURL, To: "/new", Code: "wrong"},
				}, cache.Options{})
			},
			wantStatus: http.StatusMovedPermanently,
			wantURL:    "/new",
		},
		"No Matches": {
			url: fromURL,
			mock: func(_ *payloadfakes.MockCollectionService, store cache.Store) {
				store.Set(context.TODO(), redirectCacheKey, []redirect{
					{From: "/wrong", To: "/new", Code: redirectsCode301},
				}, cache.Options{})
			},
			wantStatus: http.StatusOK,
		},
		"Redirects 301 from API": {
			url: fromURL,
			mock: func(cols *payloadfakes.MockCollectionService, store cache.Store) {
				cols.ListFunc = func(_ context.Context, _ payloadcms.Collection, _ payloadcms.ListParams, out any) (payloadcms.Response, error) {
					*out.(*payloadcms.ListResponse[redirect]) = payloadcms.ListResponse[redirect]{
						Docs: redirects,
					}
					return payloadcms.Response{}, nil
				}
			},
			wantStatus: http.StatusMovedPermanently,
			wantURL:    "/new",
		},
		"Redirects 301 from Cache": {
			url: fromURL,
			mock: func(_ *payloadfakes.MockCollectionService, store cache.Store) {
				store.Set(context.TODO(), redirectCacheKey, redirects, cache.Options{})
			},
			wantStatus: http.StatusMovedPermanently,
			wantURL:    "/new",
		},
	}
	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			app := webkit.New()
			req := httptest.NewRequest(http.MethodGet, test.url, nil)
			rr := httptest.NewRecorder()

			store := cache.NewInMemory(time.Hour)
			collections := payloadfakes.NewMockCollectionService()
			payload := &payloadcms.Client{
				Collections: collections,
			}

			if test.mock != nil {
				test.mock(collections, store)
			}

			app.Plug(redirectMiddleware(payload, store))
			app.Get(test.url, func(c *webkit.Context) error {
				return c.String(http.StatusOK, "Middleware")
			})
			app.ServeHTTP(rr, req)

			assert.Equal(t, test.wantStatus, rr.Code)
			assert.Equal(t, test.wantURL, rr.Header().Get("Location"))
		})
	}
}
