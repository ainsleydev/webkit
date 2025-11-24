package docs

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/state/manifest"
	"github.com/ainsleydev/webkit/internal/state/outputs"
	"github.com/ainsleydev/webkit/internal/templates"
)

const (
	webkitSymbolURL        = "https://github.com/ainsleydev/webkit/blob/main/resources/symbol.png?raw=true"
	resourcesDir           = "resources"
	defaultPeekapingDomain = "uptime.ainsley.dev"
)

// Readme creates the README.md file at the project root by combining
// the base template with project data from app.json.
func Readme(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	webkitOutputs := outputs.Load(input.FS)

	data := map[string]any{
		"Definition":     appDef,
		"Content":        mustLoadCustomContent(input.FS, "README.md"),
		"LogoURL":        detectLogoURL(input.FS),
		"DomainLinks":    formatDomainLinks(appDef),
		"ProviderGroups": groupByProvider(appDef),
		"CurrentYear":    time.Now().Year(),
		"Outputs":        webkitOutputs,
		"StatusPageURL":  getStatusPageURL(appDef, webkitOutputs),
		"MonitorBadges":  formatMonitorBadges(appDef, webkitOutputs),
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
			fmt.Sprintf("%s (%s)", resource.Title, resource.Type),
		)
	}

	result := make(map[string]string)
	for provider, items := range groups {
		result[string(provider)] = strings.Join(items, ", ")
	}

	return result
}

// getStatusPageURL returns the status page URL based on appdef config.
// Priority: StatusPage.Domain > StatusPage.Slug > default peekaping domain.
func getStatusPageURL(def *appdef.Definition, webkitOutputs *outputs.WebkitOutputs) string {
	if def.Monitoring.StatusPage.Domain != "" {
		return fmt.Sprintf("https://%s", def.Monitoring.StatusPage.Domain)
	}

	if def.Monitoring.StatusPage.Slug != "" {
		return fmt.Sprintf("https://%s/status/%s", defaultPeekapingDomain, def.Monitoring.StatusPage.Slug)
	}

	return fmt.Sprintf("https://%s", defaultPeekapingDomain)
}

// monitorBadge represents a single monitor's badge data for the README.
type monitorBadge struct {
	Name     string
	BadgeURL string
	Type     string
}

// formatMonitorBadges creates badge data for all monitors from outputs.
func formatMonitorBadges(def *appdef.Definition, webkitOutputs *outputs.WebkitOutputs) []monitorBadge {
	if webkitOutputs == nil || len(webkitOutputs.Monitors) == 0 {
		return nil
	}

	endpoint := webkitOutputs.PeekapingEndpoint
	if endpoint == "" {
		endpoint = fmt.Sprintf("https://%s", defaultPeekapingDomain)
	}

	badges := make([]monitorBadge, 0, len(webkitOutputs.Monitors))
	for _, m := range webkitOutputs.Monitors {
		badges = append(badges, monitorBadge{
			Name:     m.Name,
			BadgeURL: fmt.Sprintf("%s/api/v1/badge/%s/status", endpoint, m.ID),
			Type:     m.Type,
		})
	}

	return badges
}
