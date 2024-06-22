package payload

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestCacheMiddleware(t *testing.T) {
	t.Parallel()

	defaultHandler := func(c *webkit.Context) error {
		return c.String(http.StatusOK, "Cache")
	}

	tt := map[string]struct {
		url      string
		method   string
		mock     func(store cache.Store)
		handler  webkit.Handler
		assertFn func(rr *httptest.ResponseRecorder, store cache.Store)
	}{
		"Skipped": {
			url:    "/favicon.ico",
			method: http.MethodGet,
			assertFn: func(rr *httptest.ResponseRecorder, store cache.Store) {
				assert.Equal(t, "", rr.Header().Get("X-Cache"))
			},
		},
		"From Cache": {
			url:    "/page",
			method: http.MethodGet,
			mock: func(store cache.Store) {
				store.Set(context.TODO(), "page:/page", "Cache", cache.Options{})
			},
			assertFn: func(rr *httptest.ResponseRecorder, store cache.Store) {
				assert.Equal(t, "HIT", rr.Header().Get("X-Cache"))
				assert.Equal(t, "public, max-age=2419200", rr.Header().Get("Cache-Control"))
				var p string
				assert.NoError(t, store.Get(context.TODO(), "page:/page", &p))
			},
		},
		"Cache Miss": {
			url:    "/page",
			method: http.MethodGet,
			assertFn: func(rr *httptest.ResponseRecorder, store cache.Store) {
				assert.Equal(t, "MISS", rr.Header().Get("X-Cache"))
			},
		},
		"Next Error": {
			url:    "/page",
			method: http.MethodGet,
			handler: func(c *webkit.Context) error {
				return errors.New("next error")
			},
			assertFn: func(rr *httptest.ResponseRecorder, store cache.Store) {
				assert.Equal(t, "MISS", rr.Header().Get("X-Cache"))
				var p string
				assert.Error(t, store.Get(context.TODO(), "page:/page", &p))
			},
		},
		"Non 2xx": {
			url:    "/page",
			method: http.MethodGet,
			handler: func(c *webkit.Context) error {
				return c.String(http.StatusNotFound, "Not Found")
			},
			assertFn: func(rr *httptest.ResponseRecorder, store cache.Store) {
				assert.Equal(t, "MISS", rr.Header().Get("X-Cache"))
				var p string
				assert.Error(t, store.Get(context.TODO(), "page:/page", &p))
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			app := webkit.New()
			req := httptest.NewRequest(test.method, test.url, nil)
			rr := httptest.NewRecorder()

			store := cache.NewInMemory(time.Hour)

			if test.mock != nil {
				test.mock(store)
			}

			handler := defaultHandler
			if test.handler != nil {
				handler = test.handler
			}

			app.Plug(cacheMiddleware(store))
			app.Add(test.method, test.url, handler)
			app.ServeHTTP(rr, req)

			test.assertFn(rr, store)
		})
	}

	t.Run("Cache Miss", func(t *testing.T) {
		app := webkit.New()
		req := httptest.NewRequest(http.MethodGet, "/page", nil)
		rr := httptest.NewRecorder()

		store := cache.NewInMemory(time.Hour)

		app.Plug(cacheMiddleware(store))
		app.Get("/page", func(c *webkit.Context) error {
			return c.String(http.StatusOK, "Cache")
		})
		app.ServeHTTP(rr, req)

		assert.Equal(t, "MISS", rr.Header().Get("X-Cache"))
	})
}

func TestCacheBust(t *testing.T) {
	t.Parallel()
	t.Skip()

	app := webkit.New()
	req := httptest.NewRequest(http.MethodGet, "/bust", nil)
	rr := httptest.NewRecorder()

	store := cache.NewInMemory(time.Hour)
	store.Set(context.TODO(), "page:/bust", "test", cache.Options{
		Tags: []string{"payload"},
	})

	// Make sure the key is in the cache before attempting to bust it.
	var v string
	err := store.Get(context.TODO(), "page:/bust", &v)
	require.NoError(t, err)

	app.Get("/bust", cacheBust(store))

	app.ServeHTTP(rr, req)
	assert.Error(t, store.Get(context.TODO(), "page:/bust", &v))
}
