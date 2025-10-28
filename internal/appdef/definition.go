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
	Definition struct {
		WebkitVersion string     `json:"webkit_version"`
		Project       Project    `json:"project"`
		Shared        Shared     `json:"shared"`
		Resources     []Resource `json:"resources"`
		Apps          []App      `json:"apps"`
	}
	Shared struct {
		Env Environment `json:"env"`
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

	// TODO: Apply defaults and return validation errors if the user has fucked it.
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
func mergeEnvironments(envs ...Environment) Environment {
	merged := Environment{
		Dev:        make(EnvVar),
		Staging:    make(EnvVar),
		Production: make(EnvVar),
	}

	for _, env := range envs {
		merged.Dev = mergeVars(merged.Dev, env.Dev)
		merged.Staging = mergeVars(merged.Staging, env.Staging)
		merged.Production = mergeVars(merged.Production, env.Production)
	}

	return merged
}
