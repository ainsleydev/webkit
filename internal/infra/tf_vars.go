package infra

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/pkg/env"
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
		Outputs          []string       `json:"outputs,omitempty"`
	}
	// tfApp represents an application in Terraform variable format.
	tfApp struct {
		Name             string                `json:"name"`
		PlatformType     string                `json:"platform_type"`
		PlatformProvider string                `json:"platform_provider"`
		AppType          string                `json:"app_type"`
		Path             string                `json:"path"`
		Config           map[string]any        `json:"config"`
		Environment      map[string][]tfEnvVar `json:"environment,omitempty"`
		Outputs          []string              `json:"outputs,omitempty"`
	}
	// tfEnvVar represents an environment variable for Terraform
	tfEnvVar struct {
		Key    string `json:"key"`
		Value  any    `json:"value,omitempty"`
		Source string `json:"source,omitempty"`
	}
	// tfGithubConfig is used to pull image containers from GH
	// container registry.
	tfGithubConfig struct {
		User  string `json:"user"`
		Repo  string `json:"repo"`
		Token string `json:"token"`
	}
	tfDigitalOceanConfig struct {
	}
)

// tfVarsFromDefinition transforms an appdef.Definition into Terraform variables.
// It converts the app.json structure into the format expected by the
// webkit-infra Terraform modules.
func tfVarsFromDefinition(env env.Environment, def *appdef.Definition) (tfVars, error) {
	if def == nil {
		return tfVars{}, fmt.Errorf("definition cannot be nil")
	}

	vars := tfVars{
		ProjectName: def.Project.Name,
		Environment: env.String(),
		GithubConfig: tfGithubConfig{
			User: def.Project.Repo.Owner,
			Repo: def.Project.Repo.Repo,
		},
		Apps:      make([]tfApp, 0, len(def.Apps)),
		Resources: make([]tfResource, 0, len(def.Resources)),
	}

	for _, res := range def.Resources {
		vars.Resources = append(vars.Resources, tfResource{
			Name:             res.Name,
			PlatformType:     res.Type.String(),
			PlatformProvider: res.Provider.String(),
			Config:           res.Config,
			Outputs:          res.Outputs,
		})
	}

	for _, app := range def.Apps {
		tfApp := tfApp{
			Name:             app.Name,
			PlatformType:     app.Infra.Type,
			PlatformProvider: app.Infra.Provider.String(),
			AppType:          app.Type.String(),
			Config:           app.Infra.Config,
			Path:             app.Path,
		}

		var tfEnvVars []tfEnvVar
		app.MergeEnvironments(def.Shared.Env).
			Walk(func(entry appdef.EnvWalkEntry) {
				tfEnv := tfEnvVar{
					Key:    entry.Key,
					Value:  entry.Value,
					Source: entry.Source.String(),
				}
				tfEnvVars = append(tfEnvVars, tfEnv)
			})

		vars.Apps = append(vars.Apps, tfApp)
	}

	return vars, nil
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
