package payload

import (
	"github.com/ainsleyclark/go-payloadcms"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

type Adapter struct {
	kit    *webkit.Kit
	opts   Options
	client *payloadcms.Client
	cache  cache.Store
}

type Options struct {
	Redirects bool
	Settings  bool
}

func NewAdapter(
	opts Options,
	env string,
	kit *webkit.Kit,
	client *payloadcms.Client,
	cache cache.Store,
) *Adapter {
	return &Adapter{
		opts:   opts,
		kit:    kit,
		client: client,
		cache:  cache,
	}
}

func (a Adapter) initMiddleware() {
	a.kit.Plug(SettingsMiddleware(a.client, a.cache))
	a.kit.Plug(RedirectMiddleware(a.client, a.cache))
}

const (
	// CollectionMedia defines the Payload media collection slug.
	CollectionMedia payloadcms.Collection = "media"
	// CollectionUsers defines the Payload users collection slug.
	CollectionUsers payloadcms.Collection = "users"
	// CollectionRedirects defines the Payload redirects collection slug.
	CollectionRedirects payloadcms.Collection = "redirects"
)

const (
	// GlobalSettings defines the Payload settings global settings slug.
	GlobalSettings payloadcms.Global = "settings"
)
