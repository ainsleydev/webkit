package payload

import (
	"context"
	"strings"
	"time"

	"github.com/ainsleydev/webkit/pkg/markup"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
	"github.com/ainsleydev/webkit/pkg/util/stringutil"
)

// TODO:
// - Merge page meta and settings meta, where page meta takes precedence.

func Head(ctx context.Context) markup.HeadProps {
	settings := MustGetSettings(ctx)

	var pageMeta Meta
	pm, ok := ctx.Value("meta").(Meta)
	if ok {
		pageMeta = pm
	}

	props := markup.HeadProps{
		Title:       ptr.String(settings.Meta.Title),
		Description: ptr.String(settings.Meta.Description),
		Locale:      settings.Locale,
		Hash:        time.Now().Unix(),
		Private:     ptr.Bool(settings.Meta.Private),
		Canonical:   ptr.String(settings.Meta.CanonicalURL),
		Org:         schemaOrganisation(settings, "TODO - Full URL"),
		Other:       "",
	}

	if settings.Meta.Image != nil {
		props.Image = settings.Meta.Image.URL
	}

	if settings.Meta.StructuredData != nil { // Type is a map[string]any
		props.Other += "<!-- Payload Structured Data -->\n" + markup.MarshalLDJSONScript(settings.Meta.StructuredData)
	}

	if settings.CodeInjection != nil && stringutil.IsNotEmpty(settings.CodeInjection.Head) {
		props.Other += "<!-- Payload Head Code Injection -->\n" + *settings.CodeInjection.Head
	}

	return props
}

func Foot(ctx context.Context) (string, error) {
	settings, err := GetSettings(ctx)
	if err != nil {
		return "", err
	}

	if settings.CodeInjection != nil && stringutil.IsNotEmpty(settings.CodeInjection.Footer) {
		return "<!-- Payload Foot Code Injection -->\n" + *settings.CodeInjection.Footer, nil
	}

	return "", nil
}

func schemaOrganisation(settings *Settings, url string) *markup.SchemaOrgOrganisation {
	org := markup.SchemaOrgOrganisation{
		Context: "https://schema.org",
		Type:    "Organization",
		ID:      url,
		URL:     url,
	}

	if stringutil.IsNotEmpty(settings.SiteName) {
		org.LegalName = *settings.SiteName
	}

	if stringutil.IsNotEmpty(settings.TagLine) {
		org.Description = strings.ReplaceAll(*settings.TagLine, "\n", " ")
	}

	if settings.Logo != nil {
		org.Logo = settings.Logo.URL
	}

	if settings.Social != nil {
		org.SameAs = settings.Social.ToStringArray()
	}

	if settings.Address != nil {
		org.Address = markup.SchemaOrgOrganisationAddress{
			Type:            "PostalAddress",
			StreetAddress:   settings.Address.Format(),
			AddressLocality: ptr.String(settings.Address.City),
			AddressRegion:   ptr.String(settings.Address.County),
			AddressCountry:  ptr.String(settings.Address.Country),
			PostalCode:      ptr.String(settings.Address.Postcode),
		}
	}

	return &org
}
