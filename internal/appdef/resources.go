package appdef

import (
	"fmt"
	"strings"

	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

type (
	// Resource represents an infrastructure component that applications
	// depend on, such as databases, storage buckets, or caches.
	// Resources are provisioned via Terraform and their outputs are
	// made available to apps through environment variables.
	Resource struct {
		Name             string                `json:"name" required:"true" validate:"required,lowercase,alphanumdash" description:"Unique identifier for the resource (used in environment variable references)"`
		Title            string                `json:"title" required:"true" validate:"required" description:"Human-readable resource name for display purposes"`
		Type             ResourceType          `json:"type" required:"true" validate:"required,oneof=postgres s3 sqlite" description:"Type of resource to provision (postgres, s3, sqlite)"`
		Description      string                `json:"description,omitempty" validate:"omitempty,max=200" description:"Brief description of the resource's purpose and functionality"`
		Provider         ResourceProvider      `json:"provider" required:"true" validate:"required,oneof=digitalocean hetzner backblaze turso" description:"Cloud provider hosting this resource (digitalocean, hetzner, backblaze, turso)"`
		Config           Config                `json:"config" description:"Provider-specific resource configuration (e.g., size, region, version)"`
		Backup           *ResourceBackupConfig `json:"backup,omitempty" description:"Backup configuration for the resource"`
		Monitoring       *bool                 `json:"monitoring,omitempty" description:"Whether to enable uptime monitoring for this resource (defaults to true)"`
		TerraformManaged *bool                 `json:"terraformManaged,omitempty" description:"Whether this resource is managed by Terraform (defaults to true)"`
	}
	// ResourceBackupConfig defines backup behaviour for a resource.
	// Backups are enabled by default for all resources that support them.
	ResourceBackupConfig struct {
		Enabled bool `json:"enabled" description:"Whether to enable automated backups for this resource"`
	}
)

// ResourceType defines the type of resource to be provisioned.
type ResourceType string

// ResourceType constants.
const (
	ResourceTypePostgres ResourceType = "postgres"
	ResourceTypeS3       ResourceType = "s3"
	ResourceTypeSQLite   ResourceType = "sqlite"
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
	ResourceProviderHetzner      ResourceProvider = "hetzner"
	ResourceProviderBackBlaze    ResourceProvider = "backblaze"
	ResourceProviderTurso        ResourceProvider = "turso"
)

// String implements fmt.Stringer on the ResourceProvider.
func (r ResourceProvider) String() string {
	return string(r)
}

// ResourceOutput represents a single resource output with documentation.
type ResourceOutput struct {
	Name        string
	Description string
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
	ResourceTypeSQLite: {
		"id",
		"connection_url",
		"host",
		"database",
		"auth_token",
	},
}

// resourceOutputDocumentation is a global lookup of output documentation
// for resource types.
var resourceOutputDocumentation = map[ResourceType][]ResourceOutput{
	ResourceTypePostgres: {
		{Name: "connection_url", Description: "Full PostgreSQL connection string"},
		{Name: "host", Description: "Database host address"},
		{Name: "port", Description: "Database port number"},
		{Name: "database", Description: "Database name"},
		{Name: "user", Description: "Database username"},
		{Name: "password", Description: "Database password"},
	},
	ResourceTypeS3: {
		{Name: "bucket_name", Description: "S3 bucket identifier"},
		{Name: "bucket_url", Description: "Public bucket URL"},
		{Name: "region", Description: "Bucket region"},
	},
	ResourceTypeSQLite: {
		{Name: "connection_url", Description: "Full SQLite connection string"},
		{Name: "host", Description: "Turso database host"},
		{Name: "database", Description: "Database name"},
		{Name: "auth_token", Description: "Authentication token"},
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

	if r.Backup == nil {
		r.Backup = &ResourceBackupConfig{
			Enabled: true,
		}
	}

	if r.Monitoring == nil {
		r.Monitoring = ptr.BoolPtr(true)
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
	case "sqlite":
		if _, ok := r.Config["group"]; !ok {
			r.Config["group"] = "default"
		}
	}
}

// Documentation returns the available outputs for a resource type
// with descriptions, this is directly tied to Terraform outputs.
func (r ResourceType) Documentation() []ResourceOutput {
	if out, ok := resourceOutputDocumentation[r]; ok {
		return out
	}
	return []ResourceOutput{}
}

// GenerateBackupMonitor creates a push monitor for a resource's backup workflow.
// It only generates a monitor if both backup and monitoring are enabled for the resource.
// The monitor name follows the format: "{ProjectTitle} - {ResourceTitle} Backup".
// This creates a heartbeat monitor that can be pinged by CI/CD backup workflows.
func (r *Resource) GenerateBackupMonitor(projectTitle string) *Monitor {
	if r.Backup == nil || !r.Backup.Enabled {
		return nil
	}
	// Only skip if monitoring is explicitly disabled
	if r.Monitoring != nil && !*r.Monitoring {
		return nil
	}

	return &Monitor{
		Name: fmt.Sprintf("%s - %s Backup", projectTitle, r.Title),
		Type: MonitorTypePush,
	}
}
