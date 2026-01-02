package files

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/pkgjson"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/state/manifest"
)

// PackageJSON scaffolds a root JSON file to act as a
// starting point for repos.
func PackageJSON(_ context.Context, input cmdtools.CommandInput) error {
	app := input.AppDef()

	p := &pkgjson.PackageJSON{
		Name:        app.Project.Name,
		Description: app.Project.Description,
		Version:     "1.0.0",
		Private:     "false",
		License:     "BSD-3-Clause",
		Type:        "module",
		Scripts: map[string]any{
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
		Pnpm:           &pkgjson.PnpmConfig{},
		Author: &pkgjson.Author{
			Name:  "ainsley.dev LTD",
			Email: "hello@ainsley.dev",
			URL:   "https://ainsley.dev",
		},
		Maintainers: []pkgjson.Author{
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

// PackageJSONApp manipulates each app's package.json file for
// apps that use NPM. Currently, adds Docker-related scripts
// while preserving existing scripts.
func PackageJSONApp(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	for _, app := range appDef.Apps {
		if !app.ShouldUseNPM() {
			continue
		}

		pkgPath := filepath.Join(app.Path, "package.json")

		exists, err := afero.Exists(input.FS, pkgPath)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("checking if %s exists", pkgPath))
		} else if !exists {
			input.Printer().Println(fmt.Sprintf("• skipping %s - package.json not found", app.Name))
			continue
		}

		// Use shared pkgjson package for reading existing files.
		pkg, err := pkgjson.Read(input.FS, pkgPath)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("reading %s", pkgPath))
		}

		// Get or create the scripts section.
		if pkg.Scripts == nil {
			pkg.Scripts = make(map[string]any)
		}

		// Add Docker scripts with app name substitution.
		imageName := fmt.Sprintf("%s-%s", appDef.Project.Name, app.Name)
		pkg.Scripts["docker"] = "pnpm docker:build && pnpm docker:run"
		pkg.Scripts["docker:build"] = fmt.Sprintf("docker build . -t %s --progress plain --no-cache", imageName)
		pkg.Scripts["docker:run"] = fmt.Sprintf("docker run -it --init --env-file .env -p %d:%d --rm -ti %s",
			app.Build.Port, app.Build.Port, imageName)
		pkg.Scripts["docker:remove"] = fmt.Sprintf("docker image rm %s", imageName)

		// Add app-type-specific scripts.
		typeScripts := getAppTypeScripts(app.Type)
		for scriptName, scriptCommand := range typeScripts {
			pkg.Scripts[scriptName] = scriptCommand
		}

		// Write back the modified package.json using shared pkgjson package.
		if err = pkgjson.Write(input.FS, pkgPath, pkg); err != nil {
			return errors.Wrap(err, fmt.Sprintf("writing %s", pkgPath))
		}
	}

	return nil
}

// appTypeScripts maps app types to their specific package.json scripts.
// These scripts are injected into app package.json files based on the app type.
var appTypeScripts = map[appdef.AppType]map[string]string{
	appdef.AppTypePayload: {
		"migrate:check":  "node scripts/check-deps.js",
		"migrate:create": "pnpm migrate:check && NODE_ENV=production payload migrate:create",
		"migrate:status": "NODE_ENV=production payload migrate:status",
	},
	appdef.AppTypeSvelteKit: {},
	appdef.AppTypeGoLang:    {},
}

// getAppTypeScripts returns the scripts for a given app type.
// Returns an empty map if no scripts are defined for the app type.
func getAppTypeScripts(appType appdef.AppType) map[string]string {
	if scripts, ok := appTypeScripts[appType]; ok {
		return scripts
	}
	return make(map[string]string)
}
