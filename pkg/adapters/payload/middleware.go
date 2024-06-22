package payload

import (
	"strings"

	"github.com/ainsleydev/webkit/pkg/util/httputil"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

var skippable = []string{
	"/api",
	"robots.txt",
	"sitemap.xml",
}

// shouldSkipMiddleware determines if the request should skip some of
// the middleware.
func shouldSkipMiddleware(c *webkit.Context) bool {
	if httputil.IsFileRequest(c.Request) {
		return true
	}

	if c.Request.Method != "GET" {
		return true
	}

	path := c.Request.URL.Path
	for _, s := range skippable {
		if strings.Contains(path, s) {
			return true
		}
	}

	return false
}
