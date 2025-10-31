package infra

import (
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
		Config           map[string]any `json:"config"`
		Environment      []tfEnvVar     `json:"environment,omitempty"`
	}
	// tfEnvVar represents an environment variable for Terraform
	tfEnvVar struct {
		Key    string `json:"key"`
		Value  any    `json:"value,omitempty"`
		Source string `json:"source,omitempty"`
		Scope  string `json:"scope,omitempty"`
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
func tfVarsFromDefinition(env env.Environment, def *appdef.Definition) (tfVars, error) {
	if def == nil {
		return tfVars{}, errors.New("definition cannot be nil")
	}

	vars := tfVars{
		ProjectName: def.Project.Name,
		Environment: env.String(),
		Apps:        make([]tfApp, 0, len(def.Apps)),
		Resources:   make([]tfResource, 0, len(def.Resources)),
		GithubConfig: tfGithubConfig{
			Owner: def.Project.Repo.Owner,
			Repo:  def.Project.Repo.Name,
		},
	}

	for _, res := range def.Resources {
		vars.Resources = append(vars.Resources, tfResource{
			Name:             res.Name,
			PlatformType:     res.Type.String(),
			PlatformProvider: res.Provider.String(),
			Config:           normalizeConfig(res.Config),
		})
	}

	for _, app := range def.Apps {
		tfA := tfApp{
			Name:             app.Name,
			PlatformType:     app.Infra.Type,
			PlatformProvider: app.Infra.Provider.String(),
			AppType:          app.Type.String(),
			Config:           normalizeConfig(app.Infra.Config),
			Path:             app.Path,
		}

		app.MergeEnvironments(def.Shared.Env).
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

		vars.Apps = append(vars.Apps, tfA)
	}

	return vars, nil
}

// normalizeConfig recursively normalizes config values to ensure proper typing for Terraform.
//
// This is necessary because JSON unmarshaling creates []any for arrays,
// which can cause Terraform type inference issues when passed as
// the 'any' type.
func normalizeConfig(config map[string]any) map[string]any {
	if config == nil {
		return nil
	}

	normalized := make(map[string]any, len(config))
	for key, value := range config {
		normalized[key] = normalizeValue(value)
	}

	return normalized
}

// normalizeValue recursively normalizes a single config value.
func normalizeValue(value any) any {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []any:
		// Convert []interface{} to a properly typed slice.
		return normalizeSlice(v)
	case map[string]any:
		// Recursively normalize nested maps.
		return normalizeConfig(v)
	default:
		// Return primitive types as-is.
		return value
	}
}

// normalizeSlice converts []interface{} to a typed slice based on the element types.
func normalizeSlice(slice []interface{}) any {
	if len(slice) == 0 {
		return []string{} // Empty slices default to []string for Terraform compatibility.
	}

	// Check if all elements are strings.
	allStrings := true
	for _, elem := range slice {
		if _, ok := elem.(string); !ok {
			allStrings = false
			break
		}
	}

	if allStrings {
		result := make([]string, len(slice))
		for i, elem := range slice {
			result[i] = elem.(string)
		}
		return result
	}

	// Check if all elements are numbers (float64 in JSON).
	allNumbers := true
	for _, elem := range slice {
		if _, ok := elem.(float64); !ok {
			allNumbers = false
			break
		}
	}

	if allNumbers {
		result := make([]float64, len(slice))
		for i, elem := range slice {
			result[i] = elem.(float64)
		}
		return result
	}

	// Check if all elements are bools.
	allBools := true
	for _, elem := range slice {
		if _, ok := elem.(bool); !ok {
			allBools = false
			break
		}
	}

	if allBools {
		result := make([]bool, len(slice))
		for i, elem := range slice {
			result[i] = elem.(bool)
		}
		return result
	}

	// For mixed or complex types, recursively normalize each element.
	result := make([]any, len(slice))
	for i, elem := range slice {
		result[i] = normalizeValue(elem)
	}

	return result
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
