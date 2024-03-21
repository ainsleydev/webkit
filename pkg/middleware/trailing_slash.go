package middleware

import (
	"fmt"

	webkit "github.com/ainsleydev/webkit/pkg/webkit"
)

// RedirectSlashes is a middleware that will match request paths with a trailing
// slash and redirect to the same path, less the trailing slash.
//
// NOTE: RedirectSlashes middleware is *incompatible* with http.FileServer.
func RedirectSlashes(next webkit.Handler) webkit.Handler {
	return func(ctx *webkit.Context) error {
		r := ctx.Request
		path := r.URL.Path

		if len(path) > 1 && path[len(path)-1] != '/' {
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s/?%s", path, r.URL.RawQuery)
			} else {
				path = path + "/"
			}
			redirectURL := fmt.Sprintf("//%s%s", r.Host, path)
			return ctx.Redirect(301, redirectURL)
		}
		return next(ctx)
	}
}
