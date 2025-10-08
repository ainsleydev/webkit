package appdef

import (
	"fmt"
	"path/filepath"
)

type (
	App struct {
		Name        string                  `json:"name"`
		Title       string                  `json:"title"`
		Type        AppType                 `json:"type"`
		Description string                  `json:"description,omitempty"`
		Path        string                  `json:"path"`
		Build       Build                   `json:"build"`
		Infra       Infra                   `json:"infra"`
		Env         Environment             `json:"env"`
		Domains     []Domain                `json:"domains,omitempty"`
		Commands    map[Command]CommandSpec `json:"commands,omitempty" jsonschema:"oneof_type=boolean;object;string"`
	}
	Build struct {
		Dockerfile string `json:"dockerfile"`
	}
	Infra struct {
		Provider string         `json:"provider"`
		Type     string         `json:"type"`
		Config   map[string]any `json:"config"`
	}
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

func (a *App) applyDefaults() error {
	if a.Commands == nil {
		a.Commands = make(map[Command]CommandSpec)
	}

	// Get default commands for this app type
	defaults, hasDefaults := defaultCommands[a.Type]
	if !hasDefaults {
		return fmt.Errorf("no default commands defined for app type %q", a.Type)
	}

	for _, cmd := range commands {
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
