package infra

import (
	"fmt"
	"strings"

	"github.com/ainsleydev/webkit/internal/appdef"
)

type (
	// importAddress represents a single Terraform resource to import.
	importAddress struct {
		// Address is the Terraform resource address
		// (e.g., "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_cluster.this").
		Address string
		// ID is the provider-specific resource identifier.
		ID string
	}
)

// buildImportAddresses constructs the list of Terraform import addresses
// for a given resource based on its type and provider.
// The projectName is used to build the full resource name as Terraform modules do.
func buildImportAddresses(projectName string, resource *appdef.Resource, baseID string) ([]importAddress, error) {
	switch resource.Provider {
	case appdef.ResourceProviderDigitalOcean:
		return buildDigitalOceanImports(projectName, resource, baseID)
	default:
		return nil, fmt.Errorf("import not supported for provider %q", resource.Provider)
	}
}

// buildAppImportAddresses constructs the list of Terraform import addresses
// for a given app based on its platform type and provider.
// The projectName is used to build the full app name as Terraform modules do.
func buildAppImportAddresses(projectName string, app *appdef.App, appID string) ([]importAddress, error) {
	switch app.Infra.Provider {
	case appdef.ResourceProviderDigitalOcean:
		return buildDigitalOceanAppImports(projectName, app, appID)
	default:
		return nil, fmt.Errorf("import not supported for provider %q", app.Infra.Provider)
	}
}

// buildDigitalOceanImports creates import addresses for DigitalOcean resources.
func buildDigitalOceanImports(projectName string, resource *appdef.Resource, clusterID string) ([]importAddress, error) {
	switch resource.Type {
	case appdef.ResourceTypePostgres:
		return buildPostgresImports(projectName, resource, clusterID), nil
	case appdef.ResourceTypeS3:
		return buildS3Imports(resource, clusterID), nil
	default:
		return nil, fmt.Errorf("import not supported for resource type %q", resource.Type)
	}
}

// buildPostgresImports creates the import addresses for a DigitalOcean PostgreSQL database.
// The base cluster ID is used to derive the IDs for related resources.
//
// The function builds the full resource name using the same pattern as Terraform modules:
//   - Full name: ${projectName}-${resourceName} (e.g., "search-spares-db")
//   - DB prefix: lowercase with underscores (e.g., "search_spares_db")
//
// This ensures import IDs match the naming conventions used when resources are
// provisioned through webkit's Terraform modules.
func buildPostgresImports(projectName string, resource *appdef.Resource, clusterID string) []importAddress {
	// Build the full name as Terraform does: ${project_name}-${resource_name}
	// This matches platform/terraform/providers/digital_ocean/postgres/main.tf:11
	fullName := fmt.Sprintf("%s-%s", projectName, resource.Name)

	// Convert to db_prefix (lowercase with underscores) for user/db/pool names.
	// This matches platform/terraform/providers/digital_ocean/postgres/main.tf:2
	dbPrefix := strings.ToLower(strings.ReplaceAll(fullName, "-", "_"))

	baseModule := fmt.Sprintf("module.resources[\"%s\"].module.do_postgres[0]", resource.Name)

	addresses := []importAddress{
		{
			Address: fmt.Sprintf("%s.digitalocean_database_cluster.this", baseModule),
			ID:      clusterID,
		},
		{
			Address: fmt.Sprintf("%s.digitalocean_database_user.this", baseModule),
			ID:      fmt.Sprintf("%s,%s_admin", clusterID, dbPrefix),
		},
		{
			Address: fmt.Sprintf("%s.digitalocean_database_db.this", baseModule),
			ID:      fmt.Sprintf("%s,%s", clusterID, dbPrefix),
		},
		{
			Address: fmt.Sprintf("%s.digitalocean_database_connection_pool.this", baseModule),
			ID:      fmt.Sprintf("%s,%s_pool", clusterID, dbPrefix),
		},
	}

	// Check if firewall rules are configured.
	allowedIPs, hasAllowedIPs := resource.Config["allowed_ips_addr"].([]any)
	if hasAllowedIPs && len(allowedIPs) > 0 {
		addresses = append(addresses, importAddress{
			Address: fmt.Sprintf("%s.digitalocean_database_firewall.this[0]", baseModule),
			ID:      clusterID,
		})
	}

	return addresses
}

// buildS3Imports creates the import addresses for a DigitalOcean Spaces bucket.
func buildS3Imports(resource *appdef.Resource, bucketID string) []importAddress {
	baseModule := fmt.Sprintf("module.resources[\"%s\"].module.do_bucket[0]", resource.Name)
	region, ok := resource.Config["region"].(string)
	if !ok {
		region = "ams3"
	}
	id := fmt.Sprintf("%s,%s", bucketID, region)

	return []importAddress{
		{
			Address: fmt.Sprintf("%s.digitalocean_spaces_bucket.this", baseModule),
			ID:      id,
		},
		{
			Address: fmt.Sprintf("%s.digitalocean_spaces_bucket_cors_configuration.this", baseModule),
			ID:      id,
		},
		{
			Address: fmt.Sprintf("%s.digitalocean_cdn.this", baseModule),
			ID:      id,
		},
	}
}

// buildDigitalOceanAppImports creates import addresses for DigitalOcean apps.
func buildDigitalOceanAppImports(_ string, app *appdef.App, appID string) ([]importAddress, error) {
	// app.Infra.Type is a string field in the appdef.Infra struct
	platformType := app.Infra.Type

	switch platformType {
	case "container":
		// DigitalOcean App Platform
		baseModule := fmt.Sprintf("module.apps[\"%s\"].module.do_app[0]", app.Name)
		return []importAddress{
			{
				Address: fmt.Sprintf("%s.digitalocean_app.this", baseModule),
				ID:      appID,
			},
		}, nil
	case "vm":
		// DigitalOcean Droplet
		return buildDropletImports(app, appID), nil
	default:
		return nil, fmt.Errorf("import not supported for app platform type %q", platformType)
	}
}

// buildDropletImports creates the import addresses for a DigitalOcean Droplet.
// Currently only imports the droplet resource itself. The SSH key and firewall
// resources created by the module should be imported separately if they exist.
func buildDropletImports(app *appdef.App, dropletID string) []importAddress {
	baseModule := fmt.Sprintf("module.apps[\"%s\"].module.do_droplet[0]", app.Name)

	return []importAddress{
		{
			Address: fmt.Sprintf("%s.digitalocean_droplet.this", baseModule),
			ID:      dropletID,
		},
	}
}

// buildProjectImportAddress creates the import address for a DigitalOcean project.
func buildProjectImportAddress(projectID string) []importAddress {
	return []importAddress{
		{
			Address: "digitalocean_project.this",
			ID:      projectID,
		},
	}
}
