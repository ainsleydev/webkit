package secrets

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
	"github.com/ainsleydev/webkit/pkg/env"
)

// TerraformOutputProvider provides access to Terraform outputs for resource resolution.
// Maps: environment -> resource name -> output name -> value
type TerraformOutputProvider struct {
	Outputs map[env.Environment]map[string]map[string]any
}

// ResolveConfig defines the data needed in order to decrypt the
// definitions environments secrets.
type ResolveConfig struct {
	SOPSClient      sops.EncrypterDecrypter
	BaseDir         string
	TerraformOutput *TerraformOutputProvider
}

func Resolve(ctx context.Context, def *appdef.Definition, cfg ResolveConfig) error {
	// Resolve shared environment
	if err := resolveAllEnvs(ctx, cfg, &def.Shared.Env); err != nil {
		return fmt.Errorf("resolving shared env: %w", err)
	}

	// Resolve each app environment
	for i := range def.Apps {
		if err := resolveAllEnvs(ctx, cfg, &def.Apps[i].Env); err != nil {
			return fmt.Errorf("resolving app %q env: %w", def.Apps[i].Name, err)
		}
	}

	return nil
}

// resolveAllEnvs resolves all variables in an Environment (dev, staging, production)
func resolveAllEnvs(ctx context.Context, cfg ResolveConfig, enviro *appdef.Environment) error {
	return enviro.WalkE(func(entry appdef.EnvWalkEntry) error {
		for key, config := range entry.Map {
			resolveFn, ok := resolver[config.Source]
			if !ok {
				return fmt.Errorf("unknown env source type: %s", config.Source)
			}

			rc := resolveContext{
				cfg:    cfg,
				env:    entry.Environment,
				key:    key,
				config: config,
				vars:   entry.Map,
			}

			if err := resolveFn(ctx, rc); err != nil {
				return err
			}
		}
		return nil
	})
}

type resolveContext struct {
	cfg    ResolveConfig
	env    env.Environment
	key    string
	config appdef.EnvValue
	vars   appdef.EnvVar
}

type resolveFunc func(ctx context.Context, rc resolveContext) error

var resolver = map[appdef.EnvSource]resolveFunc{
	// Static value - use as-is
	appdef.EnvSourceValue: func(ctx context.Context, rc resolveContext) error {
		return nil
	},
	// Resource reference - resolve from Terraform outputs
	appdef.EnvSourceResource: func(_ context.Context, rc resolveContext) error {
		// Parse the resource reference (e.g., "db.connection_url")
		resourceName, outputName, ok := appdef.ParseResourceReference(rc.config.Value)
		if !ok {
			return fmt.Errorf("invalid resource reference format for key '%s': expected 'resource_name.output_name', got '%v'", rc.key, rc.config.Value)
		}

		// Check if Terraform outputs are provided
		if rc.cfg.TerraformOutput == nil {
			return fmt.Errorf("terraform outputs not provided: cannot resolve resource reference '%s.%s' for key '%s'", resourceName, outputName, rc.key)
		}

		// Get outputs for the current environment
		envOutputs, ok := rc.cfg.TerraformOutput.Outputs[rc.env]
		if !ok {
			return fmt.Errorf("no terraform outputs found for environment '%s' (referenced by key '%s')", rc.env, rc.key)
		}

		// Get outputs for the specific resource
		resourceOutputs, ok := envOutputs[resourceName]
		if !ok {
			return fmt.Errorf("resource '%s' not found in terraform outputs for environment '%s' (referenced by key '%s')", resourceName, rc.env, rc.key)
		}

		// Get the specific output value
		value, ok := resourceOutputs[outputName]
		if !ok {
			return fmt.Errorf("output '%s' not found for resource '%s' in terraform outputs (referenced by key '%s')", outputName, resourceName, rc.key)
		}

		// Update the variable with the resolved value
		rc.vars[rc.key] = appdef.EnvValue{
			Source: rc.config.Source,
			Value:  value,
		}

		return nil
	},
	// SOPS secret - decrypt now
	appdef.EnvSourceSOPS: func(_ context.Context, rc resolveContext) error {
		path := filepath.Join(rc.cfg.BaseDir, FilePathFromEnv(rc.env))

		// If it's an internal error and not decrypted, bail early
		// as we can't resolve much!
		resolvedMap, err := sops.DecryptFileToMap(rc.cfg.SOPSClient, path)
		if err != nil && !errors.Is(err, sops.ErrNotEncrypted) {
			return err
		}

		secret, ok := resolvedMap[rc.key]
		if !ok {
			return fmt.Errorf("secret '%s' not found", rc.key)
		}

		rc.vars[rc.key] = appdef.EnvValue{
			Source: rc.config.Source,
			Value:  secret,
		}

		return nil
	},
}
