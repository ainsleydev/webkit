package operations

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/testutil"
)

func TestCreatePackageJson(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()

	appDef := &appdef.Definition{
		Project: appdef.Project{Name: "my-website"},
		Apps: []appdef.App{
			{
				Name: "cms",
				Type: appdef.AppTypePayload,
				Path: "cms",
			},
		},
	}

	err := CreatePackageJson(t.Context(), cmdtools.CommandInput{
		FS:          fs,
		AppDefCache: appDef,
	})
	assert.NoError(t, err)

	t.Log("File Exists")
	{
		exists, err := afero.Exists(fs, "package.json")
		assert.NoError(t, err)
		assert.True(t, exists)
	}

	t.Log("Conforms to Schema")
	{
		schema, err := testutil.SchemaFromURL(t, "https://www.schemastore.org/package.json")
		require.NoError(t, err)

		file, err := afero.ReadFile(fs, "package.json")
		require.NoError(t, err)

		fmt.Print(string(file))

		err = schema.ValidateJSON(file)
		assert.NoError(t, err, "Package.json file conforms to schema")
	}
}

const file = `{
	"$schema": "https://json.schemastore.org/package.json",
	"name": "my-website",
	"version": "1.0.0",
	"private": false,
	"license": "BSD-3-Clause",
	"type": "module",
	"scripts": {
		"format": "prettier --write .",
		"lint": "eslint .",
		"lint:fix": "eslint . --fix",
		"preinstall": "npx only-allow pnpm",
		"test": "turbo test"
	},
	"devDependencies": {
		"@ainsleydev/eslint-config": "^0.0.6",
		"@ainsleydev/prettier-config": "^0.0.2",
		"@eslint/compat": "^1.4.0",
		"@payloadcms/eslint-config": "^3.28.0",
		"@payloadcms/eslint-plugin": "^3.28.0",
		"eslint": "^9.37.0",
		"eslint-plugin-perfectionist": "^4.15.1",
		"globals": "^16.0.0",
		"prettier": "^3.6.0",
		"prettier-plugin-svelte": "^3.4.0",
		"turbo": "^2.5.8",
		"typescript": "5.8.2"
	},
	"packageManager": "pnpm@10.15.0",
	"author": {
		"name": "ainsley.dev LTD",
		"email": "hello@ainsley.dev",
		"url": "https://ainsley.dev"
	},
	"maintainers": [
		{
			"name": "Ainsley Clark",
			"email": "hello@ainsley.dev",
			"url": "https://ainsley.dev"
		}
	],
	"bundledDependencies": false
}`

func TestCreatePackageJsonFromFile(t *testing.T) {
	t.Parallel()

	schema, err := testutil.SchemaFromURL(t, "https://www.schemastore.org/package.json")
	require.NoError(t, err)

	err = schema.ValidateJSON([]byte(file))
	assert.NoError(t, err, "Package.json file conforms to schema")
}
