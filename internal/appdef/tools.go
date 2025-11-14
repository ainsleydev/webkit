package appdef

// Tool represents a build tool with its installation method.
// Supports three types: "go", "pnpm", and "script".
//
// For "go" and "pnpm" types, the Name and Version fields are used to generate
// the installation command automatically.
//
// For "script" type, the Install field must contain the complete installation command.
//
// The Install field can be used with any type to override the auto-generated command.
type Tool struct {
	Type    string `json:"type" inline:"true" description:"Installation method: 'go', 'pnpm', or 'script'"`
	Name    string `json:"name,omitempty" inline:"true" description:"Package path (for 'go') or package name (for 'pnpm')"`
	Version string `json:"version,omitempty" inline:"true" description:"Version to install (e.g., 'latest', 'v0.2.543')"`
	Install string `json:"install,omitempty" inline:"true" description:"Custom installation command (overrides auto-generated command)"`
}

// defaultTools defines the default tools and their installation commands for each app type.
// These tools are automatically installed in CI/CD workflows unless
// explicitly overridden or disabled in the app's Tools configuration.
var defaultTools = map[AppType]map[string]Tool{
	AppTypeGoLang: {
		"golangci-lint": {
			Type:    "go",
			Name:    "github.com/golangci/golangci-lint/cmd/golangci-lint",
			Version: "latest",
		},
		"templ": {
			Type:    "go",
			Name:    "github.com/a-h/templ/cmd/templ",
			Version: "latest",
		},
		"sqlc": {
			Type:    "go",
			Name:    "github.com/sqlc-dev/sqlc/cmd/sqlc",
			Version: "latest",
		},
	},
	// JavaScript apps currently have no default tools.
	// Tools like eslint, prettier are typically installed via pnpm.
	AppTypeSvelteKit: {},
	AppTypePayload:   {},
}
