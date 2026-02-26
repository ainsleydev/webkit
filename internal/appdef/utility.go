package appdef

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/ainsleydev/webkit/internal/appdef/types"
)

type (
	// UtilityCI defines the CI configuration for a utility.
	// When present, a CI job will be generated for this utility.
	// Omit CI entirely to create a workspace-only utility with no CI job.
	UtilityCI struct {
		Trigger string `json:"trigger" validate:"required,oneof=pull_request" description:"Event that triggers this CI job (currently only pull_request)"`
		RunsOn  string `json:"runs_on,omitempty" description:"GitHub Actions runner (defaults to ubuntu-latest)"`
	}

	// Utility represents a non-deployed workspace member within the webkit project.
	// Utilities are included in the pnpm workspace (if JS) and optionally run on CI,
	// but are never deployed. Examples include E2E tests, shared constants,
	// benchmark suites, and CLI tools.
	Utility struct {
		Name        string                                  `json:"name" validate:"required,lowercase,alphanumdash" description:"Unique identifier for the utility (lowercase, hyphenated)"`
		Title       string                                  `json:"title" validate:"required" description:"Human-readable utility name for display purposes"`
		Description string                                  `json:"description,omitempty" validate:"omitempty,max=200" description:"Brief description of the utility's purpose and functionality"`
		Path        string                                  `json:"path" validate:"required" description:"Relative file path to the utility's source code directory"`
		Language    string                                  `json:"language" validate:"required,oneof=go js" description:"Toolchain language for CI setup and workspace inclusion (go or js)"`
		CI          *UtilityCI                              `json:"ci,omitempty" description:"CI configuration. Omit to create a workspace-only utility with no CI job"`
		Tools       map[string]Tool                         `json:"tools,omitempty" inline:"true" description:"Build tools required for CI/CD workflows"`
		Commands    *types.OrderedMap[Command, CommandSpec] `json:"commands,omitzero" jsonschema:"oneof_type=boolean;object;string" inline:"true" description:"Custom commands for the utility"`
	}
)

// HasCI returns whether this utility has CI configuration and should
// generate a CI job.
func (u *Utility) HasCI() bool {
	return u.CI != nil
}

// OrderedCommands returns the utility's commands in their defined order
// with Name populated. This includes all commands in the order they
// appear in the OrderedMap.
func (u *Utility) OrderedCommands() []CommandSpec {
	if u.Commands == nil {
		return nil
	}

	var ordered []CommandSpec
	for _, cmd := range u.Commands.Keys() {
		spec, exists := u.Commands.Get(cmd)
		if !exists {
			continue
		}
		spec.Name = cmd.String()
		ordered = append(ordered, spec)
	}

	return ordered
}

// ShouldUseNPM returns whether this utility should be included in
// the pnpm workspace. JS utilities are included, Go utilities are not.
func (u *Utility) ShouldUseNPM() bool {
	return u.Language == "js"
}

// InstallCommands returns the shell commands needed to install all tools
// for this utility. Commands are generated based on the tool's Type field:
//   - "go": generates "go install <name>@<version>"
//   - "pnpm": generates "pnpm add -g <name>@<version>"
//   - "script": uses the Install field directly
//
// If a tool provides an Install field, it overrides the auto-generated command.
// Tools are sorted alphabetically by name to ensure deterministic output.
func (u *Utility) InstallCommands() []string {
	toolNames := make([]string, 0, len(u.Tools))
	for name := range u.Tools {
		toolNames = append(toolNames, name)
	}
	sort.Strings(toolNames)

	var commands []string
	for _, name := range toolNames {
		tool := u.Tools[name]

		if tool.Install != "" {
			commands = append(commands, tool.Install)
			continue
		}

		switch tool.Type {
		case "go":
			if tool.Name != "" && tool.Version != "" {
				commands = append(commands, fmt.Sprintf("go install %s@%s", tool.Name, tool.Version))
			}
		case "pnpm":
			if tool.Name != "" && tool.Version != "" {
				commands = append(commands, fmt.Sprintf("pnpm add -g %s@%s", tool.Name, tool.Version))
			}
		case "script":
			continue
		}
	}

	return commands
}

// applyDefaults sets default values for the utility.
func (u *Utility) applyDefaults() {
	if u.Commands == nil {
		u.Commands = types.NewOrderedMap[Command, CommandSpec]()
	}
	if u.Tools == nil {
		u.Tools = make(map[string]Tool)
	}
	if u.Path != "" {
		u.Path = filepath.Clean(u.Path)
	}
	if u.CI != nil && u.CI.RunsOn == "" {
		u.CI.RunsOn = "ubuntu-latest"
	}
}
