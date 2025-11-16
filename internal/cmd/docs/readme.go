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
	baseTemplate := templates.MustLoadTemplate("docs/README.md")

	logoURL := detectLogoURL(input.FS)
	domainBadges := collectDomainBadges(input.AppDef())
	domainLinks := formatDomainLinks(input.AppDef())
	appTypeBadges := collectAppTypeBadges(input.AppDef())
	providerGroups := groupByProvider(input.AppDef())
	appsWithPorts := collectAppsWithPorts(input.AppDef())
	primaryDomainURL := getPrimaryDomainURL(input.AppDef())
	enrichedResources := enrichResources(input.AppDef().Resources)

	data := map[string]any{
		"Definition":       input.AppDef(),
		"Content":          mustLoadCustomContent(input.FS, "README.md"),
		"LogoURL":          logoURL,
		"DomainBadges":     domainBadges,
		"DomainLinks":      domainLinks,
		"AppTypeBadges":    appTypeBadges,
		"ProviderGroups":   providerGroups,
		"AppsWithPorts":    appsWithPorts,
		"PrimaryDomainURL": primaryDomainURL,
		"CurrentYear":      time.Now().Year(),
		"Resources":        enrichedResources,
	}

	err := input.Generator().Template(
		"README.md",
		baseTemplate,
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

type domainBadge struct {
	Name string
}

// collectDomainBadges returns all primary domains for badge generation.
func collectDomainBadges(def *appdef.Definition) []domainBadge {
	var badges []domainBadge
	seen := make(map[string]bool)

	for _, app := range def.Apps {
		for _, domain := range app.Domains {
			if domain.Type == appdef.DomainTypePrimary && !seen[domain.Name] {
				badges = append(badges, domainBadge{Name: domain.Name})
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
		if primaryDomain := app.PrimaryDomain(); primaryDomain != "" {
			link := fmt.Sprintf(`<a href="https://%s"><strong>%s</strong></a>`, primaryDomain, app.Title)
			links = append(links, link)
		}
	}

	return strings.Join(links, " · ")
}

// collectAppTypeBadges returns all unique app types for badge generation.
func collectAppTypeBadges(def *appdef.Definition) []domainBadge {
	seen := make(map[appdef.AppType]bool)
	var badges []domainBadge

	for _, app := range def.Apps {
		if !seen[app.Type] {
			badges = append(badges, domainBadge{Name: string(app.Type)})
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

type appWithPort struct {
	Title string
	Port  int
}

// collectAppsWithPorts returns all apps that have a build port defined.
func collectAppsWithPorts(def *appdef.Definition) []appWithPort {
	var apps []appWithPort

	for _, app := range def.Apps {
		if app.Build.Port > 0 {
			apps = append(apps, appWithPort{
				Title: app.Title,
				Port:  app.Build.Port,
			})
		}
	}

	return apps
}

// getPrimaryDomainURL returns the first primary domain URL or empty string.
func getPrimaryDomainURL(def *appdef.Definition) string {
	for _, app := range def.Apps {
		if primaryDomain := app.PrimaryDomain(); primaryDomain != "" {
			return fmt.Sprintf("https://%s", primaryDomain)
		}
	}
	return ""
}

type (
	enrichedResource struct {
		appdef.Resource
		Description string
		Outputs     []resourceOutput
	}
	resourceOutput struct {
		Name        string
		Description string
	}
)

// enrichResources adds descriptions and outputs to resources for template use.
func enrichResources(resources []appdef.Resource) []enrichedResource {
	enriched := make([]enrichedResource, len(resources))

	for i, resource := range resources {
		enriched[i] = enrichedResource{
			Resource:    resource,
			Description: getResourceDescription(resource.Type),
			Outputs:     getResourceOutputs(resource.Type),
		}
	}

	return enriched
}

// getResourceDescription returns a human-readable description for a resource type.
func getResourceDescription(resourceType appdef.ResourceType) string {
	descriptions := map[appdef.ResourceType]string{
		appdef.ResourceTypePostgres: "PostgreSQL database for application data storage.",
		appdef.ResourceTypeS3:       "S3-compatible object storage for media and assets.",
		appdef.ResourceTypeSQLite:   "SQLite database with Turso for edge deployment.",
	}

	if desc, ok := descriptions[resourceType]; ok {
		return desc
	}

	return fmt.Sprintf("%s resource.", resourceType)
}

// getResourceOutputs returns the available outputs for a resource type with descriptions.
func getResourceOutputs(resourceType appdef.ResourceType) []resourceOutput {
	outputs := map[appdef.ResourceType][]resourceOutput{
		appdef.ResourceTypePostgres: {
			{Name: "connection_url", Description: "Full PostgreSQL connection string"},
			{Name: "host", Description: "Database host address"},
			{Name: "port", Description: "Database port number"},
			{Name: "database", Description: "Database name"},
			{Name: "user", Description: "Database username"},
			{Name: "password", Description: "Database password"},
		},
		appdef.ResourceTypeS3: {
			{Name: "bucket_name", Description: "S3 bucket identifier"},
			{Name: "bucket_url", Description: "Public bucket URL"},
			{Name: "region", Description: "Bucket region"},
		},
		appdef.ResourceTypeSQLite: {
			{Name: "connection_url", Description: "Full SQLite connection string"},
			{Name: "host", Description: "Turso database host"},
			{Name: "database", Description: "Database name"},
			{Name: "auth_token", Description: "Authentication token"},
		},
	}

	if out, ok := outputs[resourceType]; ok {
		return out
	}

	return []resourceOutput{}
}
