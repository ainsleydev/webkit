package payload

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ainsleyclark/go-payloadcms"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// redirect defines the structure of a redirect within the Payload UI.
type redirect struct {
	ID        float64      `json:"id"`
	From      string       `json:"from"`
	To        string       `json:"to"`
	Code      redirectCode `json:"code"`
	UpdatedAt time.Time    `json:"updated_at"`
	CreatedAt time.Time    `json:"created_at"`
}

// redirectCode defines the available redirect codes that are
// available within the Payload UI.
type redirectCode string

// The redirects available in the select dropdown within the Payload UI.
const (
	// redirectsCode301 - Moved Permanently
	redirectsCode301 redirectCode = "301"
	// redirectsCode302 - Found
	redirectsCode302 redirectCode = "302" //nolint:unused
	// redirectsCode307 - Temporary redirect
	redirectsCode307 redirectCode = "307" //nolint:unused
	// redirectsCode308 - Permanent redirect
	redirectsCode308 redirectCode = "308" //nolint:unused
	// redirectsCode410 - Content Gone (Deleted)
	redirectsCode410 redirectCode = "410" //nolint:unused
	// redirectsCode451 - Unavailable For Legal Reasons
	redirectsCode451 redirectCode = "451" //nolint:unused
)

const redirectCacheKey = "payload_redirects"

// redirectMiddleware is a middleware that checks for redirects within the Payload CMS.
// If a redirect is found, it will redirect the user to the new location according
// to the HTTP status specified. Cache is used to store the redirects to avoid
// making requests to the Payload API on every request.
func redirectMiddleware(client *payloadcms.Client, store cache.Store) webkit.Plug {
	return func(next webkit.Handler) webkit.Handler {
		return func(c *webkit.Context) error {
			// Skip caching for non-GET requests, file requests and ignored paths.
			if shouldSkipMiddleware(c) {
				return next(c)
			}

			var (
				ctx       = c.Request.Context()
				redirects []redirect
				path      = c.Request.URL.Path
			)

			err := store.Get(ctx, redirectCacheKey, &redirects)
			if err != nil {
				lr := payloadcms.ListResponse[redirect]{}
				_, err := client.Collections.List(ctx, CollectionRedirects, payloadcms.ListParams{
					Limit: payloadcms.AllItems,
				}, &lr)
				if err != nil {
					slog.Error("Fetching redirects from Payload with URL: " + path + ", Error: " + err.Error())
					return next(c)
				}
				redirects = lr.Docs

				store.Set(ctx, redirectCacheKey, lr.Docs, cache.Options{
					Expiration: time.Second * 30,
					Tags:       []string{"payload"},
				})
			}

			for _, r := range redirects {
				if r.From != path {
					continue
				}
				code, err := strconv.Atoi(string(r.Code))
				if err != nil {
					slog.Error("Converting redirect code to integer: " + err.Error())
					code = http.StatusMovedPermanently // Still continue despite parse error
				}
				slog.Debug("Redirecting URL",
					"from", r.From,
					"to", r.To,
					"code", code,
				)
				return c.Redirect(code, r.To)
			}

			return next(c)
		}
	}
}
