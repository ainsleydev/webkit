package appdef

// defaultTools defines the default tools and versions for each app type.
// These tools are automatically installed in CI/CD workflows unless
// explicitly overridden or disabled in the app's Tools configuration.
var defaultTools = map[AppType]map[string]string{
	AppTypeGoLang: {
		"golangci-lint": "latest",
		"templ":         "latest",
		"sqlc":          "latest",
	},
	// JavaScript apps currently have no default tools.
	// Tools like eslint, prettier are typically installed via pnpm.
	AppTypeSvelteKit: {},
	AppTypePayload:   {},
}

// goToolRegistry maps well-known Go tool names to their full install paths.
// This allows users to specify short names like "golangci-lint" instead of
// the full path "github.com/golangci/golangci-lint/cmd/golangci-lint".
var goToolRegistry = map[string]string{
	"golangci-lint": "github.com/golangci/golangci-lint/cmd/golangci-lint",
	"templ":         "github.com/a-h/templ/cmd/templ",
	"sqlc":          "github.com/sqlc-dev/sqlc/cmd/sqlc",
	"buf":           "github.com/bufbuild/buf/cmd/buf",
	"wire":          "github.com/google/wire/cmd/wire",
	"mockgen":       "go.uber.org/mock/mockgen",
}
