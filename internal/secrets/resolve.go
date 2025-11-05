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

// OutputKey uniquely identifies a Terraform output value.
type OutputKey struct {
	Environment  env.Environment
	ResourceName string
	OutputName   string
}

// TerraformOutputProvider provides access to Terraform outputs for resource resolution.
type TerraformOutputProvider map[OutputKey]any

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

// ResolveForEnvironment resolves variables for a specific environment only.
// This is more efficient when you only need one environment (e.g., env generation).
func ResolveForEnvironment(ctx context.Context, def *appdef.Definition, targetEnv env.Environment, cfg ResolveConfig) error {
	// Resolve shared environment for target env
	if err := resolveSingleEnv(ctx, cfg, &def.Shared.Env, targetEnv); err != nil {
		return fmt.Errorf("resolving shared env: %w", err)
	}

	// Resolve each app environment for target env
	for i := range def.Apps {
		if err := resolveSingleEnv(ctx, cfg, &def.Apps[i].Env, targetEnv); err != nil {
			return fmt.Errorf("resolving app %q env: %w", def.Apps[i].Name, err)
		}
	}

	return nil
}

// resolveAllEnvs resolves all variables in an Environment (dev, staging, production)
func resolveAllEnvs(ctx context.Context, cfg ResolveConfig, enviro *appdef.Environment) error {
	// Resolve all three environments by calling resolveSingleEnv for each
	for _, targetEnv := range []env.Environment{env.Development, env.Staging, env.Production} {
		if err := resolveSingleEnv(ctx, cfg, enviro, targetEnv); err != nil {
			return err
		}
	}
	return nil
}

// resolveSingleEnv resolves variables for a specific environment only.
// It resolves defaults first, then environment-specific vars (following the merge pattern).
func resolveSingleEnv(ctx context.Context, cfg ResolveConfig, enviro *appdef.Environment, targetEnv env.Environment) error {
	// Resolve defaults first (they apply to the target environment)
	if err := resolveVars(ctx, cfg, enviro.Default, targetEnv); err != nil {
		return err
	}

	// Get the specific environment vars
	targetVars, err := enviro.GetVarsForEnvironment(targetEnv)
	if err != nil {
		return err
	}

	// Resolve environment-specific vars
	return resolveVars(ctx, cfg, targetVars, targetEnv)
}

// resolveVars resolves all variables in a single EnvVar map.
func resolveVars(ctx context.Context, cfg ResolveConfig, vars appdef.EnvVar, targetEnv env.Environment) error {
	for key, config := range vars {
		resolveFn, ok := resolver[config.Source]
		if !ok {
			return fmt.Errorf("unknown env source type: %s", config.Source)
		}

		rc := resolveContext{
			cfg:    cfg,
			env:    targetEnv,
			key:    key,
			config: config,
			vars:   vars,
		}

		if err := resolveFn(ctx, rc); err != nil {
			return err
		}
	}
	return nil
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
	// Static value - use as-is.
	appdef.EnvSourceValue: func(ctx context.Context, rc resolveContext) error {
		return nil
	},
	// Resource reference - resolves infrastructure outputs by querying Terraform state.
	// This allows environment variables to reference resources defined in app.json (e.g., "db.connection_url").
	appdef.EnvSourceResource: func(_ context.Context, rc resolveContext) error {
		// Don't resolve anything that's not production now.
		if rc.env != env.Production {
			return nil
		}

		resourceName, outputName, ok := appdef.ParseResourceReference(rc.config.Value)
		if !ok {
			return fmt.Errorf("invalid resource reference format for key '%s': expected 'resource_name.output_name', got '%v'", rc.key, rc.config.Value)
		}

		if rc.cfg.TerraformOutput == nil {
			return fmt.Errorf("terraform outputs not provided: cannot resolve resource reference '%s.%s' for key '%s'", resourceName, outputName, rc.key)
		}

		key := OutputKey{
			Environment:  rc.env,
			ResourceName: resourceName,
			OutputName:   outputName,
		}

		value, ok := (*rc.cfg.TerraformOutput)[key]
		if !ok {
			return fmt.Errorf("terraform output not found for environment '%s', resource '%s', output '%s' (referenced by key '%s')", rc.env, resourceName, outputName, rc.key)
		}

		rc.vars[rc.key] = appdef.EnvValue{
			Source: rc.config.Source,
			Value:  value,
		}

		return nil
	},
	// SOPS secret - decrypt now.
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
