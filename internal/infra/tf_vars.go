package infra

import (
	"fmt"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/pkg/env"
)

type (
	// TFVars represents the root structure of Terraform variables
	// that will be written to webkit.auto.tfvars.json
	TFVars struct {
		ProjectName string       `json:"project_name"`
		Environment string       `json:"environment"`
		Apps        []TFApp      `json:"apps"`
		Resources   []TFResource `json:"resources"`
	}
	// TFApp represents an application in Terraform variable format.
	TFApp struct {
		Name        string                `json:"name"`
		Type        string                `json:"type"`
		Path        string                `json:"path"`
		Infra       TFInfra               `json:"infra"`
		Environment map[string][]TFEnvVar `json:"environment,omitempty"`
	}
	// TFInfra represents infrastructure configuration for an app.
	TFInfra struct {
		Provider string         `json:"provider"`
		Type     string         `json:"type"`
		Config   map[string]any `json:"config"`
	}
	// TFResource represents a resource in Terraform variable format.
	TFResource struct {
		Name     string         `json:"name"`
		Type     string         `json:"type"`
		Provider string         `json:"provider"`
		Config   map[string]any `json:"config"`
		Outputs  []string       `json:"outputs,omitempty"`
	}
	// TFEnvVar represents an environment variable for Terraform
	TFEnvVar struct {
		Key    string `json:"key"`
		Value  any    `json:"value,omitempty"`
		Source string `json:"source,omitempty"`
	}
)

// TFVarsFromDefinition transforms an appdef.Definition into Terraform variables.
// It converts the app.json structure into the format expected by the
// webkit-infra Terraform modules.
func TFVarsFromDefinition(def *appdef.Definition, env env.Environment) (TFVars, error) {
	if def == nil {
		return TFVars{}, fmt.Errorf("definition cannot be nil")
	}

	vars := TFVars{
		ProjectName: def.Project.Name,
		Environment: env.String(),
		Apps:        make([]TFApp, 0, len(def.Apps)),
		Resources:   make([]TFResource, 0, len(def.Resources)),
	}

	// Transform resources
	for _, res := range def.Resources {
		vars.Resources = append(vars.Resources, TFResource{
			Name:     res.Name,
			Type:     res.Type.String(),
			Provider: res.Provider.String(),
			Config:   res.Config,
			Outputs:  res.Outputs,
		})
	}

	// Transform apps
	for _, app := range def.Apps {
		tfApp := TFApp{
			Name: app.Name,
			Type: app.Type.String(),
			Path: app.Path,
		}

		// Transform infra config
		if app.Infra != nil {
			tfApp.Infra = TFInfra{
				Provider: app.Infra.Provider,
				Type:     app.Infra.Type,
				Config:   app.Infra.Config,
			}
		}

		// Transform environment variables
		var tfEnvVars []TFEnvVar
		app.MergeEnvironments(def.Shared.Env).
			Walk(func(entry appdef.EnvWalkEntry) {
				tfEnv := TFEnvVar{
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
