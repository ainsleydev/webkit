package appdef

import (
	"errors"
	"fmt"
	"io"

	"github.com/goccy/go-json"
	"github.com/spf13/afero"
)

const (
	// JsonFileName defines the file name of the app manifest,
	// that should appear in the root of each webkit dir.
	JsonFileName = "app.json"
)

type (
	// Definition represents the complete webkit application configuration.
	// It defines the structure of the app.json file used to configure
	// all aspects of a webkit project including apps, resources, and infrastructure.
	Definition struct {
		Schema        string        `json:"$schema,omitempty" jsonschema:"-" description:"JSON Schema reference for IDE validation and autocomplete"`
		WebkitVersion string        `json:"webkit_version" required:"true" validate:"required" description:"The version of webkit used to generate this configuration"`
		Project       Project       `json:"project" required:"true" validate:"required" description:"Project metadata including name, title, and repository information"`
		Notifications Notifications `json:"notifications" description:"Alert and notification settings for the project"`
		Shared        Shared        `json:"shared" description:"Shared configuration that applies to all apps"`
		Resources     []Resource    `json:"resources" description:"Infrastructure resources such as databases and storage buckets"`
		Apps          []App         `json:"apps" required:"true" validate:"required,min=1,dive" minItems:"1" description:"Application definitions for all apps in the project"`
	}
	// Shared contains configuration that is shared across all applications
	// in the project, such as common environment variables.
	Shared struct {
		Env Environment `json:"env" description:"Environment variables shared across all apps"`
	}
	// SkippedItems contains information about apps and resources
	// that were filtered out due to not being Terraform managed.
	SkippedItems struct {
		Apps      []string
		Resources []string
	}
)

/************************************
	General
************************************/

func Read(root afero.Fs) (*Definition, error) {
	file, err := root.Open(JsonFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	def := &Definition{}
	if err = json.Unmarshal(data, def); err != nil {
		return nil, errors.New("unmarshalling app definition: " + err.Error())
	}

	if err = def.ApplyDefaults(); err != nil {
		return nil, err
	}

	return def, nil
}

// GithubLabels returns the labels that will appear on the
// GitHub repository by looking at the application types.
func (d *Definition) GithubLabels() []string {
	labels := []string{"webkit"}

	for _, v := range d.Apps {
		labels = append(labels, v.Type.String())
	}

	return labels
}

// ApplyDefaults ensures all required defaults are set on the Definition.
// This should be called after unmarshaling and before validation.
func (d *Definition) ApplyDefaults() error {
	for i := range d.Apps {
		if err := d.Apps[i].applyDefaults(); err != nil {
			return fmt.Errorf("applying defaults to app %q: %w", d.Apps[i].Name, err)
		}
	}

	for i := range d.Resources {
		d.Resources[i].applyDefaults()
	}

	return nil
}

// FilterTerraformManaged creates a filtered copy of the Definition containing
// only apps and resources that are managed by Terraform.
//
// Returns the filtered definition and information about what was skipped.
func (d *Definition) FilterTerraformManaged() (*Definition, SkippedItems) {
	filtered := &Definition{
		WebkitVersion: d.WebkitVersion,
		Project:       d.Project,
		Shared:        d.Shared,
		Apps:          make([]App, 0, len(d.Apps)),
		Resources:     make([]Resource, 0, len(d.Resources)),
	}

	skipped := SkippedItems{
		Apps:      make([]string, 0),
		Resources: make([]string, 0),
	}

	// Filter resources.
	for _, res := range d.Resources {
		if res.IsTerraformManaged() {
			filtered.Resources = append(filtered.Resources, res)
		} else {
			skipped.Resources = append(skipped.Resources, res.Name)
		}
	}

	// Filter apps.
	for _, app := range d.Apps {
		if app.IsTerraformManaged() {
			filtered.Apps = append(filtered.Apps, app)
		} else {
			skipped.Apps = append(skipped.Apps, app.Name)
		}
	}

	return filtered, skipped
}

/************************************
	Apps
************************************/

// ContainsGo returns true if any of the apps are marked as Go.
func (d *Definition) ContainsGo() bool {
	for _, app := range d.Apps {
		if app.Language() == "go" {
			return true
		}
	}
	return false
}

// ContainsJS returns true if any of the apps are marked as JS.
func (d *Definition) ContainsJS() bool {
	for _, app := range d.Apps {
		if app.Language() == "js" {
			return true
		}
	}
	return false
}

// HasAppType checks if the definition contains an app of the
// specified type.
func (d *Definition) HasAppType(appType AppType) bool {
	if d == nil {
		return false
	}

	for _, app := range d.Apps {
		if app.Type == appType {
			return true
		}
	}

	return false
}

// GetAppsByType returns all apps of the specified type from
// the definition.
func (d *Definition) GetAppsByType(appType AppType) []App {
	if d == nil {
		return nil
	}

	var apps []App
	for _, app := range d.Apps {
		if app.Type == appType {
			apps = append(apps, app)
		}
	}

	return apps
}

/************************************
	Env
************************************/

// MergeAllEnvironments merges shared env variables with all apps' environments.
// App-specific values take precedence over shared ones. If multiple apps define the same variable,
// the last app in the list wins.
func (d *Definition) MergeAllEnvironments() Environment {
	merged := mergeEnvironments(d.Shared.Env)

	for _, app := range d.Apps {
		merged = mergeEnvironments(merged, app.Env)
	}

	return merged
}

// mergeEnvironments merges multiple Environment structs left-to-right.
// Later environments override earlier ones.
// Default values are applied to each environment before merging environment-specific values.
func mergeEnvironments(envs ...Environment) Environment {
	merged := Environment{
		Dev:        make(EnvVar),
		Staging:    make(EnvVar),
		Production: make(EnvVar),
	}

	for _, env := range envs {
		// Apply defaults to each environment first, then apply environment-specific overrides.
		devWithDefaults := mergeVars(env.Default, env.Dev)
		stagingWithDefaults := mergeVars(env.Default, env.Staging)
		productionWithDefaults := mergeVars(env.Default, env.Production)

		// Merge into accumulated result.
		merged.Dev = mergeVars(merged.Dev, devWithDefaults)
		merged.Staging = mergeVars(merged.Staging, stagingWithDefaults)
		merged.Production = mergeVars(merged.Production, productionWithDefaults)
	}

	return merged
}
