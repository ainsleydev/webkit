package appdef

type (
	// Resource represents an infrastructure component that an application
	// depends on, such as databases, storage buckets or caches.
	Resource struct {
		Name     string               `json:"name"`
		Type     ResourceType         `json:"type"`
		Provider ResourceProvider     `json:"provider"`
		Config   map[string]any       `json:"config"` // Conforms to Terraform
		Outputs  []string             `json:"outputs"`
		Backup   ResourceBackupConfig `json:"backup,omitempty"`
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
	ResourceProviderDigitalOcean ResourceProvider = "digital-ocean"
	ResourceProviderBackBlaze    ResourceProvider = "backblaze"
)

// String implements fmt.Stringer on the ResourceProvider.
func (r ResourceProvider) String() string {
	return string(r)
}

// applyDefaults applies default values to a Resource.
func (r *Resource) applyDefaults() error {
	if r.Config == nil {
		r.Config = make(map[string]any)
	}

	if r.Outputs == nil {
		r.Outputs = []string{}
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

	return nil
}
