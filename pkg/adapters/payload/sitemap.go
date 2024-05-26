package payload

import "github.com/ainsleydev/webkit/pkg/webkit"

func (a Adapter) Sitemap() webkit.Handler {
	return func(c *webkit.Context) error {
		return nil
	}
}
