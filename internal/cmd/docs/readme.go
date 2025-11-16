package docs

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

const (
	webkitSymbolURL = "https://github.com/ainsleydev/webkit/blob/main/resources/symbol.png?raw=true"
	resourcesDir    = "resources"
)

// Readme creates the README.md file at the project root by combining
// the base template with project data from app.json.
func Readme(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	data := map[string]any{
		"Definition":     appDef,
		"Content":        mustLoadCustomContent(input.FS, "README.md"),
		"LogoURL":        detectLogoURL(input.FS),
		"DomainBadges":   collectDomainBadges(appDef),
		"DomainLinks":    formatDomainLinks(appDef),
		"AppTypeBadges":  collectAppTypeBadges(appDef),
		"ProviderGroups": groupByProvider(appDef),
		"CurrentYear":    time.Now().Year(),
	}

	err := input.Generator().Template(
		"README.md",
		templates.MustLoadTemplate("README.md"),
		data,
		scaffold.WithTracking(manifest.SourceProject()),
	)
	if err != nil {
		return errors.Wrap(err, "generating README.md")
	}

	return nil
}

// detectLogoURL checks for logo files in resources directory and returns
// the path or falls back to the WebKit symbol URL.
func detectLogoURL(fs afero.Fs) string {
	extensions := []string{"svg", "png", "jpg"}

	for _, ext := range extensions {
		logoPath := filepath.Join(resourcesDir, fmt.Sprintf("logo.%s", ext))
		exists, err := afero.Exists(fs, logoPath)
		if err == nil && exists {
			return fmt.Sprintf("./%s", logoPath)
		}
	}

	return webkitSymbolURL
}

// badge represents a shields.io badge with a name and color.
type badge struct {
	Name  string
	Color string
}

// collectDomainBadges returns all primary domains for badge generation.
func collectDomainBadges(def *appdef.Definition) []badge {
	var badges []badge
	seen := make(map[string]bool)

	for _, app := range def.Apps {
		for _, domain := range app.Domains {
			if domain.Type == appdef.DomainTypePrimary && !seen[domain.Name] {
				badges = append(badges, badge{
					Name:  domain.Name,
					Color: "", // Domain badges don't use color
				})
				seen[domain.Name] = true
			}
		}
	}

	return badges
}

// formatDomainLinks creates the HTML links for all primary domains.
func formatDomainLinks(def *appdef.Definition) string {
	var links []string

	for _, app := range def.Apps {
		if uri := app.PrimaryDomainURL(); uri != "" {
			link := fmt.Sprintf(`<a href="%s"><strong>%s</strong></a>`, uri, app.Title)
			links = append(links, link)
		}
	}

	return strings.Join(links, " · ")
}

// collectAppTypeBadges returns all unique app types for badge generation.
func collectAppTypeBadges(def *appdef.Definition) []badge {
	seen := make(map[appdef.AppType]bool)
	var badges []badge

	for _, app := range def.Apps {
		if !seen[app.Type] {
			badges = append(badges, badge{
				Name:  url.PathEscape(string(app.Type)),
				Color: "purple",
			})
			seen[app.Type] = true
		}
	}

	return badges
}

// groupByProvider groups apps and resources by their infrastructure provider.
func groupByProvider(def *appdef.Definition) map[string]string {
	groups := make(map[appdef.ResourceProvider][]string)

	for _, app := range def.Apps {
		if app.Infra.Provider != "" {
			groups[app.Infra.Provider] = append(
				groups[app.Infra.Provider],
				fmt.Sprintf("%s (App)", app.Title),
			)
		}
	}

	for _, resource := range def.Resources {
		groups[resource.Provider] = append(
			groups[resource.Provider],
			fmt.Sprintf("%s (%s)", resource.Name, resource.Type),
		)
	}

	result := make(map[string]string)
	for provider, items := range groups {
		result[string(provider)] = strings.Join(items, ", ")
	}

	return result
}
