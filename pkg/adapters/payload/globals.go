package payload

import (
	"log/slog"
	"strings"

	"github.com/ainsleyclark/go-payloadcms"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/util/httputil"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// GlobalsMiddleware is a middleware that checks for globals within the Payload CMS.
// If a global is found, it will store it in the cache and in the context for
// easy access. Cache is used to store the globals to avoid making requests to
// the Payload API on every request.
func GlobalsMiddleware[T any](client *payloadcms.Client, store cache.Store, global string) webkit.Plug {
	return func(next webkit.Handler) webkit.Handler {
		return func(c *webkit.Context) error {
			if httputil.IsFileRequest(c.Request) {
				return next(c)
			}

			var (
				ctx       = c.Request.Context()
				t         = new(T)
				globalFmt = strings.ToTitle(global)
				cacheKey  = GlobalContextKey(global)
			)

			err := store.Get(ctx, cacheKey, t)
			if err == nil {
				c.Set(cacheKey, t)
				return next(c)
			}

			slog.Debug(globalFmt + " not found in cache, fetching from Payload")

			_, err = client.Globals.Get(ctx, payloadcms.Global(global), t)
			if err != nil {
				slog.Error("Fetching " + globalFmt + " global from Payload: " + err.Error())
				return next(c)
			}

			store.Set(ctx, cacheKey, t, cache.Options{
				Expiration: cache.Forever,
				Tags:       []string{"payload"},
			})

			c.Set(cacheKey, t)

			return next(c)
		}
	}
}

// GlobalContextKey returns the cache & context key for the global that
// resides in the context.
func GlobalContextKey(global string) string {
	return "payload_" + strings.ToLower(global)
}
