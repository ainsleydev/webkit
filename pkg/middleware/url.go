package middleware

import (
	"fmt"
	"os"

	webkitctx "github.com/ainsleydev/webkit/pkg/context"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// URL is a middleware that sets the full URL of the request in the context.
// The URL can be accessed using the URLContextKey.
func URL(next webkit.Handler) webkit.Handler {
	return func(c *webkit.Context) error {
		url := os.Getenv("APP_URL")
		r := c.Request
		path := r.URL.Path
		if path == "/" {
			path = ""
		}
		if url != "" {
			url += path
		} else {
			host := r.Host
			if host == "" {
				host = "localhost"
			}
			url = fmt.Sprintf("%s://%s%s", c.Scheme(), host, path)
		}
		c.Request = c.Request.WithContext(webkitctx.WithURL(c.Context(), url))
		return next(c)
	}
}
