package appdef

import (
	"fmt"
	"sort"

	"github.com/ainsleydev/webkit/internal/appdef/types"
)

// Toolset contains tools and commands that are shared between App and Utility.
// Both types embed this struct to avoid duplicating the OrderedCommands and
// InstallCommands methods.
type Toolset struct {
	Tools    map[string]Tool                         `json:"tools,omitempty" inline:"true" description:"Build tools required for CI/CD workflows"`
	Commands *types.OrderedMap[Command, CommandSpec] `json:"commands,omitzero" jsonschema:"oneof_type=boolean;object;string" inline:"true" description:"Custom commands for linting, testing, formatting, and building"`
}

// OrderedCommands returns the commands in their defined order with Name populated.
// This includes all commands in the order they appear in the OrderedMap.
func (t *Toolset) OrderedCommands() []CommandSpec {
	if t.Commands == nil {
		return nil
	}

	var ordered []CommandSpec
	for _, cmd := range t.Commands.Keys() {
		spec, exists := t.Commands.Get(cmd)
		if !exists {
			continue
		}
		spec.Name = cmd.String()
		ordered = append(ordered, spec)
	}

	return ordered
}

// InstallCommands returns the shell commands needed to install all tools.
// Commands are generated based on the tool's Type field:
//   - "go": generates "go install <name>@<version>"
//   - "pnpm": generates "pnpm add -g <name>@<version>"
//   - "script": uses the Install field directly
//
// If a tool provides an Install field, it overrides the auto-generated command.
// Tools are sorted alphabetically by name to ensure deterministic output.
func (t *Toolset) InstallCommands() []string {
	toolNames := make([]string, 0, len(t.Tools))
	for name := range t.Tools {
		toolNames = append(toolNames, name)
	}
	sort.Strings(toolNames)

	var commands []string
	for _, name := range toolNames {
		tool := t.Tools[name]

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

// initDefaults initialises nil Commands and Tools to their zero values.
func (t *Toolset) initDefaults() {
	if t.Commands == nil {
		t.Commands = types.NewOrderedMap[Command, CommandSpec]()
	}
	if t.Tools == nil {
		t.Tools = make(map[string]Tool)
	}
}
