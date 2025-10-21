package appdef

import (
	"fmt"
	"path/filepath"
)

type (
	// App represents a single application within the WebKit project,
	// including its configuration, build settings, infrastructure, and domains.
	App struct {
		Name        string                  `json:"name"`
		Title       string                  `json:"title"`
		Type        AppType                 `json:"type"`
		Description string                  `json:"description,omitempty"`
		Path        string                  `json:"path"`
		Build       Build                   `json:"build"`
		Infra       Infra                   `json:"infra"`
		Env         Environment             `json:"env"`
		UsesNPM     *bool                   `json:"usesNPM"`
		Domains     []Domain                `json:"domains,omitzero"`
		Commands    map[Command]CommandSpec `json:"commands,omitzero" jsonschema:"oneof_type=boolean;object;string"`
	}

	// Build contains build-related configuration for an app.
	Build struct {
		Dockerfile string `json:"dockerfile"`
	}

	// Infra defines the infrastructure configuration for an app,
	// including the cloud provider and deployment settings.
	Infra struct {
		Provider ResourceProvider `json:"provider"`
		Type     string            `json:"type"`
		Config   map[string]any    `json:"config"`
	}

	// Domain represents a domain name configuration for an app.
	Domain struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Zone     string `json:"zone,omitempty"`
		Wildcard bool   `json:"wildcard,omitempty"`
	}
)

// AppType defines the type of application being run.
type AppType string

// AppType constants.
const (
	AppTypeSvelteKit AppType = "svelte-kit"
	AppTypeGoLang    AppType = "golang"
	AppTypePayload   AppType = "payload"
)

// String implements fmt.Stringer on the AppType.
func (a AppType) String() string {
	return string(a)
}

var appTypeToLanguages = map[AppType]string{
	AppTypeGoLang:    "go",
	AppTypeSvelteKit: "js",
	AppTypePayload:   "js",
}

// Language determines what language ecosystem a given app is.
// Either "go" or "js".
func (a *App) Language() string {
	return appTypeToLanguages[a.Type]
}

// DomainType defines the type of domain that should be provisioned.
type DomainType string

// DomainType constants.
const (
	DomainTypePrimary   DomainType = "primary"
	DomainTypeAlias     DomainType = "alias"
	DomainTypeUnmanaged DomainType = "unmanaged"
)

// String implements fmt.Stringer on the DomainType.
func (d DomainType) String() string {
	return string(d)
}

// OrderedCommands returns the app's commands in canonical order
// with Name populated.
func (a *App) OrderedCommands() []CommandSpec {
	var ordered []CommandSpec

	for _, cmd := range Commands {
		spec, exists := a.Commands[cmd]
		if !exists {
			// Should not happen because applyDefaults populates them
			continue
		}
		spec.Name = cmd.String() // Populate name for templates.
		ordered = append(ordered, spec)
	}

	return ordered
}

// MergeEnvironments merges the shared env with the apps,
// with the app specific variables taking precedence.
func (a *App) MergeEnvironments(shared Environment) Environment {
	return mergeEnvironments(shared, a.Env)
}

// ShouldUseNPM returns whether this app should be included in
// pnpm workspace. It checks the UsesNPM field first, and if
// not set, defaults based on Language().
func (a *App) ShouldUseNPM() bool {
	if a.UsesNPM != nil {
		return *a.UsesNPM
	}
	return a.Language() == "js"
}

// applyDefaults applies default values to the App, including default commands
// for the app type and build configuration.
func (a *App) applyDefaults() error {
	if a.Commands == nil {
		a.Commands = make(map[Command]CommandSpec)
	}

	// Get default Commands for this app type
	defaults, hasDefaults := defaultCommands[a.Type]
	if !hasDefaults {
		return fmt.Errorf("no default Commands defined for app type %q", a.Type)
	}

	for _, cmd := range Commands {
		// Skip if user has explicitly configured this command.
		if _, exists := a.Commands[cmd]; exists {
			continue
		}

		// Apply default command if available.
		if defaultCmd, ok := defaults[cmd]; ok {
			a.Commands[cmd] = CommandSpec{
				Cmd: defaultCmd,
			}
		}
	}

	if a.Build.Dockerfile == "" {
		a.Build.Dockerfile = "Dockerfile"
	}

	if a.Path != "" {
		a.Path = filepath.Clean(a.Path)
	}

	return nil
}
