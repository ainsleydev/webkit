package infra

import (
	"context"
	"fmt"
	"strings"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/platform/digitalocean"
	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/pkg/errors"
)

// FirewallUpdater handles automatic updates to PostgreSQL database firewalls
// by discovering app IPs and adding them to firewall rules.
type FirewallUpdater struct {
	apiToken string
}

// NewFirewallUpdater creates a new firewall updater.
func NewFirewallUpdater(apiToken string) *FirewallUpdater {
	return &FirewallUpdater{
		apiToken: apiToken,
	}
}

// PostgresAppMapping represents the mapping between a postgres resource
// and the apps that need access to it.
type PostgresAppMapping struct {
	Resource *appdef.Resource
	Apps     []appdef.App
}

// UpdateFirewalls updates PostgreSQL firewall rules by discovering app IPs
// and adding them to the allowed_ips_addr configuration.
//
// It returns a map of postgres resource names to the IPs that were discovered.
func (f *FirewallUpdater) UpdateFirewalls(ctx context.Context, appDef *appdef.Definition, tfOutput OutputResult, environment env.Environment) (map[string][]string, error) {
	// Find all postgres resources and the apps that reference them
	mappings := f.findPostgresAppMappings(appDef, environment)
	if len(mappings) == 0 {
		return nil, nil // No postgres resources or no apps referencing them
	}

	ipDiscovery := digitalocean.NewIPDiscovery(f.apiToken)
	result := make(map[string][]string)

	for _, mapping := range mappings {
		ips := []string{}

		// Discover IPs for each app
		for _, app := range mapping.Apps {
			appIPs, err := ipDiscovery.DiscoverAppIPs(ctx, app, tfOutput)
			if err != nil {
				// Log warning but continue with other apps
				// This allows partial success if some apps fail
				fmt.Printf("Warning: failed to discover IPs for app %q: %v\n", app.Name, err)
				continue
			}
			ips = append(ips, appIPs...)
		}

		if len(ips) > 0 {
			result[mapping.Resource.Name] = ips

			// Update the resource config with discovered IPs
			if err := f.updateResourceConfig(mapping.Resource, ips); err != nil {
				return result, errors.Wrapf(err, "updating config for postgres %q", mapping.Resource.Name)
			}
		}
	}

	return result, nil
}

// findPostgresAppMappings finds all postgres resources and the apps that
// reference them via environment variables.
func (f *FirewallUpdater) findPostgresAppMappings(appDef *appdef.Definition, environment env.Environment) []PostgresAppMapping {
	var mappings []PostgresAppMapping

	// Iterate through each postgres resource
	for i := range appDef.Resources {
		resource := &appDef.Resources[i]
		if resource.Type != appdef.ResourceTypePostgres {
			continue
		}

		var appsForResource []appdef.App

		// Find apps that reference this postgres resource
		for _, app := range appDef.Apps {
			if f.appReferencesResource(app, resource.Name, environment) {
				appsForResource = append(appsForResource, app)
			}
		}

		if len(appsForResource) > 0 {
			mappings = append(mappings, PostgresAppMapping{
				Resource: resource,
				Apps:     appsForResource,
			})
		}
	}

	return mappings
}

// appReferencesResource checks if an app references a specific resource
// by looking at its environment variables.
func (f *FirewallUpdater) appReferencesResource(app appdef.App, resourceName string, environment env.Environment) bool {
	// Merge with shared environment
	// Note: We'd need access to shared env here, but for now we'll just check app env
	found := false

	app.Env.Walk(func(entry appdef.EnvWalkEntry) {
		if entry.Environment != environment {
			return
		}

		if entry.Source != appdef.EnvSourceResource {
			return
		}

		// Parse the resource reference (e.g., "db.connection_url")
		refResource, _, ok := appdef.ParseResourceReference(entry.Value)
		if !ok {
			return
		}

		if refResource == resourceName {
			found = true
		}
	})

	return found
}

// updateResourceConfig updates the resource configuration with discovered IPs.
func (f *FirewallUpdater) updateResourceConfig(resource *appdef.Resource, ips []string) error {
	if resource.Config == nil {
		resource.Config = make(map[string]any)
	}

	// Get existing allowed IPs
	existingIPs := []string{}
	if existing, ok := resource.Config["allowed_ips_addr"]; ok {
		switch v := existing.(type) {
		case []string:
			existingIPs = v
		case []interface{}:
			for _, ip := range v {
				if ipStr, ok := ip.(string); ok {
					existingIPs = append(existingIPs, ipStr)
				}
			}
		}
	}

	// Merge with new IPs (deduplicate)
	ipSet := make(map[string]bool)
	for _, ip := range existingIPs {
		ipSet[ip] = true
	}
	for _, ip := range ips {
		ipSet[ip] = true
	}

	// Convert back to slice
	allIPs := make([]string, 0, len(ipSet))
	for ip := range ipSet {
		allIPs = append(allIPs, ip)
	}

	// Update the config
	resource.Config["allowed_ips_addr"] = allIPs

	return nil
}

// FormatUpdateSummary creates a user-friendly summary of the firewall updates.
func FormatUpdateSummary(updates map[string][]string) string {
	if len(updates) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\nℹ️  Auto-updating postgres firewalls...\n\n")

	for dbName, ips := range updates {
		sb.WriteString(fmt.Sprintf("  Database: %s\n", dbName))
		sb.WriteString(fmt.Sprintf("  ├─ Added %d IP(s): %s\n", len(ips), strings.Join(ips, ", ")))
		sb.WriteString("  └─ Firewall updated\n\n")
	}

	return sb.String()
}
