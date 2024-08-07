package payload

import (
	"net/http"
	"strings"

	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// robots returns a handler function that generates the robots.txt content
// based on the default Payload settings.
func robots(appEnv string) webkit.Handler {
	defaultRobots := func(c *webkit.Context) error {
		// Always allow robots in production if it's not found via settings
		if appEnv == env.Production {
			return c.String(http.StatusOK, "User-agent: *\nDisallow:")
		}

		// Disallow all robots in development or staging environments
		return c.String(http.StatusOK, "User-agent: *\nDisallow: /")
	}

	return func(c *webkit.Context) error {
		// Don't allow search engines to crawl if it's a Digital
		// Ocean URL: ondigitalocean.app/admin
		if strings.Contains(c.Request.URL.String(), "digitalocean") {
			return c.String(http.StatusOK, "User-agent: *\nDisallow: /")
		}

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
