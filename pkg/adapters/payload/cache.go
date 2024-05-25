package payload

import (
	"log/slog"
	"net/http"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

type CacheManager struct {
	Store cache.Store
}

func (c *CacheManager) Clear() {

}

// Need a way of obtaining the page contents so we can cache it after the
// request has been processed.
func (c *CacheManager) Middle(next webkit.Handler) webkit.Handler {
	return func(ctx *webkit.Context) error {
		if ctx.Request.Method != http.MethodGet {
			return next(ctx)
		}
		// If its HTML, not if it's css or anything
		path := ctx.Request.URL.Path

		var page string
		err := c.Store.Get(ctx.Context(), path, &page)
		if err != nil {
			slog.Debug("Cache miss: %v", err)
			return next(ctx)
		}

		ctx.Response.Header().Set("X-Cache", "HIT")

		_, err = ctx.Response.Write([]byte(page))
		if err != nil {
			return err
		}

		return nil
	}
}
