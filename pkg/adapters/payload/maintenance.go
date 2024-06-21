package payload

import (
	"strings"

	"github.com/ainsleydev/webkit/pkg/util/httputil"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// MaintenanceRendererFunc is the function that defines how the maintenance
// should be rendered on the site.
type MaintenanceRendererFunc func(c *webkit.Context, m Maintenance) error

var skippable = []string{
	"robots.txt",
	"sitemap.xml",
}

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

// MaintenanceMiddleware is a middleware that checks if the site is in maintenance mode.
// If it is, it will render the maintenance page as defined by the renderer function.
func MaintenanceMiddleware(fn MaintenanceRendererFunc) webkit.Plug {
	return func(next webkit.Handler) webkit.Handler {
		return func(c *webkit.Context) error {
			if shouldSkipMiddleware(c) {
				return next(c)
			}
			settings, err := GetSettings(c.Context())
			if err != nil {
				return next(c)
			}
			if settings.Maintenance == nil {
				return next(c)
			}
			if !settings.Maintenance.Enabled {
				return next(c)
			}
			return fn(c, *settings.Maintenance)
		}
	}
}
