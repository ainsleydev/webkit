package payload

import (
	"net/http"

	"github.com/ainsleydev/webkit/pkg/webkit"
)

// MaintenanceRendererFunc is the function that defines how the maintenance
// should be rendered on the site.
type MaintenanceRendererFunc func(c *webkit.Context, m Maintenance) error

// Maintenance defines the fields for displaying an offline page to
// the front-end when it's been enabled within PayloadCMS.
type Maintenance struct {
	Enabled bool   `json:"enabled,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

var defaultMaintenanceRenderer = func(c *webkit.Context, _ Maintenance) error {
	return c.String(http.StatusServiceUnavailable,
		"Site is under maintenance. Please check back soon",
	)
}

// maintenanceMiddleware is a middleware that checks if the site is in maintenance mode.
// If it is, it will render the maintenance page as defined by the renderer function.
func maintenanceMiddleware(fn MaintenanceRendererFunc) webkit.Plug {
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
