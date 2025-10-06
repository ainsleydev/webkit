package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal"
	"github.com/ainsleydev/webkit/internal/scaffold"
)

type (
	packageJSON struct {
		Name            string            `json:"name"`
		Version         string            `json:"version"`
		Description     string            `json:"description"`
		Private         bool              `json:"private"`
		License         string            `json:"license"`
		Type            string            `json:"type"`
		Scripts         map[string]string `json:"scripts"`
		DevDependencies map[string]string `json:"devDependencies"`
		PackageManager  string            `json:"packageManager"`
		Engines         map[string]string `json:"engines"`
		Pnpm            packagePnpm       `json:"pnpm"`
		Author          packageAuthor     `json:"author"`
		Maintainers     []packageAuthor   `json:"maintainers"`
	}
	packageAuthor struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		URL   string `json:"url"`
	}
	packagePnpm struct {
		OnlyBuiltDependencies []string `json:"onlyBuiltDependencies"`
	}
)

var createPackageJsonCmd = &cli.Command{
	Name:   "package-json",
	Action: cmdtools.WrapCommand(createPackageJson),
}

func createPackageJson(_ context.Context, input cmdtools.CommandInput) error {
	gen := scaffold.New(input.FS)
	app := input.AppDef()

	p := packageJSON{
		Name:        app.Project.Name,
		Description: app.Project.Description,
		Version:     "1.0.0",
		Private:     true,
		License:     "BSD-3-Clause",
		Type:        "module",
		Scripts: map[string]string{
			"preinstall": "npx only-allow pnpm",
			"test":       "turbo test",
			"lint":       "eslint .",
			"lint:fix":   "eslint . --fix",
			"format":     "prettier --write .",
		},
		DevDependencies: map[string]string{
			"@ainsleydev/eslint-config":   "^0.0.6",
			"@ainsleydev/prettier-config": "^0.0.2",
			"@eslint/compat":              "^1.4.0",
			"@payloadcms/eslint-config":   "^3.28.0",
			"@payloadcms/eslint-plugin":   "^3.28.0",
			"eslint":                      "^9.37.0",
			"eslint-plugin-perfectionist": "^4.15.1",
			"globals":                     "^16.0.0",
			"prettier":                    "^3.6.0",
			"prettier-plugin-svelte":      "^3.4.0",
			"turbo":                       "^2.5.8",
			"typescript":                  "5.8.2",
		},
		PackageManager: "pnpm@10.15.0",
		Engines:        nil,
		Pnpm:           packagePnpm{},
		Author: packageAuthor{
			Name:  "ainsley.dev LTD",
			Email: "hello@ainsley.dev",
			URL:   "https://ainsley.dev",
		},
		Maintainers: []packageAuthor{
			{
				Name:  "Ainsley Clark",
				Email: "hello@ainsley.dev",
				URL:   "https://ainsley.dev",
			},
		},
	}

	return gen.JSON("package.json", p)
}
