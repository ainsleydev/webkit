package payload

import (
	"context"
	"log/slog"

	"github.com/ainsleyclark/go-payloadcms"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/util/httputil"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// SettingsContextKey defines the key for obtaining the settings
// from the context.
const SettingsContextKey = "payload_settings"

// settingsCacheKey is the key used to store the settings in the cache.
const settingsCacheKey = "payload_settings"

// SettingsMiddleware defines the structure of the settings within the Payload UI.
func SettingsMiddleware(client payloadcms.GlobalsService, store cache.Store) webkit.Plug {

	return func(next webkit.Handler) webkit.Handler {
		return func(c *webkit.Context) error {
			if httputil.IsFileRequest(c.Request) {
				return next(c)
			}

			var (
				ctx      = c.Request.Context()
				settings = Settings{}
			)

			err := store.Get(ctx, settingsCacheKey, &settings)
			if err == nil {
				c.Set(SettingsContextKey, settings)
				return next(c)
			}

			slog.Debug("Settings not found in cache, fetching from Payload")

			_, err = client.Get(ctx, GlobalSettings, &settings)
			if err != nil {
				slog.Error("Fetching redirects from Payload: " + err.Error())
				return next(c)
			}

			store.Set(ctx, settingsCacheKey, settings, cache.Options{
				Expiration: cache.Forever,
				Tags:       []string{"payload"},
			})

			c.Set(SettingsContextKey, settings)

			return next(c)
		}
	}
}

func getSettings(ctx context.Context) *Settings {
	s := ctx.Value(SettingsContextKey)
	if s == nil {
		return &Settings{}
	}
	return s.(*Settings)
}
