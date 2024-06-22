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
		return settingsMiddleware(client, store)
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
		withEnv bool
		want    Settings
		wantErr bool
	}{
		"OK": {
			input: `{
				"id": 1,
				"siteName": "Example Site",
				"tagLine": "An example tagline",
				"locale": "en_GB",
				"extraField": "extraValue"
			}`,
			want: Settings{
				Id:       1,
				SiteName: ptr.StringPtr("Example Site"),
				TagLine:  ptr.StringPtr("An example tagline"),
				Locale:   "en_GB",
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

func TestMaintenance_UnmarshalJSON(t *testing.T) {
	tt := map[string]struct {
		input   string
		want    Maintenance
		wantErr bool
	}{
		"Default": {
			input: "{}",
			want: Maintenance{
				Enabled: false,
				Title:   "",
				Content: "",
			},
			wantErr: false,
		},
		"OK": {
			input: `{
				"enabled": true,
				"title": "Maintenance Title",
				"content": "Maintenance Content"
			}`,
			want: Maintenance{
				Enabled: true,
				Title:   "Maintenance Title",
				Content: "Maintenance Content",
			},
			wantErr: false,
		},
		"Invalid JSON": {
			input:   `{enabled: wrong}`,
			want:    Maintenance{},
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			var m Maintenance
			err := m.UnmarshalJSON([]byte(test.input))
			assert.Equal(t, test.wantErr, err != nil)
			assert.EqualValues(t, test.want, m)
		})
	}
}

func TestFormat(t *testing.T) {
	tt := map[string]struct {
		input SettingsAddress
		want  string
	}{
		"All Fields Present": {
			input: SettingsAddress{
				Line1:    ptr.StringPtr("123 Main St"),
				Line2:    ptr.StringPtr("Suite 500"),
				City:     ptr.StringPtr("Metropolis"),
				County:   ptr.StringPtr("Gotham"),
				Postcode: ptr.StringPtr("12345"),
				Country:  ptr.StringPtr("UK"),
			},
			want: "123 Main St, Suite 500, Metropolis, Gotham, 12345, UK",
		},
		"Some Fields Nil": {
			input: SettingsAddress{
				Line1:    ptr.StringPtr("123 Main St"),
				City:     ptr.StringPtr("Metropolis"),
				Postcode: ptr.StringPtr("12345"),
			},
			want: "123 Main St, Metropolis, 12345",
		},
		"No Fields": {
			input: SettingsAddress{},
			want:  "",
		},
		"Only Line1 And Country": {
			input: SettingsAddress{
				Line1:   ptr.StringPtr("123 Main St"),
				Country: ptr.StringPtr("UK"),
			},
			want: "123 Main St, UK",
		},
		"Only Line1": {
			input: SettingsAddress{
				Line1: ptr.StringPtr("123 Main St"),
			},
			want: "123 Main St",
		},
		"Only Country": {
			input: SettingsAddress{
				Country: ptr.StringPtr("UK"),
			},
			want: "UK",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := test.input.Format()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestToStringArray(t *testing.T) {
	tt := map[string]struct {
		input SettingsSocial
		want  []string
	}{
		"Empty": {
			input: SettingsSocial{},
			want:  []string{},
		},
		"All Fields": {
			input: SettingsSocial{
				Facebook:  ptr.StringPtr("https://facebook.com/user"),
				Instagram: ptr.StringPtr("https://instagram.com/user"),
				LinkedIn:  ptr.StringPtr("https://linkedin.com/user"),
				Tiktok:    ptr.StringPtr("https://tiktok.com/@user"),
				X:         ptr.StringPtr("https://example.com/user"),
				Youtube:   ptr.StringPtr("https://youtube.com/user"),
			},
			want: []string{
				"https://facebook.com/user",
				"https://instagram.com/user",
				"https://linkedin.com/user",
				"https://tiktok.com/@user",
				"https://example.com/user",
				"https://youtube.com/user",
			},
		},
		"Some Fields Empty": {
			input: SettingsSocial{
				Facebook: ptr.StringPtr("https://facebook.com/user"),
				Tiktok:   ptr.StringPtr("https://tiktok.com/@user"),
			},
			want: []string{
				"https://facebook.com/user",
				"https://tiktok.com/@user",
			},
		},
	}
	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := test.input.ToStringArray()
			assert.ElementsMatch(t, test.want, got)
		})
	}
}
