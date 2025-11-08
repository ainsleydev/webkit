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

	t.Run("Empty folder with only gitkeep", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "cms", Type: appdef.AppTypePayload, Path: "./apps/cms"},
			},
		}

		fs := afero.NewMemMapFs()
		input := setup(t, fs, appDef)

		// Create empty public folder with only .gitkeep.
		publicPath := filepath.Join("./apps/cms/public")
		err := fs.MkdirAll(publicPath, 0o755)
		require.NoError(t, err)

		gitkeepPath := filepath.Join(publicPath, ".gitkeep")
		err = afero.WriteFile(fs, gitkeepPath, []byte{}, 0o644)
		require.NoError(t, err)

		err = PublicFolder(t.Context(), input)
		require.NoError(t, err)

		// Verify .gitkeep still exists.
		exists, err := afero.Exists(input.FS, gitkeepPath)
		require.NoError(t, err)
		assert.True(t, exists, "expected .gitkeep to exist in empty folder")
	})

	t.Run("Idempotent execution", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "cms", Type: appdef.AppTypePayload, Path: "./apps/cms"},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		// Run first time.
		err := PublicFolder(t.Context(), input)
		require.NoError(t, err)

		gitkeepPath := filepath.Join("./apps/cms/public", ".gitkeep")
		exists, err := afero.Exists(input.FS, gitkeepPath)
		require.NoError(t, err)
		assert.True(t, exists)

		// Run second time - should be idempotent.
		err = PublicFolder(t.Context(), input)
		require.NoError(t, err)

		exists, err = afero.Exists(input.FS, gitkeepPath)
		require.NoError(t, err)
		assert.True(t, exists)
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

func TestFolderHasFiles(t *testing.T) {
	t.Parallel()

	t.Run("Non-existent folder", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		hasFiles, err := folderHasFiles(fs, "./non-existent")
		require.NoError(t, err)
		assert.False(t, hasFiles)
	})

	t.Run("Empty folder", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll("./empty", 0o755)
		require.NoError(t, err)

		hasFiles, err := folderHasFiles(fs, "./empty")
		require.NoError(t, err)
		assert.False(t, hasFiles)
	})

	t.Run("Folder with only gitkeep", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll("./with-gitkeep", 0o755)
		require.NoError(t, err)

		gitkeepPath := filepath.Join("./with-gitkeep", ".gitkeep")
		err = afero.WriteFile(fs, gitkeepPath, []byte{}, 0o644)
		require.NoError(t, err)

		hasFiles, err := folderHasFiles(fs, "./with-gitkeep")
		require.NoError(t, err)
		assert.False(t, hasFiles, "folder with only .gitkeep should be considered empty")
	})

	t.Run("Folder with files", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll("./with-files", 0o755)
		require.NoError(t, err)

		filePath := filepath.Join("./with-files", "file.txt")
		err = afero.WriteFile(fs, filePath, []byte("content"), 0o644)
		require.NoError(t, err)

		hasFiles, err := folderHasFiles(fs, "./with-files")
		require.NoError(t, err)
		assert.True(t, hasFiles)
	})

	t.Run("Folder with gitkeep and other files", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll("./mixed", 0o755)
		require.NoError(t, err)

		gitkeepPath := filepath.Join("./mixed", ".gitkeep")
		err = afero.WriteFile(fs, gitkeepPath, []byte{}, 0o644)
		require.NoError(t, err)

		filePath := filepath.Join("./mixed", "file.txt")
		err = afero.WriteFile(fs, filePath, []byte("content"), 0o644)
		require.NoError(t, err)

		hasFiles, err := folderHasFiles(fs, "./mixed")
		require.NoError(t, err)
		assert.True(t, hasFiles, "folder with .gitkeep and other files should have files")
	})
}
