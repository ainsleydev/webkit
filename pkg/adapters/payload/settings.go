package payload

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/ainsleyclark/go-payloadcms"
	"github.com/goccy/go-json"
	"github.com/perimeterx/marshmallow"

	"github.com/ainsleydev/webkit/pkg/cache"
	"github.com/ainsleydev/webkit/pkg/markup"
	schemaorg "github.com/ainsleydev/webkit/pkg/markup/schema"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
	"github.com/ainsleydev/webkit/pkg/util/stringutil"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// SettingsContextKey defines the key for obtaining the settings
// from the context.
const SettingsContextKey = "payload_settings"

// ErrSettingsNotFound is returned when the settings are not found in the context.
var ErrSettingsNotFound = errors.New("settings not found in context")

// GetSettings is a helper function to get the settings from the context.
// If the settings are not found, it returns an error.
func GetSettings(ctx context.Context) (*Settings, error) {
	s := ctx.Value(SettingsContextKey)
	if s == nil {
		return nil, ErrSettingsNotFound
	}
	return s.(*Settings), nil
}

// WithSettings is a helper function to set the Payload Settings in the context.
func WithSettings(ctx context.Context, s *Settings) context.Context {
	return context.WithValue(ctx, SettingsContextKey, s) //nolint
}

// MustGetSettings is a helper function to get the settings from the context.
// If the settings are not found, it logs an error and returns nil.
func MustGetSettings(ctx context.Context) *Settings {
	s, err := GetSettings(ctx)
	if err != nil {
		slog.Error("Obtaining settings: " + err.Error())
	}
	// Guard Check for nil panics
	if s == nil {
		s = &Settings{}
	}
	return s
}

// settingsMiddleware defines the structure of the settings within the Payload UI.
func settingsMiddleware(client *payloadcms.Client, store cache.Store) webkit.Plug {
	return globalsMiddleware[Settings](client, store, "settings")
}

// Settings defines the common global collection type within Payload
// that allows users to change site settings.
type Settings struct {
	ID            float64                `json:"id"`
	SiteName      *string                `json:"siteName,omitempty"`
	TagLine       *string                `json:"tagLine,omitempty"`
	Locale        string                 `json:"locale,omitempty"` // In en_GB format
	Logo          *Media                 `json:"logo,omitempty"`
	Meta          SettingsMeta           `json:"meta"`
	Robots        *string                `json:"robots,omitempty"`
	CodeInjection *SettingsCodeInjection `json:"codeInjection,omitempty"`
	Maintenance   *Maintenance           `json:"maintenance,omitempty"`
	Contact       *SettingsContact       `json:"contact,omitempty"`
	Social        *SettingsSocial        `json:"social,omitempty"`
	Address       *SettingsAddress       `json:"address,omitempty"`
	UpdatedAt     *time.Time             `json:"updatedAt,omitempty"`
	CreatedAt     *time.Time             `json:"createdAt,omitempty"`
	Extra         map[string]any         `json:"-"` // Extra fields that are not defined in the struct.
}

// SettingsMeta defines the data generated by the Meta plugin from Payload
// along with additional fields such as Private & Canonical.
//
// The SEO plugin appears in the majority of collections and in both
// the Global Settings and Page level fields.
type SettingsMeta struct {
	Title          *string `json:"title,omitempty"`
	Description    *string `json:"description,omitempty"`
	Image          *Media  `json:"image,omitempty"`
	Private        *bool   `json:"private,omitempty"`
	CanonicalURL   *string `json:"canonicalURL,omitempty"`
	StructuredData any     `json:"structuredData,omitempty"`
}

// SettingsCodeInjection defines the fields for injecting code into the head
// or foot of the frontend.
type SettingsCodeInjection struct {
	Footer *string `json:"footer,omitempty"`
	Head   *string `json:"head,omitempty"`
}

// SettingsContact defines the fields for contact details for the company.
type SettingsContact struct {
	Email     *string `json:"email,omitempty"`
	Telephone *string `json:"telephone,omitempty"`
}

// SettingsAddress defines the fields for a company address.
type SettingsAddress struct {
	Line1    *string `json:"line1,omitempty"`
	Line2    *string `json:"line2,omitempty"`
	City     *string `json:"city,omitempty"`
	Country  *string `json:"country,omitempty"`
	County   *string `json:"county,omitempty"`
	Postcode *string `json:"postcode,omitempty"`
}

// SettingsSocial defines the fields for social media links.
type SettingsSocial struct {
	Facebook  *string `json:"facebook,omitempty"`
	Instagram *string `json:"instagram,omitempty"`
	LinkedIn  *string `json:"linkedIn,omitempty"`
	Tiktok    *string `json:"tiktok,omitempty"`
	X         *string `json:"x,omitempty"`
	Youtube   *string `json:"youtube,omitempty"`
}

