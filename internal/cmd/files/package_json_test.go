package files

import (
	"encoding/json"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/util/testutil"
)

func TestPackageJSON(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
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

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := PackageJSON(t.Context(), input)
		assert.NoError(t, err)

		t.Log("File Exists")
		{
			exists, err := afero.Exists(input.FS, "package.json")
			assert.NoError(t, err)
			assert.True(t, exists)
		}

		t.Log("Conforms to Schema")
		{
			schema, err := testutil.SchemaFromURL(t, "https://www.schemastore.org/package.json")
			require.NoError(t, err)

			file, err := afero.ReadFile(input.FS, "package.json")
			require.NoError(t, err)

			err = schema.ValidateJSON(file)
			assert.NoError(t, err, "Package.json file conforms to schema")
		}
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), &appdef.Definition{})

		got := PackageJSON(t.Context(), input)
		assert.Error(t, got)
	})
}

func TestPackageJSONApp(t *testing.T) {
	t.Parallel()

	t.Run("Adds scripts to NPM app", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{Name: "my-website"},
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "cms",
					Build: appdef.Build{
						Port: 3000,
					},
				},
			},
		}
		require.NoError(t, appDef.ApplyDefaults())

		input := setup(t, fs, appDef)

		require.NoError(t, afero.WriteFile(fs, "cms/package.json", []byte(`{
	"name": "cms",
	"version": "1.0.0",
	"scripts": {
		"dev": "payload dev",
		"build": "payload build"
	}
}`), 0o644))

		err := PackageJSONApp(t.Context(), input)
		assert.NoError(t, err)

		t.Log("File exists and contains Docker scripts")
		{
			data, err := afero.ReadFile(fs, "cms/package.json")
			require.NoError(t, err)

			var pkg map[string]any
			require.NoError(t, json.Unmarshal(data, &pkg))

			scripts, ok := pkg["scripts"].(map[string]any)
			require.True(t, ok)

			assert.Equal(t, "pnpm docker:build && pnpm docker:run", scripts["docker"])
			assert.Equal(t, "docker build . -t cms-web --progress plain --no-cache", scripts["docker:build"])
			assert.Equal(t, "docker run -it --init --env-file .env -p 3000:3000 --rm -ti cms-web", scripts["docker:run"])
			assert.Equal(t, "docker image rm cms-web", scripts["docker:remove"])
		}

		t.Log("Existing scripts preserved")
		{
			data, err := afero.ReadFile(fs, "cms/package.json")
			require.NoError(t, err)

			var pkg map[string]any
			require.NoError(t, json.Unmarshal(data, &pkg))

			scripts, ok := pkg["scripts"].(map[string]any)
			require.True(t, ok)

			assert.Equal(t, "payload dev", scripts["dev"])
			assert.Equal(t, "payload build", scripts["build"])
		}
	})

	t.Run("Skips non-NPM app", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{Name: "my-website"},
			Apps: []appdef.App{
				{
					Name: "api",
					Type: appdef.AppTypeGoLang,
					Path: "api",
					Build: appdef.Build{
						Port: 8080,
					},
				},
			},
		}
		require.NoError(t, appDef.ApplyDefaults())

		input := setup(t, fs, appDef)

		require.NoError(t, afero.WriteFile(fs, "api/package.json", []byte(`{
	"name": "api",
	"version": "1.0.0"
}`), 0o644))

		err := PackageJSONApp(t.Context(), input)
		assert.NoError(t, err)

		t.Log("Package.json unchanged")
		{
			data, err := afero.ReadFile(fs, "api/package.json")
			require.NoError(t, err)

			var pkg map[string]any
			require.NoError(t, json.Unmarshal(data, &pkg))

			_, hasScripts := pkg["scripts"]
			assert.False(t, hasScripts)
		}
	})

	t.Run("Skips when package.json missing", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{Name: "my-website"},
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "web",
					Build: appdef.Build{
						Port: 3001,
					},
				},
			},
		}
		require.NoError(t, appDef.ApplyDefaults())

		input := setup(t, fs, appDef)

		err := PackageJSONApp(t.Context(), input)
		assert.NoError(t, err)

		exists, _ := afero.Exists(fs, "web/package.json")
		assert.False(t, exists)
	})

	t.Run("Creates scripts section if missing", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{Name: "my-website"},
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "web",
					Build: appdef.Build{
						Port: 3001,
					},
				},
			},
		}
		require.NoError(t, appDef.ApplyDefaults())

		input := setup(t, fs, appDef)

		require.NoError(t, afero.WriteFile(fs, "web/package.json", []byte(`{
	"name": "web",
	"version": "1.0.0"
}`), 0o644))

		err := PackageJSONApp(t.Context(), input)
		assert.NoError(t, err)

		t.Log("Scripts section created with Docker scripts")
		{
			data, err := afero.ReadFile(fs, "web/package.json")
			require.NoError(t, err)

			var pkg map[string]any
			require.NoError(t, json.Unmarshal(data, &pkg))

			scripts, ok := pkg["scripts"].(map[string]any)
			require.True(t, ok)
			assert.NotEmpty(t, scripts)
		}
	})

	t.Run("Multiple apps with different ports", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{Name: "my-website"},
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "cms",
					Build: appdef.Build{
						Port: 3000,
					},
				},
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "web",
					Build: appdef.Build{
						Port: 3001,
					},
				},
			},
		}
		require.NoError(t, appDef.ApplyDefaults())

		input := setup(t, fs, appDef)

		require.NoError(t, afero.WriteFile(fs, "cms/package.json", []byte(`{"name": "cms"}`), 0o644))
		require.NoError(t, afero.WriteFile(fs, "web/package.json", []byte(`{"name": "web"}`), 0o644))

		err := PackageJSONApp(t.Context(), input)
		assert.NoError(t, err)

		t.Log("CMS has port 3000")
		{
			data, err := afero.ReadFile(fs, "cms/package.json")
			require.NoError(t, err)

			var pkg map[string]any
			require.NoError(t, json.Unmarshal(data, &pkg))

			scripts := pkg["scripts"].(map[string]any)
			assert.Contains(t, scripts["docker:run"], "3000:3000")
			assert.Contains(t, scripts["docker:build"], "cms-web")
		}

		t.Log("Web has port 3001")
		{
			data, err := afero.ReadFile(fs, "web/package.json")
			require.NoError(t, err)

			var pkg map[string]any
			require.NoError(t, json.Unmarshal(data, &pkg))

			scripts := pkg["scripts"].(map[string]any)
			assert.Contains(t, scripts["docker:run"], "3001:3001")
			assert.Contains(t, scripts["docker:build"], "web-web")
		}
	})
}
