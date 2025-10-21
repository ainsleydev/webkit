package files

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func TestPnpmWorkspace(t *testing.T) {
	t.Parallel()

	t.Run("No Apps", func(t *testing.T) {
		t.Context()

		appDef := &appdef.Definition{
			Apps: []appdef.App{},
		}

		input := cmdtools.CommandInput{
			FS:          afero.NewMemMapFs(),
			AppDefCache: appDef,
			Manifest:    manifest.NewTracker(),
		}

		got := PnpmWorkspace(t.Context(), input)
		assert.NoError(t, got)
	})

	t.Run("No NPM Apps", func(t *testing.T) {
		t.Context()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:    "cms",
					Type:    appdef.AppTypePayload,
					Path:    "./apps/cms",
					UsesNPM: ptr.BoolPtr(false),
				},
			},
		}

		input := cmdtools.CommandInput{
			FS:          afero.NewMemMapFs(),
			AppDefCache: appDef,
			Manifest:    manifest.NewTracker(),
		}

		got := PnpmWorkspace(t.Context(), input)
		assert.NoError(t, got)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:    "api",
					Type:    appdef.AppTypeGoLang,
					Path:    "./apps/api",
					UsesNPM: ptr.BoolPtr(true), // Go app but explicitly uses NPM
				},
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "./apps/cms",
				},
				{
					Name:    "web",
					Type:    appdef.AppTypeSvelteKit,
					Path:    "./apps/web",
					UsesNPM: ptr.BoolPtr(false), // JS app but explicitly doesn't use NPM
				},
			},
		}

		fs := afero.NewMemMapFs()
		input := cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: appDef,
			Manifest:    manifest.NewTracker(),
		}

		err := PnpmWorkspace(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(fs, "pnpm-workspace.yaml")

		t.Log("Verify File")
		{
			require.NoError(t, err)
			assert.NotEmpty(t, file)
		}

		t.Log("Verify Packages")
		{
			var workspace map[string]any
			err = yaml.Unmarshal(file, &workspace)
			require.NoError(t, err)

			packages, ok := workspace["packages"]
			require.True(t, ok)
			assert.Len(t, packages, 2)
			assert.Contains(t, packages, "./apps/api")    // Go app with UsesNPM: true
			assert.Contains(t, packages, "./apps/cms")    // JS app (default)
			assert.NotContains(t, packages, "./apps/web") // JS app with UsesNPM: false
		}
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "cms", Type: appdef.AppTypePayload, Path: "./apps/cms"},
			},
		}

		input := cmdtools.CommandInput{
			FS:          afero.NewReadOnlyFs(afero.NewMemMapFs()),
			AppDefCache: appDef,
			Manifest:    manifest.NewTracker(),
		}

		got := PnpmWorkspace(t.Context(), input)
		assert.Error(t, got)
	})
}
