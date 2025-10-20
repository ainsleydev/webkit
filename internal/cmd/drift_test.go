// internal/cmd/drift_test.go

package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/mocks"
)

func TestDrift(t *testing.T) {
	t.Parallel()

	t.Run("No Drift - No Manifest", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{
			Project: appdef.Project{
				Name: "test",
				Repo: appdef.GitHubRepo{Owner: "test", Name: "test"},
			},
		})

		err := drift(t.Context(), input)
		assert.Error(t, err, "should error when no manifest exists")
	})

	t.Run("Update Error", func(t *testing.T) {
		t.Parallel()
		t.Skip("TODO")
	})

	t.Run("FS Error", func(t *testing.T) {
		t.Parallel()

		mock := mocks.NewMockFS(gomock.NewController(t))
		mock.EXPECT().
			Open(gomock.Any()).
			Return(nil, fmt.Errorf("open error"))

		input := setup(t, mock, &appdef.Definition{})

		err := drift(t.Context(), input)
		assert.Error(t, err)
	})

	t.Run("No Drift - Files Match", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test",
				Repo: appdef.GitHubRepo{Owner: "test", Name: "test"},
			},
		}

		// Run update first to generate all files
		input := setup(t, fs, appDef)
		err := update(t.Context(), input)
		require.NoError(t, err)

		// Now check drift
		input, buf := setupWithPrinter(t, fs, appDef)
		err = drift(t.Context(), input)

		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "No drift detected")
		assert.Contains(t, buf.String(), "all files are up to date")
	})

	t.Run("Drift - Manual Modification", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test",
				Repo: appdef.GitHubRepo{Owner: "test", Name: "test"},
			},
		}

		// Run update to generate files
		input := setup(t, fs, appDef)
		err := update(t.Context(), input)
		require.NoError(t, err)

		// Manually modify a file
		err = afero.WriteFile(fs, ".gitignore", []byte("# User modified"), 0o644)
		require.NoError(t, err)

		// Check drift
		input, buf := setupWithPrinter(t, fs, appDef)
		err = drift(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, buf.String(), "Manual modifications detected")
		assert.Contains(t, buf.String(), ".gitignore")
		assert.Contains(t, buf.String(), "Run 'webkit update' to sync all files")
	})

	t.Run("Drift - Outdated From app.json Change", func(t *testing.T) {
		t.Skip()
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test",
				Repo: appdef.GitHubRepo{Owner: "test", Name: "test"},
			},
		}

		// Run update with initial app.json
		input := setup(t, fs, appDef)
		err := update(t.Context(), input)
		require.NoError(t, err)

		// Change app.json (simulate by changing the definition)
		appDef.Project.Name = "test-renamed"

		// Check drift with new app.json
		input, buf := setupWithPrinter(t, fs, appDef)
		err = drift(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, buf.String(), "Outdated files detected")
		assert.Contains(t, buf.String(), "Run 'webkit update' to sync all files")
	})

	t.Run("Drift - Missing File", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test",
				Repo: appdef.GitHubRepo{Owner: "test", Name: "test"},
			},
		}

		// Run update to generate files
		input := setup(t, fs, appDef)
		err := update(t.Context(), input)
		require.NoError(t, err)

		// Delete a file
		err = fs.Remove(".gitignore")
		require.NoError(t, err)

		// Check drift
		input, buf := setupWithPrinter(t, fs, appDef)
		err = drift(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, buf.String(), "Missing files detected")
		assert.Contains(t, buf.String(), ".gitignore")
		assert.Contains(t, buf.String(), "Run 'webkit update' to sync all files")
	})

	t.Run("Drift - Orphaned File", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test",
				Repo: appdef.GitHubRepo{Owner: "test", Name: "test"},
			},
			Apps: []appdef.App{
				{Name: "cms", Type: appdef.AppTypePayload, Path: "services/cms"},
			},
		}

		// Run update with app
		input := setup(t, fs, appDef)
		err := update(t.Context(), input)
		require.NoError(t, err)

		// Remove app from definition
		appDef.Apps = []appdef.App{}

		// Check drift - should detect orphaned workflow file
		input, buf := setupWithPrinter(t, fs, appDef)
		err = drift(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, buf.String(), "Orphaned files detected")
		assert.Contains(t, buf.String(), "Run 'webkit update' to sync all files")
	})
}
