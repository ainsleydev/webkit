package payloadcms

import "github.com/ainsleydev/webkit/pkg/webkit"

func RedirectMiddleware(next webkit.Handler) webkit.Handler {
	return func(ctx *webkit.Context) error {
		return next(ctx)
	}
}
