package payload

import (
	"errors"
	"os"
	"time"

	"github.com/ainsleyclark/go-payloadcms"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

type Adapter struct {
	*payloadcms.Client
	kit                *webkit.Kit
	cache              cache.Store
	baseURL            string
	apiKey             string
	env                map[string]string
	maintenanceHandler MaintenanceRendererFunc
	globalMiddlewares  []func(*payloadcms.Client, cache.Store) webkit.Plug
}

func NewAdapter(options ...Option) (*Adapter, error) {
	a := &Adapter{
		cache:              cache.NewInMemory(time.Hour * 24), // Default cache store.
		maintenanceHandler: defaultMaintenanceRenderer,
	}

	// Apply all the functional options to configure the client.
	for _, opt := range options {
		opt(a)
	}

	// Ensure all required fields are set.
	if err := a.validate(); err != nil {
		return nil, err
	}

	// Instantiate the Payload HTTP Client
	client, err := payloadcms.New(
		payloadcms.WithBaseURL(a.baseURL),
		payloadcms.WithAPIKey(a.apiKey),
	)
	if err != nil {
		return nil, err
	}
	a.Client = client

	// Set the Payload URL in the environment just in case it's not defined
	// in the env file. Used for media URLs and other utilities.
	if err = os.Setenv(EnvPayloadURL, a.baseURL); err != nil {
		return nil, err
	}

	a.attachHandlers()

	return a, nil
}

func (a Adapter) validate() error {
	if a.kit == nil {
		return errors.New("kit is required")
	}
	return nil
}

func (a Adapter) attachHandlers() {
	a.kit.Get("/robots.txt", robots(a.env[env.AppEnvironmentKey]))
	a.kit.Get("/sitemap.xml", sitemap())
	a.kit.Plug(redirectMiddleware(a.Client, a.cache))
	a.kit.Plug(settingsMiddleware(a.Client, a.cache))
	a.kit.Plug(maintenanceMiddleware(a.maintenanceHandler))
	for _, m := range a.globalMiddlewares {
		a.kit.Plug(m(a.Client, a.cache))
	}
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

const (
	EnvPayloadURL = "PAYLOAD_URL"
)
