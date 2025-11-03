package infra

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/pkg/env"
)

// Add at top of file:
const (
	envScopeSecret  = "SECRET"
	envScopeGeneral = "GENERAL"
)

type (
	// tfVars represents the root structure of Terraform variables
	// that will be written to webkit.auto.tfvars.json
	tfVars struct {
		ProjectName  string         `json:"project_name"`
		Environment  string         `json:"environment"`
		GithubConfig tfGithubConfig `json:"github_config"`
		Apps         []tfApp        `json:"apps"`
		Resources    []tfResource   `json:"resources"`
	}
	// tfResource represents a resource in Terraform variable format.
	tfResource struct {
		Name             string         `json:"name"`
		PlatformType     string         `json:"platform_type"`
		PlatformProvider string         `json:"platform_provider"`
		Config           map[string]any `json:"config"`
	}
	// tfApp represents an application in Terraform variable format.
	tfApp struct {
		Name             string         `json:"name"`
		PlatformType     string         `json:"platform_type"`
		PlatformProvider string         `json:"platform_provider"`
		AppType          string         `json:"app_type"`
		Path             string         `json:"path"`
		ImageTag         string         `json:"image_tag,omitempty"`
		Config           map[string]any `json:"config"`
		Environment      []tfEnvVar     `json:"env_vars,omitempty"`
		Domains          []tfDomain     `json:"domains,omitempty"`
	}
	// tfDomain represents a domain configuration for Terraform.
	tfDomain struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Zone     string `json:"zone,omitempty"`
		Wildcard bool   `json:"wildcard,omitempty"`
	}
	// tfEnvVar represents an environment variable for Terraform
	tfEnvVar struct {
		Key    string `json:"key"`
		Value  any    `json:"value,omitempty"`
		Source string `json:"source,omitempty"`
		Scope  string `json:"type,omitempty"`
	}
	// tfGithubConfig is used to pull image containers from GH
	// container registry.
	tfGithubConfig struct {
		Owner string `json:"owner"`
		Repo  string `json:"repo"`
	}
)

// tfVarsFromDefinition transforms an appdef.Definition into Terraform variables.
// It should directly map what's defined in /platform/base/variables.tf, if it
// doesn't, then provisioning probably won't work.
func (t *Terraform) tfVarsFromDefinition(ctx context.Context, env env.Environment) (tfVars, error) {
	if t.appDef == nil {
		return tfVars{}, errors.New("definition cannot be nil")
	}

	vars := tfVars{
		ProjectName: t.appDef.Project.Name,
		Environment: env.String(),
		Apps:        make([]tfApp, 0, len(t.appDef.Apps)),
		Resources:   make([]tfResource, 0, len(t.appDef.Resources)),
		GithubConfig: tfGithubConfig{
			Owner: t.appDef.Project.Repo.Owner,
			Repo:  t.appDef.Project.Repo.Name,
		},
	}

	for _, res := range t.appDef.Resources {
		vars.Resources = append(vars.Resources, tfResource{
			Name:             res.Name,
			PlatformType:     res.Type.String(),
			PlatformProvider: res.Provider.String(),
			Config:           encodeConfigForTerraform(res.Config),
		})
	}

	for _, app := range t.appDef.Apps {
		tfA := tfApp{
			Name:             app.Name,
			PlatformType:     app.Infra.Type,
			PlatformProvider: app.Infra.Provider.String(),
			AppType:          app.Type.String(),
			Config:           encodeConfigForTerraform(app.Infra.Config),
			Path:             app.Path,
		}

		// Determine the image tag for container-based apps.
		if app.Infra.Type == "container" {
			tfA.ImageTag = t.determineImageTag(ctx, app.Name)
		}

		app.MergeEnvironments(t.appDef.Shared.Env).
			Walk(func(entry appdef.EnvWalkEntry) {
				if entry.Environment != env {
					return
				}
				scope := envScopeSecret
				if entry.Source == appdef.EnvSourceValue {
					scope = envScopeGeneral
				}
				tfA.Environment = append(tfA.Environment, tfEnvVar{
					Key:    entry.Key,
					Value:  entry.Value,
					Source: entry.Source.String(),
					Scope:  scope,
				})
			})

		for _, domain := range app.Domains {
			tfA.Domains = append(tfA.Domains, tfDomain{
				Name:     domain.Name,
				Type:     domain.Type.String(),
				Zone:     domain.Zone,
				Wildcard: domain.Wildcard,
			})
		}

		vars.Apps = append(vars.Apps, tfA)
	}

	return vars, nil
}

// encodeConfigForTerraform prepares configuration maps for Terraform consumption.
//
// This function solves two critical issues with Terraform's type system:
//
//  1. Type Consistency: Terraform requires all elements in a list to have the same type.
//     When some configs are nil and others are {}, Terraform sees incompatible types.
//     Solution: Convert nil configs to empty maps for consistency.
//
//  2. Array Encoding: Terraform's 'any' type has trouble with heterogeneous structures
//     when resources have different config shapes (e.g., one has arrays, another doesn't).
//     The webkit Terraform modules work around this by expecting arrays as JSON strings
//     and using jsondecode() to parse them back.
//
//     Example in platform/terraform/modules/resources/main.tf:
//     allowed_ips_addr = try(jsondecode(var.platform_config.allowed_ips_addr), [])
//
//     This pattern allows all config values to be strings/primitives, avoiding type
//     inference issues across list elements.
//
// Returns an empty map for nil input to ensure type consistency.
func encodeConfigForTerraform(config map[string]any) map[string]any {
	if config == nil {
		return map[string]any{}
	}

	encoded := make(map[string]any, len(config))
	for key, value := range config {
		encoded[key] = encodeConfigValue(value)
	}
	return encoded
}

// encodeConfigValue encodes a single config value for Terraform.
// Arrays/slices are JSON-encoded as strings. Other types pass through unchanged.
func encodeConfigValue(value any) any {
	if value == nil {
		return nil
	}

	// Check if the value is a slice/array (regardless of element type).
	switch v := value.(type) {
	case []any, []string, []int, []float64, []bool:
		// Encode arrays as JSON strings for Terraform's jsondecode().
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			// If marshaling fails, return as-is (shouldn't happen with standard types).
			return value
		}
		return string(jsonBytes)
	default:
		// Primitives (strings, numbers, bools) and other types pass through.
		return value
	}
}

// writeTFVarsFile writes Terraform variables to a JSON file.
// Terraform automatically loads *.auto.tfvars.json files.
func (t *Terraform) writeTFVarsFile(vars tfVars) error {
	data, err := json.MarshalIndent(vars, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal tf vars")
	}

	path := filepath.Join(filepath.Join(t.tmpDir, "base"), "webkit.auto.tfvars.json")
	if err = afero.WriteFile(t.fs, path, data, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to write tf vars file")
	}

	return nil
}
