package appdef

import (
	"encoding/json"
	"fmt"
)

// Tool represents a build tool with its version and installation command.
type Tool struct {
	Version string // Version to install (e.g., "latest", "v0.2.543")
	Install string // Complete installation command
}

// ToolSpec defines a tool configuration that can be specified as either
// a string (version override) or an object (full configuration).
// Examples:
//   - "latest"           - uses default install command with latest version
//   - "v0.2.543"         - uses default install command with specific version
//   - {"version": "...", "install": "..."} - full control over installation
type ToolSpec struct {
	Version string `json:"version,omitempty" description:"Version to install (e.g., 'latest', 'v0.2.543')"`
	Install string `json:"install,omitempty" description:"Custom installation command"`
}

// UnmarshalJSON implements json.Unmarshaler to handle string or object formats.
func (t *ToolSpec) UnmarshalJSON(data []byte) error {
	// Try string (version override)
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		// Special case: empty string or "disabled" means disabled
		if str == "" || str == "disabled" {
			t.Version = "disabled"
			return nil
		}
		t.Version = str
		return nil
	}

	// Try object (full control)
	type Alias ToolSpec
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("invalid tool format: expected string or object")
	}

	*t = ToolSpec(aux)
	return nil
}

// defaultTools defines the default tools and their installation commands for each app type.
// These tools are automatically installed in CI/CD workflows unless
// explicitly overridden or disabled in the app's Tools configuration.
var defaultTools = map[AppType]map[string]Tool{
	AppTypeGoLang: {
		"golangci-lint": {
			Version: "latest",
			Install: "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
		},
		"templ": {
			Version: "latest",
			Install: "go install github.com/a-h/templ/cmd/templ@latest",
		},
		"sqlc": {
			Version: "latest",
			Install: "go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest",
		},
	},
	// JavaScript apps currently have no default tools.
	// Tools like eslint, prettier are typically installed via pnpm.
	AppTypeSvelteKit: {},
	AppTypePayload:   {},
}

// goToolRegistry maps well-known Go tool names to their full install paths.
// This allows users to override versions while keeping the default install method.
// When a user specifies a version override for a known tool, we construct the
// install command using this registry.
var goToolRegistry = map[string]string{
	"golangci-lint": "github.com/golangci/golangci-lint/cmd/golangci-lint",
	"templ":         "github.com/a-h/templ/cmd/templ",
	"sqlc":          "github.com/sqlc-dev/sqlc/cmd/sqlc",
	"buf":           "github.com/bufbuild/buf/cmd/buf",
	"wire":          "github.com/google/wire/cmd/wire",
	"mockgen":       "go.uber.org/mock/mockgen",
}
