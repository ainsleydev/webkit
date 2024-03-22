package middleware

import (
	"compress/gzip"

	"github.com/klauspost/compress/gzhttp"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

// Gzip returns a middleware that compresses HTTP responses using gzip compression.
// It wraps the provided handler, adding gzip compression to responses
// based on the specified configuration.
func Gzip(next webkit.Handler) webkit.Handler {
	wrapper, _ := gzhttp.NewWrapper( //nolint Only returns on validation error.
		// Use the best compression level so that the gzip header is always added.
		gzhttp.CompressionLevel(gzip.BestCompression),
		// Compress responses larger than 512 bytes
		gzhttp.MinSize(512),
		// TODO: Add more content types to compress such as CSS & JS.
		gzhttp.ContentTypes([]string{
			"text/plain",
			"text/html",
		}),
	)
	return webkit.WrapHandler(wrapper(next))
}
