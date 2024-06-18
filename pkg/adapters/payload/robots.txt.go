package payload

import (
	"net/http"

	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// Robots returns a handler function that generates the robots.txt content
// based on the default Payload settings.
func (a Adapter) Robots() webkit.Handler {
	defaultRobots := func(c *webkit.Context) error {
		// Always allow robots in production if it's not found via settings
		if env.IsProduction() {
			return c.String(http.StatusOK, "User-agent: *\nDisallow:")
		}

		// Disallow all robots in development or staging environments
		return c.String(http.StatusOK, "User-agent: *\nDisallow: /")
	}

	return func(c *webkit.Context) error {
		settings, err := GetSettings(c.Context())
		if err != nil {
			return defaultRobots(c)
		}

		robots := settings.Robots
		if robots != nil {
			return c.String(http.StatusOK, *robots)
		}

		return defaultRobots(c)
	}
}
