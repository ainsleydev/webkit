package infra

import (
	"fmt"
	"strings"

	"github.com/ainsleydev/webkit/internal/appdef"
)

type (
	// importAddress represents a single Terraform resource to import.
	importAddress struct {
		// Address is the Terraform resource address (e.g., "module.resources[\"db\"].module.do_postgres[0].digitalocean_database_cluster.this").
		Address string
		// ID is the provider-specific resource identifier.
		ID string
	}
)

// buildImportAddresses constructs the list of Terraform import addresses
// for a given resource based on its type and provider.
func buildImportAddresses(resource *appdef.Resource, baseID string) ([]importAddress, error) {
	switch resource.Provider {
	case appdef.ResourceProviderDigitalOcean:
		return buildDigitalOceanImports(resource, baseID)
	default:
		return nil, fmt.Errorf("import not supported for provider %q", resource.Provider)
	}
}

// buildDigitalOceanImports creates import addresses for DigitalOcean resources.
func buildDigitalOceanImports(resource *appdef.Resource, clusterID string) ([]importAddress, error) {
	switch resource.Type {
	case appdef.ResourceTypePostgres:
		return buildPostgresImports(resource, clusterID), nil
	case appdef.ResourceTypeS3:
		return buildS3Imports(resource, clusterID), nil
	default:
		return nil, fmt.Errorf("import not supported for resource type %q", resource.Type)
	}
}

// buildPostgresImports creates the import addresses for a DigitalOcean PostgreSQL database.
// The base cluster ID is used to derive the IDs for related resources.
func buildPostgresImports(resource *appdef.Resource, clusterID string) []importAddress {
	dbPrefix := strings.ToLower(strings.ReplaceAll(resource.Name, "-", "_"))
	resourceName := fmt.Sprintf("%s-db", resource.Name)

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

	_ = resourceName // May be needed for future use.

	return addresses
}

// buildS3Imports creates the import addresses for a DigitalOcean Spaces bucket.
func buildS3Imports(resource *appdef.Resource, bucketID string) []importAddress {
	baseModule := fmt.Sprintf("module.resources[\"%s\"].module.do_bucket[0]", resource.Name)

	return []importAddress{
		{
			Address: fmt.Sprintf("%s.digitalocean_spaces_bucket.this", baseModule),
			ID:      bucketID,
		},
	}
}
