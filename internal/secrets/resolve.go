package secrets

import (
	"fmt"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
)

// ResolvedEnvVar represents an environment variable that has been
// fully resolved (decrypted if necessary) and is ready for Terraform.
type ResolvedEnvVar struct {
	Key   string
	Value string
	Type  EnvVarType
}

// EnvVarType indicates how Terraform should handle the variable
type EnvVarType string

const (
	EnvVarTypeGeneral EnvVarType = "GENERAL"
	EnvVarTypeSecret  EnvVarType = "SECRET"
)

// Resolver decrypts and resolves environment variables from the app definition
type Resolver struct {
	sopsClient sops.EncrypterDecrypter
}

// NewResolver creates a new environment variable resolver with SOPS support
func NewResolver() (*Resolver, error) {
	// Initialize age provider (checks SOPS_AGE_KEY env var or ~/.config/webkit/age.key)
	ageProvider, err := age.NewProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize age provider: %w", err)
	}

	return &Resolver{
		sopsClient: sops.NewClient(ageProvider),
	}, nil
}

// ResolveEnvironmentVariables resolves all environment variables for a given environment.
// It handles three source types:
// - "value": static string value (pass through)
// - "resource": Terraform resource reference (formatted as "resource:name.output")
// - "sops": encrypted secret (decrypted using SOPS)
func (r *Resolver) ResolveEnvironmentVariables(envVars appdef.EnvVar) (appdef.EnvVar, error) {
	resolved := make([]ResolvedEnvVar, 0, len(envVars))

	for key, config := range envVars {
		var resolvedVar ResolvedEnvVar
		var err error

		switch config.Source {
		case appdef.EnvSourceValue:
			// Static value - use as-is
			resolvedVar = ResolvedEnvVar{
				Key:   key,
				Value: config.Value,
				Type:  EnvVarTypeGeneral,
			}

		case appdef.EnvSourceResource:
			// Resource reference - format for Terraform to resolve later
			// e.g., "db.connection_url" → "resource:db.connection_url"
			resolvedVar = ResolvedEnvVar{
				Key:   key,
				Value: fmt.Sprintf("resource:%s", config.Value),
				Type:  EnvVarTypeGeneral,
			}

		case appdef.EnvSourceSOPS:
			// SOPS secret - decrypt now
			//resolvedVar, err = r.decryptSOPSSecret(key, config.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt %s: %w", key, err)
			}

		default:
			return nil, fmt.Errorf("unknown env var source type: %s for key: %s", config.Source, key)
		}

		resolved = append(resolved, resolvedVar)
	}

	return nil, nil
	//return resolved, nil
}
