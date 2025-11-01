package files

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
)

// PackageJSON scaffolds a root JSON file to act as a
// starting point for repos.
func PackageJSON(_ context.Context, input cmdtools.CommandInput) error {
	app := input.AppDef()

	p := packageJSON{
		Name:        app.Project.Name,
		Description: app.Project.Description,
		Version:     "1.0.0",
		Private:     "false",
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

	return input.Generator().JSON("package.json", p,
		scaffold.WithTracking(manifest.SourceProject()),
		scaffold.WithScaffoldMode(),
	)
}

type (
	packageJSON struct {
		Name            string            `json:"name"`
		Version         string            `json:"version"`
		Description     string            `json:"description,omitempty"`
		Private         string            `json:"private"`
		License         string            `json:"license"`
		Type            string            `json:"type"`
		Scripts         map[string]string `json:"scripts"`
		DevDependencies map[string]string `json:"devDependencies"`
		PackageManager  string            `json:"packageManager"`
		Engines         map[string]string `json:"engines,omitempty"`
		Pnpm            packagePnpm       `json:"pnpm,omitzero"`
		Author          packageAuthor     `json:"author"`
		Maintainers     []packageAuthor   `json:"maintainers"`
	}
	packageAuthor struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		URL   string `json:"url"`
	}
	packagePnpm struct {
		OnlyBuiltDependencies []string `json:"onlyBuiltDependencies,omitempty"`
	}
)

// PackageJSONApp manipulates each app's package.json file for apps that use NPM.
// Currently adds Docker-related scripts while preserving existing scripts and other package.json fields.
func PackageJSONApp(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	for _, app := range appDef.Apps {
		// Skip apps that don't use NPM.
		if !app.ShouldUseNPM() {
			continue
		}

		// Construct the path to the app's package.json.
		pkgPath := filepath.Join(app.Path, "package.json")

		// Check if package.json exists.
		exists, err := afero.Exists(input.FS, pkgPath)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("checking if %s exists", pkgPath))
		}

		// Skip if package.json doesn't exist.
		if !exists {
			input.Printer().Println(fmt.Sprintf("• skipping %s - package.json not found", app.Name))
			continue
		}

		// Read existing package.json.
		data, err := afero.ReadFile(input.FS, pkgPath)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("reading %s", pkgPath))
		}

		// Parse the existing package.json.
		var pkg map[string]any
		if err := json.Unmarshal(data, &pkg); err != nil {
			return errors.Wrap(err, fmt.Sprintf("parsing %s", pkgPath))
		}

		// Get or create the scripts section.
		scripts, ok := pkg["scripts"].(map[string]any)
		if !ok {
			scripts = make(map[string]any)
			pkg["scripts"] = scripts
		}

		// Add Docker scripts with app name substitution.
		imageName := fmt.Sprintf("%s-web", app.Name)
		scripts["docker"] = "pnpm docker:build && pnpm docker:run"
		scripts["docker:build"] = fmt.Sprintf("docker build . -t %s --progress plain --no-cache", imageName)
		scripts["docker:run"] = fmt.Sprintf("docker run -it --init --env-file .env -p %d:%d --rm -ti %s",
			app.Build.Port, app.Build.Port, imageName)
		scripts["docker:remove"] = fmt.Sprintf("docker image rm %s", imageName)

		// Write back the modified package.json.
		if err := input.Generator().JSON(pkgPath, pkg, scaffold.WithTracking(manifest.SourceApp(app.Name))); err != nil {
			return errors.Wrap(err, fmt.Sprintf("writing %s", pkgPath))
		}

		input.Printer().Println(fmt.Sprintf("• added Docker scripts to %s", pkgPath))
	}

	return nil
}
