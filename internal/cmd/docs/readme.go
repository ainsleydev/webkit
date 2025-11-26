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

type (
	// ReadmeFrontMatter contains front matter metadata for README templates.
	ReadmeFrontMatter struct {
		Logo *LogoConfig `yaml:"logo,omitempty" json:"logo,omitempty"`
	}

	// LogoConfig contains logo display configuration.
	LogoConfig struct {
		Width  int `yaml:"width,omitempty" json:"width,omitempty"`
		Height int `yaml:"height,omitempty" json:"height,omitempty"`
	}

	// ReadmeContent contains parsed front matter and content.
	ReadmeContent struct {
		Meta    ReadmeFrontMatter
		Content string
	}

	// Logo contains the complete logo information including URL and dimensions.
	Logo struct {
		URL    string
		Width  int
		Height int
	}
)

// Readme creates the README.md file at the project root by combining
// the base template with project data from app.json.
func Readme(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	webkitOutputs := outputs.Load(input.FS)

	readmeContent, err := loadReadmeContent(input.FS)
	if err != nil {
		return errors.Wrap(err, "loading README content")
	}

	data := map[string]any{
		"Definition":     appDef,
		"Content":        readmeContent.Content,
		"Logo":           buildLogo(input.FS, readmeContent),
		"DomainLinks":    formatDomainLinks(appDef),
		"ProviderGroups": groupByProvider(appDef),
		"CurrentYear":    time.Now().Year(),
		"Outputs":        webkitOutputs,
		"StatusPageURL":  getStatusPageURL(appDef),
		"DashboardURL":   getDashboardURL(webkitOutputs),
		"MonitorBadges":  formatMonitorBadges(webkitOutputs),
	}

	err = input.Generator().Template(
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

// loadReadmeContent loads README content and parses front matter if present.
func loadReadmeContent(fs afero.Fs) (*ReadmeContent, error) {
	var meta ReadmeFrontMatter
	content, err := parseContentWithFrontMatter(
		fs,
		filepath.Join(customDocsDir, "README.md"),
		&meta,
	)
	if err != nil {
		return nil, err
	}

	return &ReadmeContent{
		Meta:    meta,
		Content: content,
	}, nil
}

// buildLogo constructs a Logo combining the detected URL and front matter dimensions.
func buildLogo(fs afero.Fs, content *ReadmeContent) Logo {
	logo := Logo{
		URL: detectLogoURL(fs),
	}

	if content.Meta.Logo != nil {
		logo.Width = content.Meta.Logo.Width
		logo.Height = content.Meta.Logo.Height
	}

	return logo
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
func getStatusPageURL(def *appdef.Definition) string {
	if def.Monitoring.StatusPage.Domain != "" {
		// Peekaping doesn't support https for CNAME's yet, might be good
		// to open an issue on their repo.
		return fmt.Sprintf("http://%s", def.Monitoring.StatusPage.Domain)
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

// getPeekapingEndpoint returns the Peekaping endpoint from outputs or the default.
func getPeekapingEndpoint(webkitOutputs *outputs.WebkitOutputs) string {
	if webkitOutputs != nil && webkitOutputs.Peekaping.Endpoint != "" {
		return webkitOutputs.Peekaping.Endpoint
	}
	return fmt.Sprintf("https://%s", defaultPeekapingDomain)
}

// getDashboardURL returns the dashboard URL filtered by project tag.
// Returns empty string if outputs is nil. The project tag is URL-encoded
// to handle special characters.
func getDashboardURL(webkitOutputs *outputs.WebkitOutputs) string {
	if webkitOutputs == nil {
		return ""
	}

	endpoint := getPeekapingEndpoint(webkitOutputs)
	projectTag := webkitOutputs.Peekaping.ProjectTag

	if projectTag == "" {
		return fmt.Sprintf("%s/monitors", endpoint)
	}

	return fmt.Sprintf("%s/monitors?tags=%s", endpoint, url.QueryEscape(projectTag))
}

// formatMonitorBadges creates badge data for all monitors from outputs.
func formatMonitorBadges(webkitOutputs *outputs.WebkitOutputs) []monitorBadge {
	if webkitOutputs == nil || len(webkitOutputs.Monitors) == 0 {
		return nil
	}

	endpoint := getPeekapingEndpoint(webkitOutputs)
	labels := "style=flat&upLabel=up&downLabel=down&pendingLabel=pending&maintenanceLabel=maintenance&pausedLabel=paused"

	badges := make([]monitorBadge, 0, len(webkitOutputs.Monitors))
	for _, m := range webkitOutputs.Monitors {
		badges = append(badges, monitorBadge{
			Name:     m.Name,
			BadgeURL: fmt.Sprintf("%s/api/v1/badge/%s/status?%s", endpoint, m.ID, labels),
			Type:     m.Type,
		})
	}

	return badges
}