// UnmarshalJSON unmarshalls the JSON data into the Settings type.
// This method is used to extract known fields and assign the remaining
// fields to the Extra map.
func (s *Settings) UnmarshalJSON(data []byte) error {
	var temp Settings
	result, err := marshmallow.Unmarshal(data,
		&temp,
		marshmallow.WithExcludeKnownFieldsFromMap(true),
	)
	if err != nil {
		return err
	}

	*s = temp
	s.Extra = result

	return nil
}

// UnmarshalJSON implements the custom unmarshalling logic for the Maintenance struct
// to make sure enabled is set to false if it's not present in the JSON data.
func (m *Maintenance) UnmarshalJSON(data []byte) error {
	m.Enabled = false
	m.Title = ""
	m.Content = ""

	// Define a struct to temporarily unmarshal into
	type alias Maintenance
	temp := (*alias)(m)

	return json.Unmarshal(data, &temp)
}

// Format returns the address as a comma-delimited string, excluding nil fields.
func (a SettingsAddress) Format() string {
	var parts []string
	addr := []*string{a.Line1, a.Line2, a.City, a.County, a.Country, a.Postcode}
	for _, field := range addr {
		if field != nil && *field != "" {
			parts = append(parts, *field)
		}
	}
	return strings.Join(parts, ", ")
}

// StringArray returns the social media links as an array of strings.
func (s SettingsSocial) StringArray() []string {
	var parts []string
	fields := []*string{s.Facebook, s.Instagram, s.LinkedIn, s.Tiktok, s.X, s.Youtube}
	for _, field := range fields {
		if field != nil && *field != "" {
			parts = append(parts, *field)
		}
	}
	return parts
}

// OpenGraph transforms the settings into an Open Graph object
// for use in the head of the frontend.
func (s *Settings) OpenGraph(url string) *markup.OpenGraph {
	m := &markup.OpenGraph{
		Type:        "website",
		SiteName:    ptr.String(s.SiteName),
		Title:       ptr.String(s.Meta.Title),
		Description: ptr.String(s.Meta.Description),
		URL:         url,
		Locale:      s.Locale,
	}
	if s.Meta.Image != nil {
		img := markup.OpengraphImage{
			URL:         s.Meta.Image.URL,
			ContentType: s.Meta.Image.MimeType,
			Width:       int(ptr.Float64(s.Meta.Image.Width)),
			Height:      int(ptr.Float64(s.Meta.Image.Height)),
		}
		if s.Meta.Image.Extra["alt"] != nil {
			img.Alt = s.Meta.Image.Extra["alt"].(string)
		}
		m.Image = append(m.Image, img)
	}
	return m
}

// TwitterCard transforms the settings into a Twitter Card
// for use in the head of the frontend.
func (s *Settings) TwitterCard() *markup.TwitterCard {
	card := &markup.TwitterCard{
		Title:       ptr.String(s.Meta.Title),
		Description: ptr.String(s.Meta.Description),
	}

	if s.Meta.Image != nil {
		card.Image = s.Meta.Image.URL
		if s.Meta.Image.Extra["alt"] != nil {
			card.ImageAlt = s.Meta.Image.Extra["alt"].(string)
		}
	}

	if s.Social == nil || stringutil.IsEmpty(s.Social.X) {
		return card
	}

	u, err := url.Parse(*s.Social.X)
	if err != nil {
		slog.Error("Parsing Twitter URL: " + err.Error())
		return card
	}

	if u.Path == "" {
		return card
	}

	// Assumes that the username is the last part of the path
	// For example: https://x.com/user (@user Tag)
	p := path.Base(strings.TrimSuffix(u.Path, "/"))
	p = strings.TrimPrefix(p, "@")
	card.Site = "@" + p
	card.Creator = "@" + p

	return card
}

// SchemaOrganisation transforms the settings into a Schema.org Organisation
// structure for use in the head of the frontend.
func (s *Settings) SchemaOrganisation(url string) *schemaorg.Organisation {
	org := schemaorg.Organisation{
		ID:  url,
		URL: url,
	}

	if stringutil.IsNotEmpty(s.SiteName) {
		org.LegalName = *s.SiteName
	}

	if stringutil.IsNotEmpty(s.TagLine) {
		org.Description = strings.ReplaceAll(*s.TagLine, "\n", " ")
	}

	if s.Logo != nil {
		org.Logo = s.Logo.URL
	}

	if s.Social != nil {
		org.SameAs = s.Social.StringArray()
	}

	if s.Address != nil {
		org.Address = schemaorg.Address{
			StreetAddress: s.Address.Format(),
			Locality:      ptr.String(s.Address.City),
			Region:        ptr.String(s.Address.County),
			Country:       ptr.String(s.Address.Country),
			PostalCode:    ptr.String(s.Address.Postcode),
		}
	}

	return &org
}
