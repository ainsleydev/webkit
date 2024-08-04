package payload

import (
	"context"
	"log/slog"

	"dario.cat/mergo"

	webkitctx "github.com/ainsleydev/webkit/pkg/context"
	"github.com/ainsleydev/webkit/pkg/markup"
	schemaorg "github.com/ainsleydev/webkit/pkg/markup/schema"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
	"github.com/ainsleydev/webkit/pkg/util/stringutil"
)

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

	url, ok := webkitctx.URL(ctx)
	if !ok {
		slog.Error("Error getting full URL from context")
		return markup.HeadProps{}
	}

	err = mergo.Merge(&settings.Meta, pageMeta, mergo.WithOverride, mergo.WithoutDereference)
	if err != nil {
		slog.Error("Merging page meta with settings meta: " + err.Error())
	}

	// ctx, ok := ctx.Value("navigation").(context.Context)

	props := markup.HeadProps{
		Title:        ptr.String(settings.Meta.Title),
		Description:  ptr.String(settings.Meta.Description),
		Locale:       settings.Locale,
		Private:      ptr.Bool(settings.Meta.Private),
		Canonical:    ptr.String(settings.Meta.CanonicalURL),
		OpenGraph:    settings.OpenGraph(url),
		Twitter:      settings.TwitterCard(),
		Organisation: settings.SchemaOrganisation(url),
		// Navigation:   ,
	}

	if settings.Meta.Image != nil {
		props.Image = settings.Meta.Image.URL
	}

	if settings.Meta.StructuredData != nil { // Type is a map[string]any
		props.Other += "<!-- Global Structured Data -->\n" + schemaorg.ToLDJSONScript(settings.Meta.StructuredData)
	}

	if pageMeta.StructuredData != nil { // Type is a map[string]any
		props.Other += "<!-- Page Structured Data -->\n" + schemaorg.ToLDJSONScript(pageMeta.StructuredData)
	}

	if settings.CodeInjection != nil && stringutil.IsNotEmpty(settings.CodeInjection.Head) {
		props.Other += "<!-- Global Head Code Injection -->\n" + *settings.CodeInjection.Head
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
