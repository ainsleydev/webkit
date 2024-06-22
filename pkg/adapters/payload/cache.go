package payload

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/util/httputil"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// cachePageExpiry is the time that a page will be cached for.
const cachePageExpiry = time.Hour * 24 * 7 * 4

// cacheMiddleware is a middleware increases performance of the application
// by caching full HTML pages instead of calling the Payload API on
// every request.
//
// If the method is not GET or the request is for a file, the request will be passed
// to the next http handler in the chain.
func cacheMiddleware(store cache.Store) webkit.Plug {
	return func(next webkit.Handler) webkit.Handler {
		return func(c *webkit.Context) error {
			ctx := c.Request.Context()

			// Skip caching for non-GET requests, file requests and ignored paths.
			if shouldSkipMiddleware(c) {
				return next(c)
			}

			cacheKey := fmt.Sprintf("page:%s", c.Request.URL.RequestURI())

			var page string
			if err := store.Get(ctx, cacheKey, &page); err == nil {
				// Cache hit, serve from cache
				c.Set("cache_hit", "HIT")
				c.Response.Header().Set("X-Cache", "HIT")
				c.Response.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(cachePageExpiry.Seconds())))
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

			if !httputil.Is2xx(rr.Status) {
				return nil
			}

			// Store the response in cache for future page requests.
			store.Set(ctx, cacheKey, rr.Body.String(), cache.Options{
				Expiration: cachePageExpiry,
				Tags:       []string{"payload"},
			})

			return nil
		}
	}
}

// CacheBust is a handler that can be used to clear the cache for a specific page.
func cacheBust(store cache.Store) webkit.Handler {
	type webhookRequest struct {
		Slug string `json:"slug,omitempty"`
	}

	return func(c *webkit.Context) error {
		ctx := c.Request.Context()

		//all, err := io.ReadAll(c.Request.Body)
		//if err != nil {
		//	return err
		//}
		//
		//fmt.Println(string(all))

		var w webhookRequest
		if err := c.BindJSON(&w); err != nil {
			return err
		}

		slog.Info("Payload cache busting for: " + w.Slug)

		// TODO:
		//
		// We need a way for the cache to be invalidated when a new page is published, edited or
		// deleted, we could do this by using Payload Hooks. At the moment we're just
		// flushing everything instead of invalidating a specific page.
		var err error
		switch w.Slug {
		case string(CollectionRedirects):
			err = store.Delete(ctx, redirectCacheKey)
		default:
			store.Invalidate(ctx, []string{"payload"})
		}

		if err != nil {
			slog.Error("Cache busting error: " + err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusAccepted)
	}
}
