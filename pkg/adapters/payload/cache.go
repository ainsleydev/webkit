package payload

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/util/httputil"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// TODO: We need a way for the cache to be invalidated when a new page is published, edited or
// deleted, we could do this by using Payload Hooks.

// cachePageExpiry is the time that a page will be cached for.
const cachePageExpiry = time.Hour * 24 * 7 * 4

// CacheBust is a handler that can be used to clear the cache for a specific page.
func CacheBust(store cache.Store) webkit.Handler {
	return func(c *webkit.Context) error {
		_ = c.Request.Context()
		return nil
	}
}

// https://github.com/ainsleydev/audits.com/blob/691badc3cc142f13122a3ed6e86b4a0046824916/backend/config/plugins.ts#L69

// CacheMiddleware is a middleware increases performance of the application
// by caching full HTML pages instead of calling the Payload API on
// every request.
// If the method is not GET or the request is for a file, the request will be passed
// to the next http handler in the chain.
func CacheMiddleware(store cache.Store, ignorePaths []string) webkit.Plug {
	return func(next webkit.Handler) webkit.Handler {
		return func(c *webkit.Context) error {
			ctx := c.Request.Context()

			if c.Request.Method != http.MethodGet {
				return next(c)
			}
			if httputil.IsFileRequest(c.Request) {
				return next(c)
			}

			cacheKey := fmt.Sprintf("page:%s", c.Request.URL.RequestURI())

			var page string
			if err := store.Get(ctx, cacheKey, &page); err == nil {
				// Cache hit, serve from cache
				c.Set("cache_hit", "HIT")
				c.Response.Header().Set("X-Cache", "HIT")
				c.Response.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%f", cachePageExpiry.Seconds()))
				return c.HTML(http.StatusOK, page)
			}

			rr := httputil.NewResponseRecorder(c.Response)
			c.Set("cache_hit", "MISS")
			c.Response.Header().Set("X-Cache", "MISS")
			c.Response = rr

			// Process next request in chain.
			if err := next(c); err != nil {
				return err
			}

			if rr.Status != http.StatusOK {
				return nil
			}

			// Store the response in cache for future page requests.
			return store.Set(ctx, cacheKey, rr.Body.String(), cache.Options{
				Expiration: cachePageExpiry,
				Tags:       []string{"payload"},
			})
		}
	}
}
