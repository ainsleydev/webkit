package appdef

import (
	"github.com/ainsleydev/webkit/pkg/env"
)

type (
	// Environment contains environment-specific variable configurations.
	Environment struct {
		Dev        EnvVar `json:"dev,omitempty"`
		Staging    EnvVar `json:"staging,omitempty"`
		Production EnvVar `json:"production,omitempty"`
	}
	// EnvVar is a map of variable names to their configurations.
	EnvVar map[string]EnvValue
	// EnvValue represents a single environment variable configuration
	EnvValue struct {
		Source EnvSource `json:"source"`          // See below
		Value  any       `json:"value,omitempty"` // Used for "value" and "resource" sources
		Path   string    `json:"path,omitempty"`  // Used for "sops" source (format: "file:key")
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

// EnvWalkEntry holds the details of a single environment variable during iteration.
type EnvWalkEntry struct {
	Environment env.Environment
	Key         string
	Value       any
	Source      EnvSource
	Path        string
	Map         EnvVar
}

// EnvironmentWalker defines a function that processes one environment entry.
type EnvironmentWalker func(entry EnvWalkEntry)

// EnvironmentWalkerE defines a function that processes one environment entry and may return an error.
type EnvironmentWalkerE func(entry EnvWalkEntry) error

// Walk iterates over all environments and calls fn for each environment variable.
func (e Environment) Walk(fn EnvironmentWalker) {
	_ = e.walkEnvs(func(envName env.Environment, vars EnvVar) error {
		for key, val := range vars {
			fn(EnvWalkEntry{
				Environment: envName,
				Key:         key,
				Value:       val.Value,
				Source:      val.Source,
				Path:        val.Path,
				Map:         vars,
			})
		}
		return nil
	})
}

// WalkE iterates over all environments and calls fn for each environment variable.
// If fn returns an error, iteration stops and the error is returned.
func (e Environment) WalkE(fn EnvironmentWalkerE) error {
	return e.walkEnvs(func(envName env.Environment, vars EnvVar) error {
		for key, val := range vars {
			if err := fn(EnvWalkEntry{
				Environment: envName,
				Key:         key,
				Value:       val.Value,
				Source:      val.Source,
				Path:        val.Path,
				Map:         vars,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

// walkEnvs is the internal helper that iterates over environments.
func (e Environment) walkEnvs(fn func(envName env.Environment, vars EnvVar) error) error {
	envs := []struct {
		name env.Environment
		vars EnvVar
	}{
		{env.Development, e.Dev},
		{env.Staging, e.Staging},
		{env.Production, e.Production},
	}

	for _, envData := range envs {
		if envData.vars == nil {
			continue
		}
		if err := fn(envData.name, envData.vars); err != nil {
			return err
		}
	}

	return nil
}

// mergeVars merges `override` into `base`, with `override`
// taking precedence (usually app/shared).
func mergeVars(base, override EnvVar) EnvVar {
	if base == nil {
		base = make(EnvVar)
	}
	if override == nil {
		override = make(EnvVar)
	}
	for k, v := range override {
		base[k] = v
	}
	return base
}
