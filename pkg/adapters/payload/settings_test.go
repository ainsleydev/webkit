package payload

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/ainsleyclark/go-payloadcms"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

func TestSettingsMiddleware(t *testing.T) {
	t.Parallel()

	GlobalMiddlewareTestHelper(t, func(client *payloadcms.Client, store cache.Store) webkit.Plug {
		return SettingsMiddleware(client, store)
	})
}

func TestGetSettings(t *testing.T) {
	t.Parallel()

	s := &Settings{
		Id:       123,
		SiteName: ptr.StringPtr("Site Name"),
	}

	t.Run("OK", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), SettingsContextKey, s)
		got, err := GetSettings(ctx)
		assert.NoError(t, err)
		assert.Equal(t, s, got)
	})

	t.Run("Error", func(t *testing.T) {
		got, err := GetSettings(context.TODO())
		assert.Error(t, err, ErrSettingsNotFound)
		assert.Nil(t, got)
	})
}

func TestMustGetSettings(t *testing.T) {
	t.Parallel()

	s := &Settings{
		Id:       123,
		SiteName: ptr.StringPtr("Site Name"),
	}

	t.Run("OK", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), SettingsContextKey, s)
		got, err := GetSettings(ctx)
		assert.NoError(t, err)
		assert.Equal(t, s, got)
	})

	t.Run("Error", func(t *testing.T) {
		var buf bytes.Buffer
		slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))

		got := MustGetSettings(context.TODO())

		assert.Contains(t, buf.String(), ErrSettingsNotFound.Error())
		assert.Nil(t, got)
	})
}

func TestSettings_UnmarshalJSON(t *testing.T) {
	tt := map[string]struct {
		input   string
		want    Settings
		wantErr bool
	}{
		"OK": {
			input: `{
				"id": 1,
				"siteName": "Example Site",
				"tagLine": "An example tagline",
				"locale": "en_GB",
				"logo": 10,
				"extraField": "extraValue"
			}`,
			want: Settings{
				Id:       1,
				SiteName: ptr.StringPtr("Example Site"),
				TagLine:  ptr.StringPtr("An example tagline"),
				Locale:   "en_GB",
				Logo:     ptr.IntPtr(10),
				Extra: map[string]any{
					"extraField": "extraValue",
				},
			},
			wantErr: false,
		},
		"Invalid JSON": {
			input:   `{id: 1, siteName: "Example Site"}`,
			want:    Settings{},
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			var s Settings
			err := s.UnmarshalJSON([]byte(test.input))
			assert.Equal(t, test.wantErr, err != nil)
			assert.EqualValues(t, test.want, s)
		})
	}
}
