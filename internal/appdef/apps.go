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
		Env         Env                     `json:"env"`
		Commands    map[Command]CommandSpec `json:"commands,omitempty" jsonschema:"oneof_type=boolean;object;string"`
		DependsOn   []string                `json:"depends_on,omitempty"`
	}
	Build struct {
		Dockerfile string `json:"dockerfile"`
	}
	Infra struct {
		Provider string `json:"provider"`
		Type     string `json:"type"`
		Config   struct {
			Size          string   `json:"size,omitempty"`
			Region        string   `json:"region"`
			Domain        string   `json:"domain"`
			SshKeys       []string `json:"ssh_keys,omitempty"`
			InstanceCount int      `json:"instance_count,omitempty"`
			EnvFromShared bool     `json:"env_from_shared,omitempty"`
		} `json:"config"`
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
