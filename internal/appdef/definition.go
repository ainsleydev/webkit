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
	Project struct {
		Name        string `json:"name" jsonschema:"required,pattern=^[a-z0-9-]+$,title=Project Name,description=Machine-readable project name (kebab-case),example=my-website"`
		Title       string `json:"title" jsonschema:"required,title=Project Title,description=Human-readable project title,example=My Website"`
		Description string `json:"description" jsonschema:"required,title=Description,description=Brief description of the project"`
		Repo        string `json:"repo" jsonschema:"required,format=uri,title=Repository,description=Git repository URL,example=git@github.com:ainsley/my-website.git"`
	}
	Shared struct {
		Env Environment `json:"env"`
	}
)

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
	if err := json.Unmarshal(data, def); err != nil {
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
		if err := d.Resources[i].applyDefaults(); err != nil {
			return fmt.Errorf("applying defaults to resource %q: %w", d.Resources[i].Name, err)
		}
	}

	return nil
}

/**** Env *****/

// MergeAllEnvironments merges shared environment variables with all apps' environments.
// App-specific values take precedence over shared ones. If multiple apps define the same variable,
// the last app in the list wins.
func (d *Definition) MergeAllEnvironments() Environment {
	merged := d.mergeEnvironments(d.Shared.Env)

	for _, app := range d.Apps {
		merged = d.mergeEnvironments(merged, app.Env)
	}

	return merged
}

// MergeAppEnvironment merges shared environment variables with a single app's environment.
// The app's variables take precedence over the shared ones.
func (d *Definition) MergeAppEnvironment(appName string) (Environment, bool) {
	var app *App
	for i := range d.Apps {
		if d.Apps[i].Name == appName {
			app = &d.Apps[i]
			break
		}
	}

	if app == nil {
		return Environment{}, false
	}

	return d.mergeEnvironments(d.Shared.Env, app.Env), true
}

// mergeEnvironments merges multiple Environment structs left-to-right.
// Later environments override earlier ones.
func (d *Definition) mergeEnvironments(envs ...Environment) Environment {
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
