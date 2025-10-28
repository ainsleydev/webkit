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
	// SkippedItems contains information about apps and resources
	// that were not included in Terraform variables.
	SkippedItems struct {
		Apps      []string
		Resources []string
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
//
// Returns the generated Terraform variables, information about skipped items,
// and an error if any occurred.
func tfVarsFromDefinition(env env.Environment, def *appdef.Definition) (tfVars, SkippedItems, error) {
	if def == nil {
		return tfVars{}, SkippedItems{}, errors.New("definition cannot be nil")
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

	skipped := SkippedItems{
		Apps:      make([]string, 0),
		Resources: make([]string, 0),
	}

	for _, res := range def.Resources {
		// Skip resources that are not managed by Terraform.
		if !res.IsTerraformManaged() {
			skipped.Resources = append(skipped.Resources, res.Name)
			continue
		}

		vars.Resources = append(vars.Resources, tfResource{
			Name:             res.Name,
			PlatformType:     res.Type.String(),
			PlatformProvider: res.Provider.String(),
			Config:           res.Config,
		})
	}

	for _, app := range def.Apps {
		// Skip apps that are not managed by Terraform.
		if !app.IsTerraformManaged() {
			skipped.Apps = append(skipped.Apps, app.Name)
			continue
		}

		tfA := tfApp{
			Name:             app.Name,
			PlatformType:     app.Infra.Type,
			PlatformProvider: app.Infra.Provider.String(),
			AppType:          app.Type.String(),
			Config:           app.Infra.Config,
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

	return vars, skipped, nil
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
