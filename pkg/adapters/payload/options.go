package payload

import (
	"github.com/ainsleyclark/go-payloadcms"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

//type Options struct {
//	Redirects bool
//	Settings  bool
//}

// Option is a functional option type that allows us to configure the Client.
type Option func(a *Adapter)

func WithCache(cache cache.Store) Option {
	return func(a *Adapter) {
		a.cache = cache
	}
}

func WithWebkit(kit *webkit.Kit) Option {
	return func(a *Adapter) {
		a.kit = kit
	}
}

// WithBaseURL is a functional option to set the base URL of the Payload API.
// Example: https://api.payloadcms.com
func WithBaseURL(url string) Option {
	return func(a *Adapter) {
		a.baseURL = url
	}
}

func WithEnvirons(environs map[string]string) Option {
	return func(a *Adapter) {
		a.env = environs
	}
}

func WithMaintenanceHandler(fn MaintenanceRendererFunc) Option {
	return func(a *Adapter) {
		a.maintenanceHandler = fn
	}
}

// WithAPIKey is a functional option to set the API key for the Payload API.
// To get an API key, visit: https://payloadcms.com/docs/rest-api/overview#authentication
//
// Usually, you can obtain one by enabling auth on the users type, and
// visiting the users collection in the Payload dashboard.
func WithAPIKey(apiKey string) Option {
	return func(a *Adapter) {
		a.apiKey = apiKey
	}
}

// WithGlobalMiddleware -- TODO
func WithGlobalMiddleware[T any](global string) Option {
	return func(a *Adapter) {
		a.globalMiddlewares = append(a.globalMiddlewares, func(client *payloadcms.Client, store cache.Store) webkit.Plug {
			return globalsMiddleware[T](client, store, global)
		})
	}
}
