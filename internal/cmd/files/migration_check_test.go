package files

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestMigrationCheckScript(t *testing.T) {
	t.Parallel()

	t.Run("No Apps", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{
			Apps: []appdef.App{},
		})

		err := MigrationCheckScript(t.Context(), input)
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

		err := MigrationCheckScript(t.Context(), input)
		assert.NoError(t, err)

		// Verify no scripts were created.
		for _, app := range appDef.Apps {
			scriptPath := filepath.Join(app.Path, "scripts", "check-deps.cjs")
			exists, err := afero.Exists(input.FS, scriptPath)
			require.NoError(t, err)
			assert.False(t, exists, "expected %s not to exist", scriptPath)
		}
	})

	t.Run("Creates Script For Payload Apps", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "cms", Type: appdef.AppTypePayload, Path: "./apps/cms"},
				{Name: "web", Type: appdef.AppTypeSvelteKit, Path: "./apps/web"},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := MigrationCheckScript(t.Context(), input)
		require.NoError(t, err)

		// Verify script was created for Payload app.
		cmsScriptPath := filepath.Join("./apps/cms/scripts", "check-deps.cjs")
		exists, err := afero.Exists(input.FS, cmsScriptPath)
		require.NoError(t, err)
		assert.True(t, exists, "expected %s to exist", cmsScriptPath)

		// Verify script content.
		content, err := afero.ReadFile(input.FS, cmsScriptPath)
		require.NoError(t, err)
		assert.Contains(t, string(content), "pnpm-lock.yaml")
		assert.Contains(t, string(content), ".pnpm")
		assert.Contains(t, string(content), "Dependencies out of sync")

		// Verify script was not created for SvelteKit app.
		webScriptPath := filepath.Join("./apps/web/scripts", "check-deps.cjs")
		exists, err = afero.Exists(input.FS, webScriptPath)
		require.NoError(t, err)
		assert.False(t, exists, "expected %s not to exist", webScriptPath)
	})

	t.Run("Multiple Payload Apps", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "cms", Type: appdef.AppTypePayload, Path: "./apps/cms"},
				{Name: "admin", Type: appdef.AppTypePayload, Path: "./apps/admin"},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := MigrationCheckScript(t.Context(), input)
		require.NoError(t, err)

		// Verify script was created for both Payload apps.
		for _, app := range appDef.Apps {
			scriptPath := filepath.Join(app.Path, "scripts", "check-deps.cjs")
			exists, err := afero.Exists(input.FS, scriptPath)
			require.NoError(t, err)
			assert.True(t, exists, "expected %s to exist", scriptPath)
		}
	})

	t.Run("Idempotent Execution", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "cms", Type: appdef.AppTypePayload, Path: "./apps/cms"},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		// Run first time.
		err := MigrationCheckScript(t.Context(), input)
		require.NoError(t, err)

		scriptPath := filepath.Join("./apps/cms/scripts", "check-deps.cjs")
		content1, err := afero.ReadFile(input.FS, scriptPath)
		require.NoError(t, err)

		// Run second time - should be idempotent.
		err = MigrationCheckScript(t.Context(), input)
		require.NoError(t, err)

		content2, err := afero.ReadFile(input.FS, scriptPath)
		require.NoError(t, err)

		assert.Equal(t, content1, content2, "script content should remain the same")
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "cms", Type: appdef.AppTypePayload, Path: "./apps/cms"},
			},
		}

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), appDef)

		err := MigrationCheckScript(t.Context(), input)
		assert.Error(t, err)
	})
}
