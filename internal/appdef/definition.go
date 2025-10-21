// Package appdef provides types and functions for parsing and working with
// the WebKit application definition file (app.json).
package appdef

import (
	"errors"
	"fmt"
	"io"

	"github.com/goccy/go-json"
	"github.com/spf13/afero"
)

const (
	// JsonFileName defines the file name of the app manifest
	// that should appear in the root of each WebKit directory.
	JsonFileName = "app.json"
)

type (
	// Definition represents the complete WebKit application manifest,
	// including project metadata, shared configuration, resources, and apps.
	Definition struct {
		WebkitVersion string     `json:"webkit_version"`
		Project       Project    `json:"project"`
		Shared        Shared     `json:"shared"`
		Resources     []Resource `json:"resources"`
		Apps          []App      `json:"apps"`
	}

	// Shared contains configuration that is shared across all apps,
	// such as environment variables.
	Shared struct {
		Env Environment `json:"env"`
	}
)

// Read reads and parses the WebKit application definition from the filesystem.
// It unmarshals the JSON file, applies default values, and returns the complete
// definition or an error if parsing fails.
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
