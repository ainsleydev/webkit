package middleware

import (
	"compress/gzip"
	"github.com/ainsleydev/webkit/pkg/webkit"
	"github.com/klauspost/compress/gzhttp"
	"net/http"
)

// Gzip returns a middleware that compresses HTTP responses using gzip compression.
// It wraps the provided handler, adding gzip compression to responses
// based on the specified configuration.
func Gzip(next webkit.Handler) webkit.Handler {
	//return webkit.WrapMiddelewareHandler(next, middleware.Compress(5))

	wrapper, _ := gzhttp.NewWrapper( //nolint Only returns on validation error.
		gzhttp.CompressionLevel(gzip.BestCompression),
		gzhttp.MinSize(0), // Compress responses larger than 512 bytes
		gzhttp.ContentTypes([]string{
			"text/html",
			"text/css",
			"text/plain",
			"text/javascript",
			"application/javascript",
			"application/x-javascript",
			"application/json",
			"application/atom+xml",
			"application/rss+xml",
			"image/svg+xml",
		}),
	)

	return func(c *webkit.Context) (err error) {
		wrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Response = w
			c.Request = r
			err = next(c)
		})).ServeHTTP(c.Response, c.Request)
		return
	}
}
