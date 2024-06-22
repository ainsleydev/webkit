package middleware

import (
	"fmt"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

// URLContextKey is the key used to retrieve the full URL in the context.
const URLContextKey = "url"

// URL is a middleware that sets the full URL of the request in the context.
// The URL can be accessed using the URLContextKey.
func URL(next webkit.Handler) webkit.Handler {
	return func(c *webkit.Context) error {
		u := c.Request.URL
		url := fmt.Sprintf("%s://%s%s", c.Scheme(), u.Host, u.Path)
		c.Set(URLContextKey, url)
		return next(c)
	}
}
