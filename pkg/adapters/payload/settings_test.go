package payload

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/ainsleyclark/go-payloadcms"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/markup"
	schemaorg "github.com/ainsleydev/webkit/pkg/markup/schema"
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
		ID:       123,
		SiteName: ptr.StringPtr("Site Name"),
	}

	t.Run("OK", func(t *testing.T) {
		t.Parallel()
		ctx := context.WithValue(context.Background(), SettingsContextKey, s)
		got, err := GetSettings(ctx)
		assert.NoError(t, err)
		assert.Equal(t, s, got)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		got, err := GetSettings(context.TODO())
		assert.Error(t, err, ErrSettingsNotFound)
		assert.Nil(t, got)
	})
}

func TestWithSettings(t *testing.T) {
	t.Parallel()

	s := &Settings{
		ID:       123,
		SiteName: ptr.StringPtr("Site Name"),
	}

	ctx := WithSettings(context.Background(), s)
	got, err := GetSettings(ctx)
	assert.NoError(t, err)
	assert.Equal(t, s, got)
}

func TestMustGetSettings(t *testing.T) {
	s := &Settings{
		ID:       123,
		SiteName: ptr.StringPtr("Site Name"),
	}

	t.Run("OK", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), SettingsContextKey, s)
		got := MustGetSettings(ctx)
		assert.Equal(t, s, got)
	})

	t.Run("Nil", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), SettingsContextKey, nil)
		got := MustGetSettings(ctx)
		assert.Equal(t, &Settings{}, got)
	})

	t.Run("Error", func(t *testing.T) {
		var buf bytes.Buffer
		slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))

		got := MustGetSettings(context.TODO())

		assert.Contains(t, buf.String(), ErrSettingsNotFound.Error())
		assert.Equal(t, &Settings{}, got)
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
				ID:       1,
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
		"All Tabs Present": {
			input: SettingsAddress{
				Line1:    ptr.StringPtr("123 Main St"),
				Line2:    ptr.StringPtr("Suite 500"),
				City:     ptr.StringPtr("Metropolis"),
				County:   ptr.StringPtr("Gotham"),
				Postcode: ptr.StringPtr("12345"),
				Country:  ptr.StringPtr("UK"),
			},
			want: "123 Main St, Suite 500, Metropolis, Gotham, UK, 12345",
		},
		"Some Tabs Nil": {
			input: SettingsAddress{
				Line1:    ptr.StringPtr("123 Main St"),
				City:     ptr.StringPtr("Metropolis"),
				Postcode: ptr.StringPtr("12345"),
			},
			want: "123 Main St, Metropolis, 12345",
		},
		"No Tabs": {
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
		"All Tabs": {
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
		"Some Tabs Empty": {
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
			got := test.input.StringArray()
			assert.ElementsMatch(t, test.want, got)
		})
	}
}

func TestSettings_MarkupOpenGraph(t *testing.T) {
	t.Parallel()

	url := "https://example.com"

	tt := map[string]struct {
		input Settings
		want  markup.OpenGraph
	}{
		"Default": {
			input: Settings{},
			want: markup.OpenGraph{
				Type: "website",
				URL:  url,
			},
		},
		"Full": {
			input: Settings{
				SiteName: ptr.StringPtr("Example Site"),
				Locale:   "en_GB",
				Meta: SettingsMeta{
					Title:       ptr.StringPtr("Title"),
					Description: ptr.StringPtr("Description"),
					Image: &Media{
						URL:      "https://example.com/image.jpg",
						MimeType: "image/jpeg",
						Width:    ptr.Float64Ptr(1200),
						Height:   ptr.Float64Ptr(630),
						Extra: map[string]interface{}{
							"alt": "Alt",
						},
					},
				},
			},
			want: markup.OpenGraph{
				Type:        "website",
				SiteName:    "Example Site",
				Title:       "Title",
				Description: "Description",
				URL:         url,
				Locale:      "en_GB",
				Image: []markup.OpengraphImage{
					{
						URL:         "https://example.com/image.jpg",
						ContentType: "image/jpeg",
						Width:       1200,
						Height:      630,
						Alt:         "Alt",
					},
				},
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.OpenGraph(url)
			assert.Equal(t, &test.want, got)
		})
	}
}

