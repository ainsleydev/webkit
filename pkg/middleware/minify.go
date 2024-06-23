package middleware

import (
	"log/slog"

	"github.com/tdewolff/minify/v2"

	"github.com/tdewolff/minify/v2/html"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

// Minify is a middleware that minifies the response body of the request.
// Only text/html content type is minified. Assets such as CSS & JS
// should be minified when built (obviously).
//
// IMPORTANT NOTE: This middleware should be used after any caching
// middleware to ensure that the response is minified before being cached.
// In the event that this is called before the full page cache, it will
// add unnecessary overhead to the response times.
//
//	app := webkit.New()
//	app.Plug(middleware.Cache)        // <--<< Cache should come before Minifier
//	app.Plug(middleware.Minify)
//	r.Get("/", handler)
func Minify(next webkit.Handler) webkit.Handler {
	m := minify.New()
	m.Add("text/html", &html.Minifier{
		KeepComments:        true,
		KeepSpecialComments: true,
		KeepDocumentTags:    true,
		KeepEndTags:         true,
	})

	return func(c *webkit.Context) error {
		// ResponseWriter minifies any writes to the http.ResponseWriter.
		mw := m.ResponseWriter(c.Response, c.Request)
		defer func() {
			if err := mw.Close(); err != nil {
				slog.Error("Failed to close minified response writer: " + err.Error())
			}
		}()

		// Replace the original response writer with the minified one
		c.Response = mw

		return next(c)
	}
}
