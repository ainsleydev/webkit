package payload

import (
	"net/http"

	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// Robots returns a handler function that generates the robots.txt content
// based on the default Payload settings.
func (a Adapter) Robots() webkit.Handler {
	return func(c *webkit.Context) error {
		settings := getSettings(c)

		if env.IsProduction() && settings.Robots == nil {
			return c.String(http.StatusOK, "User-agent: *\nDisallow:")
		} else if settings.Robots == nil {
			return c.String(http.StatusOK, "User-agent: *\nDisallow: /")
		}

		return c.String(http.StatusOK, *settings.Robots)
	}
}
