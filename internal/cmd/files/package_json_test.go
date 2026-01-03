package files

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/fsext"
	"github.com/ainsleydev/webkit/internal/mocks"
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

	t.Run("HTML characters not escaped", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{Name: "my-website", Description: "Test with > and < characters"},
			Apps:    []appdef.App{},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := PackageJSON(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, "package.json")
		require.NoError(t, err)

		fileContent := string(file)
		assert.Contains(t, fileContent, "Test with > and < characters", "HTML characters should not be escaped")
		assert.NotContains(t, fileContent, "\\u003e", "Should not contain escaped >")
		assert.NotContains(t, fileContent, "\\u003c", "Should not contain escaped <")
	})

	t.Run("Field ordering correct", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name:        "my-website",
				Description: "My project description",
			},
			Apps: []appdef.App{},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := PackageJSON(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, "package.json")
		require.NoError(t, err)

		var pkg map[string]any
		require.NoError(t, json.Unmarshal(file, &pkg))

		t.Log("Verify field order in raw JSON")
		{
			fileContent := string(file)
			nameIdx := bytes.Index(file, []byte(`"name"`))
			descIdx := bytes.Index(file, []byte(`"description"`))
			licenseIdx := bytes.Index(file, []byte(`"license"`))
			privateIdx := bytes.Index(file, []byte(`"private"`))
			typeIdx := bytes.Index(file, []byte(`"type"`))
			versionIdx := bytes.Index(file, []byte(`"version"`))

			assert.Greater(t, descIdx, nameIdx, "description should come after name")
			assert.Greater(t, licenseIdx, descIdx, "license should come after description")
			assert.Greater(t, privateIdx, licenseIdx, "private should come after license")
			assert.Greater(t, typeIdx, privateIdx, "type should come after private")
			assert.Greater(t, versionIdx, typeIdx, "version should come after type")

			_ = fileContent
		}
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
			assert.Equal(t, "docker build . -t my-website-cms --progress plain --no-cache", scripts["docker:build"])
			assert.Equal(t, "docker run -it --init --env-file .env -p 3000:3000 --rm -ti my-website-cms", scripts["docker:run"])
			assert.Equal(t, "docker image rm my-website-cms", scripts["docker:remove"])
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

		exists := fsext.Exists(fs, "web/package.json")
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
			assert.Contains(t, scripts["docker:build"], "my-website-cms")
		}

		t.Log("Web has port 3001")
		{
			data, err := afero.ReadFile(fs, "web/package.json")
			require.NoError(t, err)

			var pkg map[string]any
			require.NoError(t, json.Unmarshal(data, &pkg))

			scripts := pkg["scripts"].(map[string]any)
			assert.Contains(t, scripts["docker:run"], "3001:3001")
			assert.Contains(t, scripts["docker:build"], "my-website-web")
		}
	})

	t.Run("Exists check error", func(t *testing.T) {
		mock := mocks.NewMockFS(gomock.NewController(t))
		mock.EXPECT().Stat(gomock.Any()).Return(nil, errors.New("stat error"))

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

		input := setup(t, mock, appDef)

		err := PackageJSONApp(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "checking if")
	})

	t.Run("Read file error", func(t *testing.T) {
		mock := mocks.NewMockFS(gomock.NewController(t))
		mock.EXPECT().Stat(gomock.Any()).Return(nil, nil)
		mock.EXPECT().Open(gomock.Any()).Return(nil, errors.New("read error"))

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

		input := setup(t, mock, appDef)

		err := PackageJSONApp(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "reading")
	})

	t.Run("Invalid JSON error", func(t *testing.T) {
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

		require.NoError(t, afero.WriteFile(fs, "cms/package.json", []byte(`{invalid json`), 0o644))

		err := PackageJSONApp(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "parsing")
	})

	t.Run("Write error", func(t *testing.T) {
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

		require.NoError(t, afero.WriteFile(fs, "cms/package.json", []byte(`{"name": "cms"}`), 0o644))

		input := setup(t, afero.NewReadOnlyFs(fs), appDef)

		err := PackageJSONApp(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "writing")
	})

	t.Run("Payload app includes migration scripts", func(t *testing.T) {
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
		"dev": "payload dev"
	}
}`), 0o644))

		err := PackageJSONApp(t.Context(), input)
		assert.NoError(t, err)

		t.Log("Payload migration scripts added")
		{
			data, err := afero.ReadFile(fs, "cms/package.json")
			require.NoError(t, err)

			var pkg map[string]any
			require.NoError(t, json.Unmarshal(data, &pkg))

			scripts, ok := pkg["scripts"].(map[string]any)
			require.True(t, ok)

			assert.Equal(t, "NODE_ENV=production payload migrate:create", scripts["migrate:create"])
			assert.Equal(t, "NODE_ENV=production payload migrate:status", scripts["migrate:status"])
		}
	})

	t.Run("SvelteKit app does not include migration scripts", func(t *testing.T) {
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
	"version": "1.0.0",
	"scripts": {
		"dev": "vite dev"
	}
}`), 0o644))

		err := PackageJSONApp(t.Context(), input)
		assert.NoError(t, err)

		t.Log("No migration scripts for SvelteKit")
		{
			data, err := afero.ReadFile(fs, "web/package.json")
			require.NoError(t, err)

			var pkg map[string]any
			require.NoError(t, json.Unmarshal(data, &pkg))

			scripts, ok := pkg["scripts"].(map[string]any)
			require.True(t, ok)

			_, hasMigrateCreate := scripts["migrate:create"]
			_, hasMigrateStatus := scripts["migrate:status"]
			assert.False(t, hasMigrateCreate)
			assert.False(t, hasMigrateStatus)
		}
	})
}

func TestGetAppTypeScripts(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input appdef.AppType
		want  map[string]string
	}{
		"Payload includes migration scripts": {
			input: appdef.AppTypePayload,
			want: map[string]string{
				"migrate:create": "NODE_ENV=production payload migrate:create",
				"migrate:status": "NODE_ENV=production payload migrate:status",
			},
		},
		"SvelteKit returns empty map": {
			input: appdef.AppTypeSvelteKit,
			want:  map[string]string{},
		},
		"GoLang returns empty map": {
			input: appdef.AppTypeGoLang,
			want:  map[string]string{},
		},
		"Unknown type returns empty map": {
			input: appdef.AppType("unknown"),
			want:  map[string]string{},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := getAppTypeScripts(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
