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

// ResolveConfig defines the data needed in order to decrypt the
// definitions environments secrets.
type ResolveConfig struct {
	SOPSClient sops.EncrypterDecrypter
	BaseDir    string
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
	// We haven't implemented this yet, but hopefully we can do
	// so when we get Terraform outputs.
	appdef.EnvSourceResource: func(_ context.Context, rc resolveContext) error {
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
			Path:   rc.config.Path,
		}

		return nil
	},
}
