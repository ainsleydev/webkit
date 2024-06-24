package payload

import (
	"github.com/ainsleyclark/go-payloadcms"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// Option is a functional option type that allows us to configure the Client.
type Option func(a *Adapter)

// WithWebkit is a functional option to set the Webkit instance.
func WithWebkit(kit *webkit.Kit) Option {
	return func(a *Adapter) {
		a.kit = kit
	}
}

// WithCache is a functional option to set the cache store for the adapter.
func WithCache(cache cache.Store) Option {
	return func(a *Adapter) {
		a.cache = cache
	}
}

// WithBaseURL is a functional option to set the base URL of the Payload API.
// Example: https://api.payloadcms.com
func WithBaseURL(url string) Option {
	return func(a *Adapter) {
		a.baseURL = url
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

// WithMaintenanceHandler is a functional option to set the maintenance handler
// for the adapter. The maintenance handler is called when the site is in
// maintenance mode.
func WithMaintenanceHandler(fn MaintenanceRendererFunc) Option {
	return func(a *Adapter) {
		a.maintenanceHandler = fn
	}
}

// WithGlobalMiddleware is a functional option to set the global middleware for the adapter.
// The global middleware is applied to all requests and can be used to inject common data
// into the context.
//
// Global data can be accessed by using GlobalsContextKey(global)
//
// Example: payload.WithGlobalMiddleware[types.Navigation]("navigation")
func WithGlobalMiddleware[T any](global string) Option {
	return func(a *Adapter) {
		a.globalMiddlewares = append(a.globalMiddlewares, func(client *payloadcms.Client, store cache.Store) webkit.Plug {
			return globalsMiddleware[T](client, store, global)
		})
	}
}

// WithNavigation is a functional option to set the navigation global middleware for the adapter.
// The navigation middleware is used to inject the navigation data into the context.
func WithNavigation() Option {
	return func(a *Adapter) {
		a.globalMiddlewares = append(a.globalMiddlewares, func(client *payloadcms.Client, store cache.Store) webkit.Plug {
			return globalsMiddleware[Navigation](client, store, string(GlobalNavigation))
		})
	}
}
