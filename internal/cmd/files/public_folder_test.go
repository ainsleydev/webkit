package files

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestPublicFolder(t *testing.T) {
	t.Parallel()

	t.Run("No Apps", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{
			Apps: []appdef.App{},
		})

		err := PublicFolder(t.Context(), input)
		assert.NoError(t, err)
	})

	t.Run("No Payload Apps", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "web", Type: appdef.AppTypeSvelteKit, Path: "./apps/web"},
				{Name: "api", Type: appdef.AppTypeGoLang, Path: "./apps/api"},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := PublicFolder(t.Context(), input)
		assert.NoError(t, err)

		// Verify no public folders were created.
		for _, app := range appDef.Apps {
			publicPath := filepath.Join(app.Path, "public")
			exists, err := afero.DirExists(input.FS, publicPath)
			require.NoError(t, err)
			assert.False(t, exists, "expected %s not to exist", publicPath)
		}
	})

	t.Run("Creates Public Folder For Payload Apps", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "cms", Type: appdef.AppTypePayload, Path: "./apps/cms"},
				{Name: "web", Type: appdef.AppTypeSvelteKit, Path: "./apps/web"},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := PublicFolder(t.Context(), input)
		require.NoError(t, err)

		// Verify public folder was created for Payload app.
		cmsPublicPath := filepath.Join("./apps/cms/public", ".gitkeep")
		exists, err := afero.Exists(input.FS, cmsPublicPath)
		require.NoError(t, err)
		assert.True(t, exists, "expected %s to exist", cmsPublicPath)

		// Verify public folder was not created for SvelteKit app.
		webPublicPath := filepath.Join("./apps/web/public")
		exists, err = afero.DirExists(input.FS, webPublicPath)
		require.NoError(t, err)
		assert.False(t, exists, "expected %s not to exist", webPublicPath)
	})

	t.Run("Skips Existing Public Folder", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "cms", Type: appdef.AppTypePayload, Path: "./apps/cms"},
			},
		}

		fs := afero.NewMemMapFs()
		input := setup(t, fs, appDef)

		// Create existing public folder with a file.
		publicPath := filepath.Join("./apps/cms/public")
		err := fs.MkdirAll(publicPath, 0o755)
		require.NoError(t, err)

		existingFile := filepath.Join(publicPath, "existing.txt")
		err = afero.WriteFile(fs, existingFile, []byte("existing content"), 0o644)
		require.NoError(t, err)

		err = PublicFolder(t.Context(), input)
		require.NoError(t, err)

		// Verify existing file is still there.
		exists, err := afero.Exists(input.FS, existingFile)
		require.NoError(t, err)
		assert.True(t, exists, "expected existing file to remain")

		// Verify .gitkeep was not created.
		gitkeepPath := filepath.Join(publicPath, ".gitkeep")
		exists, err = afero.Exists(input.FS, gitkeepPath)
		require.NoError(t, err)
		assert.False(t, exists, "expected .gitkeep not to be created when folder exists")
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "cms", Type: appdef.AppTypePayload, Path: "./apps/cms"},
			},
		}

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), appDef)

		err := PublicFolder(t.Context(), input)
		assert.Error(t, err)
	})
}
