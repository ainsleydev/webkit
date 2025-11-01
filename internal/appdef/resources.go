package appdef

import (
	"fmt"
	"strings"

	"github.com/ainsleydev/webkit/pkg/env"
)

type (
	// Resource represents an infrastructure component that an application
	// depends on, such as databases, storage buckets or caches.
	Resource struct {
		Name             string               `json:"name"`
		Type             ResourceType         `json:"type"`
		Provider         ResourceProvider     `json:"provider"`
		Config           map[string]any       `json:"config"` // Conforms to Terraform
		Backup           ResourceBackupConfig `json:"backup,omitempty"`
		TerraformManaged *bool                `json:"terraformManaged,omitempty"`
	}
	// ResourceBackupConfig defines optional backup behavior for a resource.
	// Backup is enabled by default.
	ResourceBackupConfig struct {
		Enabled bool `json:"enabled"`
	}
)

// ResourceType defines the type of resource to be provisioned.
type ResourceType string

// ResourceType constants.
const (
	ResourceTypePostgres ResourceType = "postgres"
	ResourceTypeS3       ResourceType = "s3"
)

// String implements fmt.Stringer on the ResourceType.
func (r ResourceType) String() string {
	return string(r)
}

// ResourceProvider defines a provider of cloud infra.
type ResourceProvider string

// ResourceProvider constants.
const (
	ResourceProviderDigitalOcean ResourceProvider = "digitalocean"
	ResourceProviderBackBlaze    ResourceProvider = "b2"
)

// String implements fmt.Stringer on the ResourceProvider.
func (r ResourceProvider) String() string {
	return string(r)
}

// requiredOutputs is a global lookup of all required outputs
// from a resource type.
var requiredOutputs = map[ResourceType][]string{
	ResourceTypePostgres: {
		"id",
		"connection_url",
		"host",
		"port",
		"database",
		"user",
		"password",
	},
	ResourceTypeS3: {
		"id",
		"bucket_name",
		"bucket_url",
		"region",
	},
}

// Outputs returns the required outputs for a resource type for Terraform.
//
// These should always be exported to GitHub secrets regardless
// of user config defined in the app definition.
func (r ResourceType) Outputs() []string {
	if outputs, ok := requiredOutputs[r]; ok {
		return outputs
	}
	return nil
}

// GitHubSecretName returns the GitHub secret name for a resource output.
// Format: TF_{ENVIRONMENT}_{RESOURCE_NAME}_{OUTPUT_NAME} (uppercase)
//
// Example:
//
//	resource.GitHubSecretName(env.Production, "connection_url")
//	↓
//	"TF_PROD_DB_CONNECTION_URL"
func (r *Resource) GitHubSecretName(environment env.Environment, output string) string {
	return fmt.Sprintf("TF_%s_%s_%s",
		strings.ToUpper(environment.Short()),
		strings.ToUpper(strings.ReplaceAll(r.Name, "-", "_")),
		strings.ToUpper(output))
}

// IsTerraformManaged returns whether this resource should be managed by Terraform.
// It defaults to true when the field is nil or explicitly set to true.
func (r *Resource) IsTerraformManaged() bool {
	if r.TerraformManaged == nil {
		return true
	}
	return *r.TerraformManaged
}

// applyDefaults applies default values to a Resource.
func (r *Resource) applyDefaults() {
	if r.Config == nil {
		r.Config = make(map[string]any)
	}

	r.Backup = ResourceBackupConfig{
		Enabled: true,
	}

	// Apply type-specific defaults
	// TODO: These types should be nicely hardcoded.
	switch r.Type {
	case "postgres":
		if _, ok := r.Config["engine_version"]; !ok {
			r.Config["engine_version"] = "17"
		}
	case "s3":
		if _, ok := r.Config["acl"]; !ok {
			r.Config["acl"] = "private"
		}
	}
}
