package payload

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ainsleyclark/go-payloadcms"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// Redirect defines the structure of a redirect within the Payload UI.
type Redirect struct {
	Id        float64      `json:"id"`
	From      string       `json:"from"`
	To        string       `json:"to"`
	Code      RedirectCode `json:"code"`
	UpdatedAt time.Time    `json:"updated_at"`
	CreatedAt time.Time    `json:"created_at"`
}

// RedirectCode defines the available redirect codes that are
// available within the Payload UI.
type RedirectCode string

// The redirects available in the select dropdown within the Payload UI.
const (
	// RedirectsCode301 - Moved Permanently
	RedirectsCode301 RedirectCode = "301"
	// RedirectsCode302 - Found
	RedirectsCode302 RedirectCode = "302"
	// RedirectsCode307 - Temporary Redirect
	RedirectsCode307 RedirectCode = "307"
	// RedirectsCode308 - Permanent Redirect
	RedirectsCode308 RedirectCode = "308"
	// RedirectsCode410 - Content Gone (Deleted)
	RedirectsCode410 RedirectCode = "410"
	// RedirectsCode451 - Unavailable For Legal Reasons
	RedirectsCode451 RedirectCode = "451"
)

const redirectCacheKey = "payload_redirects"

// RedirectMiddleware is a middleware that checks for redirects within the Payload CMS.
// If a redirect is found, it will redirect the user to the new location according
// to the HTTP status specified. Cache is used to store the redirects to avoid
// making requests to the Payload API on every request.
func RedirectMiddleware(client *payloadcms.Client, store cache.Store) webkit.Plug {
	return func(next webkit.Handler) webkit.Handler {
		return func(c *webkit.Context) error {
			var (
				ctx       = c.Request.Context()
				redirects []Redirect
				path      = c.Request.URL.Path
			)

			err := store.Get(ctx, redirectCacheKey, &redirects)
			if err != nil || env.IsDevelopment() {
				slog.Debug("Redirects not found in cache, fetching from Payload")

				lr := payloadcms.ListResponse[Redirect]{}
				_, err := client.Collections.List(ctx, CollectionRedirects, payloadcms.ListParams{
					Limit: payloadcms.AllItems,
				}, &lr)

				if err != nil {
					slog.Error("Fetching redirects from Payload: " + err.Error())
					return next(c)
				}
				redirects = lr.Docs

				store.Set(ctx, redirectCacheKey, lr.Docs, cache.Options{
					Expiration: time.Second * 30,
					Tags:       []string{"payload"},
				})
			}

			for _, r := range redirects {
				if r.From != strings.TrimSuffix(path, "/") {
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
