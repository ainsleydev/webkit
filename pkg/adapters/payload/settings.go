package payload

import (
	"log/slog"

	"github.com/ainsleydev/webkit/pkg/apis/payloadcms"
	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/util/httputil"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

const settingsCacheKey = "payload_settings"

const SettingsContextKey = "payload_settings"

func SettingsMiddleware(client *payloadcms.Client, store cache.Store) webkit.Plug {
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
			if err != nil {
				slog.Debug("Settings not found in cache, fetching from Payload")

				_, err := client.Globals.Get(ctx, "settings", &settings)
				if err != nil {
					slog.Error("Fetching redirects from Payload: " + err.Error())
					return nil
				}

				err = store.Set(ctx, settingsCacheKey, settings, cache.Options{
					Expiration: cache.Forever,
					Tags:       []string{"payload"},
				})
				if err != nil {
					slog.Error("Setting settings in cache: " + err.Error())
				}
			}

			c.Set(SettingsContextKey, settings)

			return next(c)
		}
	}
}