func TestSettings_MarkupTwitterCard(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Settings
		want  markup.TwitterCard
	}{
		"Default": {
			input: Settings{},
			want:  markup.TwitterCard{},
		},
		"Simple": {
			input: Settings{
				Meta: SettingsMeta{
					Title:       ptr.StringPtr("Title"),
					Description: ptr.StringPtr("Description"),
				},
			},
			want: markup.TwitterCard{
				Title:       "Title",
				Description: "Description",
			},
		},
		"With Image": {
			input: Settings{
				Meta: SettingsMeta{
					Title:       ptr.StringPtr("Title"),
					Description: ptr.StringPtr("Description"),
					Image: &Media{
						URL: "https://example.com/image.jpg",
						Extra: map[string]interface{}{
							"alt": "Alt",
						},
					},
				},
			},
			want: markup.TwitterCard{
				Title:       "Title",
				Description: "Description",
				Image:       "https://example.com/image.jpg",
				ImageAlt:    "Alt",
			},
		},
		"With Creator": {
			input: Settings{
				Social: &SettingsSocial{
					X: ptr.StringPtr("https://example.com/user"),
				},
			},
			want: markup.TwitterCard{
				Site:    "@user",
				Creator: "@user",
			},
		},
		"Parse URL Error": {
			input: Settings{
				Social: &SettingsSocial{
					X: ptr.StringPtr(":invalidScheme://example.com"),
				},
			},
			want: markup.TwitterCard{},
		},
		"No Path": {
			input: Settings{
				Social: &SettingsSocial{
					X: ptr.StringPtr("https://x.com"),
				},
			},
			want: markup.TwitterCard{},
		},
		"Full": {
			input: Settings{
				SiteName: ptr.StringPtr("Example Site"),
				Locale:   "en_GB",
				Social: &SettingsSocial{
					X: ptr.StringPtr("https://example.com/user"),
				},
				Meta: SettingsMeta{
					Title:       ptr.StringPtr("Title"),
					Description: ptr.StringPtr("Description"),
					Image: &Media{
						URL:      "https://example.com/image.jpg",
						MimeType: "image/jpeg",
						Width:    ptr.Float64Ptr(1200),
						Height:   ptr.Float64Ptr(630),
						Extra: map[string]interface{}{
							"alt": "Alt",
						},
					},
				},
			},
			want: markup.TwitterCard{
				Site:        "@user",
				Creator:     "@user",
				Title:       "Title",
				Description: "Description",
				Image:       "https://example.com/image.jpg",
				ImageAlt:    "Alt",
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.TwitterCard()
			assert.EqualValues(t, &test.want, got)
		})
	}
}

func TestSettings_MarkupSchemaOrganisation(t *testing.T) {
	t.Parallel()

	url := "https://example.com"

	tt := map[string]struct {
		input Settings
		want  schemaorg.Organisation
	}{
		"Default": {
			input: Settings{},
			want: schemaorg.Organisation{
				ID:  url,
				URL: url,
			},
		},
		"Full": {
			input: Settings{
				SiteName: ptr.StringPtr("Site"),
				TagLine:  ptr.StringPtr("TagLine\nNew Line"),
				Logo:     &Media{URL: "https://example.com/logo.png"},
				Social: &SettingsSocial{
					Facebook:  ptr.StringPtr("https://facebook.com/example"),
					Instagram: ptr.StringPtr("https://instagram.com/example"),
					LinkedIn:  ptr.StringPtr("https://linkedin.com/example"),
					Tiktok:    ptr.StringPtr("https://tiktok.com/@example"),
					X:         ptr.StringPtr("https://x.com/example"),
					Youtube:   ptr.StringPtr("https://youtube.com/example"),
				},
				Address: &SettingsAddress{
					Line1:    ptr.StringPtr("Line 1"),
					Line2:    ptr.StringPtr("Line 2"),
					City:     ptr.StringPtr("City"),
					County:   ptr.StringPtr("County"),
					Postcode: ptr.StringPtr("12345"),
					Country:  ptr.StringPtr("Country"),
				},
			},
			want: schemaorg.Organisation{
				ID:          url,
				URL:         url,
				LegalName:   "Site",
				Description: "TagLine New Line",
				Logo:        "https://example.com/logo.png",
				SameAs: []string{
					"https://facebook.com/example",
					"https://instagram.com/example",
					"https://linkedin.com/example",
					"https://tiktok.com/@example",
					"https://x.com/example",
					"https://youtube.com/example",
				},
				Address: schemaorg.Address{
					StreetAddress: "Line 1, Line 2, City, County, Country, 12345",
					Locality:      "City",
					Region:        "County",
					Country:       "Country",
					PostalCode:    "12345",
				},
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.SchemaOrganisation(url)
			assert.Equal(t, &test.want, got)
		})
	}
}
