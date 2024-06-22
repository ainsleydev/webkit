package payload

import (
	"context"
	"log/slog"

	"dario.cat/mergo"

	"github.com/ainsleydev/webkit/pkg/markup"
	"github.com/ainsleydev/webkit/pkg/middleware"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
	"github.com/ainsleydev/webkit/pkg/util/stringutil"
)

// TODO:
// - Merge page meta and settings meta, where page meta takes precedence.

type Navigation struct {
	Items []NavigationItem `json:"items"`
}

type NavigationItem struct {
	Label string `json:"label"`
	Link  string `json:"link"`
}

const ContextKeyPageMeta = "payload_page_meta"

func Head(ctx context.Context) markup.HeadProps {
	settings, err := GetSettings(ctx)
	if err != nil {
		return markup.HeadProps{}
	}

	pageMeta := &SettingsMeta{}
	pm, ok := ctx.Value(ContextKeyPageMeta).(*SettingsMeta)
	if ok {
		pageMeta = pm
	}

	url, ok := ctx.Value(middleware.URLContextKey).(string)
	if !ok {
		slog.Error("Error getting full URL from context under key: " + middleware.URLContextKey)
		return markup.HeadProps{}
	}

	err = mergo.Merge(&settings.Meta, pageMeta, mergo.WithOverride, mergo.WithoutDereference)
	if err != nil {
		slog.Error("Merging page meta with settings meta: " + err.Error())
	}

	//ctx, ok := ctx.Value("navigation").(context.Context)

	props := markup.HeadProps{
		Title:        ptr.String(settings.Meta.Title),
		Description:  ptr.String(settings.Meta.Description),
		Locale:       settings.Locale,
		Private:      ptr.Bool(settings.Meta.Private),
		Canonical:    ptr.String(settings.Meta.CanonicalURL),
		OpenGraph:    settings.MarkupOpenGraph(url),
		Twitter:      settings.MarkupTwitterCard(),
		Organisation: settings.MarkupSchemaOrganisation(url),
		//Navigation:   ,
	}

	if settings.Meta.Image != nil {
		props.Image = settings.Meta.Image.URL
	}

	if settings.Meta.StructuredData != nil { // Type is a map[string]any
		props.Other += "<!-- Global Structured Data -->\n" + markup.MarshalLDJSONScript(settings.Meta.StructuredData)
	}

	if pageMeta.StructuredData != nil { // Type is a map[string]any
		props.Other += "<!-- Page Structured Data -->\n" + markup.MarshalLDJSONScript(pageMeta.StructuredData)
	}

	if settings.CodeInjection != nil && stringutil.IsNotEmpty(settings.CodeInjection.Head) {
		props.Other += "<!-- Global Head Code Injection -->\n" + *settings.CodeInjection.Head
	}

	// TODO: We need to add code injection for page meta here.

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
