package appdef

import (
	"fmt"
	"strings"

	"github.com/ainsleydev/webkit/pkg/env"
)

type (
	// Environment contains environment-specific variable configurations.
	// Variables can be defined per environment (dev, staging, production)
	// or set once in 'default' to apply across all environments.
	Environment struct {
		Default    EnvVar `json:"default,omitempty" inline:"true" description:"Environment variables that apply to all environments (dev, staging, production)"`
		Dev        EnvVar `json:"dev,omitempty" inline:"true" description:"Environment variables specific to the development environment"`
		Staging    EnvVar `json:"staging,omitempty" inline:"true" description:"Environment variables specific to the staging environment"`
		Production EnvVar `json:"production,omitempty" inline:"true" description:"Environment variables specific to the production environment"`
	}
	// EnvVar is a map of variable names to their value configurations.
	// Each key is the environment variable name, and the value defines
	// where the variable's value comes from (static value, resource output, or secret).
	EnvVar map[string]EnvValue
	// EnvValue represents a single environment variable configuration.
	// It specifies both the source type and the value/reference for the variable.
	EnvValue struct {
		Source EnvSource `json:"source" jsonschema:"required" description:"Source type for the variable value (value, resource, sops)"`
		// Value holds the actual value or reference depending on the source type:
		// - "value": A static string (e.g., "https://api.example.com")
		// - "resource": A Terraform resource reference (e.g., "db.connection_url")
		// - "sops": The variable name/key to lookup in the SOPS file (e.g., "API_KEY")
		Value any `json:"value,omitempty" description:"The value or reference for this variable (format depends on source type)"`
	}
)

// EnvSource defines the type of application being run.
type EnvSource string

const (
	// EnvSourceValue is a static string value (default).
	// Example: "https://api.example.com"
	EnvSourceValue EnvSource = "value"

	// EnvSourceResource references a Terraform resource output.
	// Example: "db.connection_url"
	EnvSourceResource EnvSource = "resource"

	// EnvSourceSOPS is an encrypted secret stored in a SOPS file.
	// Example: "secrets/production.yaml:API_KEY"
	EnvSourceSOPS EnvSource = "sops"
)

// String implements fmt.Stringer on the EnvSource.
func (e EnvSource) String() string {
	return string(e)
}

// EnvWalkEntry holds the details of a single env variable during iteration.
type EnvWalkEntry struct {
	Environment env.Environment
	Key         string
	Value       any
	Source      EnvSource
	Map         EnvVar
}

// EnvironmentWalker defines a function that processes one env entry.
type EnvironmentWalker func(entry EnvWalkEntry)

// EnvironmentWalkerE defines a function that processes one env entry and may return an error.
type EnvironmentWalkerE func(entry EnvWalkEntry) error

// Walk iterates over all environments and calls fn for each env variable.
func (e Environment) Walk(fn EnvironmentWalker) {
	_ = e.walkEnvs(func(envName env.Environment, vars EnvVar) error {
		for key, val := range vars {
			fn(EnvWalkEntry{
				Environment: envName,
				Key:         key,
				Value:       val.Value,
				Source:      val.Source,
				Map:         vars,
			})
		}
		return nil
	})
}

// WalkE iterates over all environments and calls fn for each env variable.
// If fn returns an error, iteration stops and the error is returned.
func (e Environment) WalkE(fn EnvironmentWalkerE) error {
	return e.walkEnvs(func(envName env.Environment, vars EnvVar) error {
		for key, val := range vars {
			if err := fn(EnvWalkEntry{
				Environment: envName,
				Key:         key,
				Value:       val.Value,
				Source:      val.Source,
				Map:         vars,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

// walkEnvs is the internal helper that iterates over environments.
// It walks over Default vars first (for each environment), then environment-specific vars.
// This ensures defaults apply to all environments, with environment-specific values overriding.
// The original maps are passed to the walker, allowing mutations.
func (e Environment) walkEnvs(fn func(envName env.Environment, vars EnvVar) error) error {
	// First, walk over default vars for each environment.
	if len(e.Default) > 0 {
		for _, envName := range []env.Environment{env.Development, env.Staging, env.Production} {
			if err := fn(envName, e.Default); err != nil {
				return err
			}
		}
	}

	// Then walk over environment-specific vars.
	envs := []struct {
		name env.Environment
		vars EnvVar
	}{
		{env.Development, e.Dev},
		{env.Staging, e.Staging},
		{env.Production, e.Production},
	}

	for _, envData := range envs {
		if envData.vars == nil || len(envData.vars) == 0 {
			continue
		}
		if err := fn(envData.name, envData.vars); err != nil {
			return err
		}
	}

	return nil
}

// GetVarsForEnvironment returns the EnvVar map for the specified environment.
// Returns an error if the environment is unknown.
func (e Environment) GetVarsForEnvironment(target env.Environment) (EnvVar, error) {
	switch target {
	case env.Development:
		return e.Dev, nil
	case env.Staging:
		return e.Staging, nil
	case env.Production:
		return e.Production, nil
	default:
		return nil, fmt.Errorf("unknown environment: %s", target)
	}
}

// ParseResourceReference parses a resource reference string
// (e.g., "db.connection_url").
//
// Resource references follow the format: "resource_name.output_name".
func ParseResourceReference(value any) (resourceName, outputName string, ok bool) {
	valueStr, isString := value.(string)
	if !isString {
		return "", "", false
	}

	parts := strings.SplitN(valueStr, ".", 2)
	if len(parts) != 2 || (parts[0] == "" || parts[1] == "") {
		return "", "", false
	}

	return parts[0], parts[1], true
}

// mergeVars merges `override` into `base`, with `override`
// taking precedence (usually app/shared).
// Returns a new map without mutating the inputs.
func mergeVars(base, override EnvVar) EnvVar {
	result := make(EnvVar)

	// Copy base first.
	for k, v := range base {
		result[k] = v
	}

	// Apply overrides.
	for k, v := range override {
		result[k] = v
	}

	return result
}
