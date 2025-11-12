package appdef

import (
	"fmt"
	"path/filepath"
)

type (
	// App represents a single application within the webkit project.
	// Each app has its own build configuration, infrastructure requirements,
	// environment variables, and deployment settings. Apps can be of different
	// types (Payload CMS, SvelteKit, GoLang) and are deployed independently.
	App struct {
		Name             string                  `json:"name" validate:"required,lowercase,alphanumdash" description:"Unique identifier for the app (lowercase, hyphenated)"`
		Title            string                  `json:"title" validate:"required" description:"Human-readable app name for display purposes"`
		Type             AppType                 `json:"type" validate:"required,oneof=svelte-kit golang payload" description:"Application type (payload, svelte-kit, golang)"`
		Description      string                  `json:"description,omitempty" validate:"omitempty,max=200" description:"Brief description of the app's purpose and functionality"`
		Path             string                  `json:"path" validate:"required" description:"Relative file path to the app's source code directory"`
		Build            Build                   `json:"build" description:"Build configuration for Docker containerisation"`
		Infra            Infra                   `json:"infra" validate:"required" description:"Infrastructure and deployment configuration"`
		Env              Environment             `json:"env" description:"Environment variables specific to this app"`
		UsesNPM          *bool                   `json:"usesNPM" description:"Whether this app should be included in the pnpm workspace (auto-detected if not set)"`
		TerraformManaged *bool                   `json:"terraformManaged,omitempty" description:"Whether this app's infrastructure is managed by Terraform (defaults to true)"`
		Domains          []Domain                `json:"domains,omitzero" description:"Domain configurations for accessing this app"`
		Tools            map[string]string       `json:"tools,omitempty" description:"Build tools and their versions required for CI/CD (e.g., golangci-lint: latest)"`
		Commands         map[Command]CommandSpec `json:"commands,omitzero" jsonschema:"oneof_type=boolean;object;string" inline:"true" description:"Custom commands for linting, testing, formatting, and building"`
	}
	// Build defines Docker build configuration for containerised applications.
	// These settings control how the app is built and exposed in container environments.
	Build struct {
		Dockerfile string `json:"dockerfile" description:"Path to the Dockerfile relative to the app directory"`
		Port       int    `json:"port,omitempty" validate:"omitempty,min=1,max=65535" description:"Port number the app listens on inside the container"`
		Release    *bool  `json:"release,omitempty" description:"Whether to build and release this app in CI/CD (defaults to true)"`
	}
	// Infra defines infrastructure and deployment configuration for an app.
	// This includes the cloud provider, deployment type (VM, container, etc.),
	// and provider-specific configuration options.
	Infra struct {
		Provider ResourceProvider `json:"provider" validate:"required" description:"Cloud infrastructure provider (digitalocean, backblaze)"`
		Type     string           `json:"type" validate:"required" description:"Infrastructure type (vm, app, container, function)"`
		Config   map[string]any   `json:"config" description:"Provider-specific infrastructure configuration options"`
	}
	// Domain represents a domain name configuration for accessing an app.
	// Domains can be primary, aliases, or unmanaged depending on your DNS setup.
	Domain struct {
		Name     string     `json:"name" validate:"required" description:"Domain name without protocol (e.g., 'example.com' or 'api.example.com')"`
		Type     DomainType `json:"type" validate:"omitempty,oneof=primary alias unmanaged" description:"Domain type (primary, alias, unmanaged)"`
		Zone     string     `json:"zone,omitempty" description:"DNS zone for the domain if different from default"`
		Wildcard bool       `json:"wildcard,omitempty" description:"Whether this is a wildcard domain (e.g., '*.example.com')"`
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

// IsTerraformManaged returns whether this app should be managed by Terraform.
// It defaults to true when the field is nil or explicitly set to true.
func (a *App) IsTerraformManaged() bool {
	if a.TerraformManaged == nil {
		return true
	}
	return *a.TerraformManaged
}

// ShouldRelease returns whether this app should be built and released in CI/CD.
// It defaults to true when the field is nil or explicitly set to true.
func (a *App) ShouldRelease() bool {
	if a.Build.Release == nil {
		return true
	}
	return *a.Build.Release
}

// PrimaryDomain returns the primary domain for this app.
// It first looks for a domain with type "primary" in the Domains array.
// If no primary domain is found, it returns the first domain in the array.
// If the Domains array is empty, it returns an empty string.
func (a *App) PrimaryDomain() string {
	// First try to find a primary domain
	for _, domain := range a.Domains {
		if domain.Type == DomainTypePrimary {
			return domain.Name
		}
	}
	// Fallback to first domain if exists
	if len(a.Domains) > 0 {
		return a.Domains[0].Name
	}
	return ""
}

// ResolvedTools returns the app's tools with defaults merged in.
// Default tools for the app type are included first, then overridden
// by any explicitly configured tools. Tools can be disabled by setting
// their version to an empty string or "disabled".
func (a *App) ResolvedTools() map[string]string {
	tools := make(map[string]string)

	// Start with defaults for app type.
	if defaults, ok := defaultTools[a.Type]; ok {
		for k, v := range defaults {
			tools[k] = v
		}
	}

	// Override with explicit tools.
	for k, v := range a.Tools {
		if v == "" || v == "disabled" {
			delete(tools, k)
		} else {
			tools[k] = v
		}
	}

	return tools
}

// InstallCommands returns the shell commands needed to install
// all resolved tools for this app. For Go tools, it uses the
// goToolRegistry to map tool names to their full install paths.
func (a *App) InstallCommands() []string {
	var commands []string

	for tool, version := range a.ResolvedTools() {
		if installPath, ok := goToolRegistry[tool]; ok {
			commands = append(commands, fmt.Sprintf("go install %s@%s", installPath, version))
		} else {
			// Assume tool is already a full install path.
			commands = append(commands, fmt.Sprintf("go install %s@%s", tool, version))
		}
	}

	return commands
}

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

	if a.Build.Port == 0 {
		a.Build.Port = a.defaultPort()
	}

	if a.Path != "" {
		a.Path = filepath.Clean(a.Path)
	}

	return nil
}

// defaultPort returns the default port for the app based on its type.
// - Payload CMS: 3000
// - SvelteKit: 3001
// - GoLang: 8080
func (a *App) defaultPort() int {
	switch a.Type {
	case AppTypePayload:
		return 3000
	case AppTypeSvelteKit:
		return 3001
	case AppTypeGoLang:
		return 8080
	default:
		return 3000
	}
}
