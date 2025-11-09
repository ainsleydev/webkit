package appdef

import (
	"fmt"
	"strings"

	"github.com/ainsleydev/webkit/pkg/env"
)

type (
	// Resource represents an infrastructure component that applications
	// depend on, such as databases, storage buckets, or caches.
	// Resources are provisioned via Terraform and their outputs are
	// made available to apps through environment variables.
	Resource struct {
		Name             string               `json:"name" required:"true" validate:"required,lowercase,alphanumdash" description:"Unique identifier for the resource (used in environment variable references)"`
		Type             ResourceType         `json:"type" required:"true" validate:"required,oneof=postgres s3" description:"Type of resource to provision (postgres, s3)"`
		Provider         ResourceProvider     `json:"provider" required:"true" validate:"required,oneof=digitalocean backblaze" description:"Cloud provider hosting this resource (digitalocean, backblaze)"`
		Config           map[string]any       `json:"config" description:"Provider-specific resource configuration (e.g., size, region, version)"`
		Backup           ResourceBackupConfig `json:"backup,omitempty" description:"Backup configuration for the resource"`
		TerraformManaged *bool                `json:"terraformManaged,omitempty" description:"Whether this resource is managed by Terraform (defaults to true)"`
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
	ResourceProviderBackBlaze    ResourceProvider = "backblaze"
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
