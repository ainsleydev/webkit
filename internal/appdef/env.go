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

// EnvironmentWalker defines the type for walking a collection of
// environment variables.
type EnvironmentWalker func(env env.Environment, name string, value EnvValue)

// Walk walks through each non-nil environment (dev, staging, production),
// calling fn(envName, envVars) for each one.
func (e Environment) Walk(fn EnvironmentWalker) {
	if e.Dev != nil {
		for name, val := range e.Dev {
			fn(env.Development, name, val)
		}
	}
	if e.Staging != nil {
		for name, val := range e.Staging {
			fn(env.Staging, name, val)
		}
	}
	if e.Production != nil {
		for name, val := range e.Production {
			fn(env.Production, name, val)
		}
	}
}

// EnvironmentWalkerContext contains the context for walking environment variables
type EnvironmentWalkerContext struct {
	Env   string
	Name  string
	Value EnvValue
	Vars  EnvVar // The map being iterated over
}

// EnvironmentWalker defines the type for walking a collection of
// environment variables.
//type EnvironmentWalker func(ctx EnvironmentWalkerContext) error
//
//// Walk walks through each non-nil environment (dev, staging, production),
//// calling fn(envName, envVars) for each one.
//func (e Environment) Walk(fn EnvironmentWalker) error {
//	if e.Dev != nil {
//		for name, val := range e.Dev {
//			if err := fn(EnvironmentWalkerContext{
//				Env:   env.Development,
//				Name:  name,
//				Value: val,
//				Vars:  e.Dev,
//			}); err != nil {
//				return err
//			}
//		}
//	}
//	if e.Staging != nil {
//		for name, val := range e.Staging {
//			if err := fn(EnvironmentWalkerContext{
//				Env:   env.Staging,
//				Name:  name,
//				Value: val,
//				Vars:  e.Staging,
//			}); err != nil {
//				return err
//			}
//		}
//	}
//	if e.Production != nil {
//		for name, val := range e.Production {
//			if err := fn(EnvironmentWalkerContext{
//				Env:   env.Production,
//				Name:  name,
//				Value: val,
//				Vars:  e.Production,
//			}); err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}

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
