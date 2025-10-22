package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/manifest"
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

func TestFormatDriftOutput(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   []manifest.DriftEntry
		format  string
		wantErr bool
	}{
		"Text format": {
			input:   []manifest.DriftEntry{},
			format:  "text",
			wantErr: false,
		},
		"Markdown format": {
			input:   []manifest.DriftEntry{},
			format:  "markdown",
			wantErr: false,
		},
		"JSON format": {
			input:   []manifest.DriftEntry{},
			format:  "json",
			wantErr: false,
		},
		"Invalid format": {
			input:   []manifest.DriftEntry{},
			format:  "invalid",
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := formatDriftOutput(test.input, test.format)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestFormatDriftAsText(t *testing.T) {
	t.Parallel()

	t.Run("No drift", func(t *testing.T) {
		t.Parallel()

		output := formatDriftAsText([]manifest.DriftEntry{})
		assert.Contains(t, output, "No drift detected")
		assert.Contains(t, output, "all files are up to date")
	})

	t.Run("Modified files", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{Path: "file1.txt", Type: manifest.DriftReasonModified},
			{Path: "file2.txt", Type: manifest.DriftReasonModified},
		}

		output := formatDriftAsText(drifted)
		assert.Contains(t, output, "Manual modifications detected")
		assert.Contains(t, output, "file1.txt")
		assert.Contains(t, output, "file2.txt")
		assert.Contains(t, output, "Run 'webkit update' to sync all files")
	})

	t.Run("Outdated files", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{Path: ".env", Type: manifest.DriftReasonOutdated},
		}

		output := formatDriftAsText(drifted)
		assert.Contains(t, output, "Outdated files detected")
		assert.Contains(t, output, ".env")
		assert.Contains(t, output, "app.json changed")
	})

	t.Run("Missing files", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{Path: "missing.txt", Type: manifest.DriftReasonNew},
		}

		output := formatDriftAsText(drifted)
		assert.Contains(t, output, "Missing files detected")
		assert.Contains(t, output, "missing.txt")
	})

	t.Run("Orphaned files", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{Path: "orphan.txt", Type: manifest.DriftReasonDeleted},
		}

		output := formatDriftAsText(drifted)
		assert.Contains(t, output, "Orphaned files detected")
		assert.Contains(t, output, "orphan.txt")
	})

	t.Run("Mixed drift types", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{Path: "modified.txt", Type: manifest.DriftReasonModified},
			{Path: "outdated.txt", Type: manifest.DriftReasonOutdated},
			{Path: "missing.txt", Type: manifest.DriftReasonNew},
			{Path: "orphan.txt", Type: manifest.DriftReasonDeleted},
		}

		output := formatDriftAsText(drifted)
		assert.Contains(t, output, "Manual modifications detected")
		assert.Contains(t, output, "Outdated files detected")
		assert.Contains(t, output, "Missing files detected")
		assert.Contains(t, output, "Orphaned files detected")
	})
}

func TestFormatDriftAsMarkdown(t *testing.T) {
	t.Parallel()

	t.Run("No drift", func(t *testing.T) {
		t.Parallel()

		output := formatDriftAsMarkdown([]manifest.DriftEntry{})
		assert.Contains(t, output, "## WebKit Drift Detection")
		assert.Contains(t, output, "No drift detected")
		assert.Contains(t, output, "all files are up to date")
	})

	t.Run("Modified files - singular", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{Path: "file1.txt", Type: manifest.DriftReasonModified},
		}

		output := formatDriftAsMarkdown(drifted)
		assert.Contains(t, output, "## WebKit Drift Detection")
		assert.Contains(t, output, "Manual modifications detected (1 file)")
		assert.Contains(t, output, "`file1.txt`")
		assert.Contains(t, output, "**Action Required:**")
		assert.Contains(t, output, "`webkit update`")
	})

	t.Run("Modified files - plural", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{Path: "file1.txt", Type: manifest.DriftReasonModified},
			{Path: "file2.txt", Type: manifest.DriftReasonModified},
		}

		output := formatDriftAsMarkdown(drifted)
		assert.Contains(t, output, "Manual modifications detected (2 files)")
		assert.Contains(t, output, "`file1.txt`")
		assert.Contains(t, output, "`file2.txt`")
	})

	t.Run("Outdated files", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{Path: ".env", Type: manifest.DriftReasonOutdated},
			{Path: "config.yaml", Type: manifest.DriftReasonOutdated},
		}

		output := formatDriftAsMarkdown(drifted)
		assert.Contains(t, output, "Outdated files detected (2 files)")
		assert.Contains(t, output, "app.json changed")
		assert.Contains(t, output, "`.env`")
		assert.Contains(t, output, "`config.yaml`")
	})

	t.Run("Missing files", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{Path: "missing.txt", Type: manifest.DriftReasonNew},
		}

		output := formatDriftAsMarkdown(drifted)
		assert.Contains(t, output, "Missing files detected (1 file)")
		assert.Contains(t, output, "`missing.txt`")
	})

	t.Run("Orphaned files", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{Path: "orphan.txt", Type: manifest.DriftReasonDeleted},
		}

		output := formatDriftAsMarkdown(drifted)
		assert.Contains(t, output, "Orphaned files detected (1 file)")
		assert.Contains(t, output, "`orphan.txt`")
	})

	t.Run("Mixed drift types", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{Path: "modified.txt", Type: manifest.DriftReasonModified},
			{Path: "outdated.txt", Type: manifest.DriftReasonOutdated},
			{Path: "missing.txt", Type: manifest.DriftReasonNew},
			{Path: "orphan.txt", Type: manifest.DriftReasonDeleted},
		}

		output := formatDriftAsMarkdown(drifted)
		assert.Contains(t, output, "### ⚠ Manual modifications detected")
		assert.Contains(t, output, "### 📝 Outdated files detected")
		assert.Contains(t, output, "### 📄 Missing files detected")
		assert.Contains(t, output, "### 🗑️ Orphaned files detected")
		assert.Contains(t, output, "---")
	})
}

func TestFormatDriftAsJSON(t *testing.T) {
	t.Parallel()

	t.Run("No drift", func(t *testing.T) {
		t.Parallel()

		output, err := formatDriftAsJSON([]manifest.DriftEntry{})
		require.NoError(t, err)
		assert.Contains(t, output, `"drift_detected": false`)
		assert.Contains(t, output, `"total_files": 0`)
	})

	t.Run("With drift", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{Path: "file1.txt", Type: manifest.DriftReasonModified},
			{Path: "file2.txt", Type: manifest.DriftReasonOutdated},
		}

		output, err := formatDriftAsJSON(drifted)
		require.NoError(t, err)
		assert.Contains(t, output, `"drift_detected": true`)
		assert.Contains(t, output, `"total_files": 2`)
		assert.Contains(t, output, `"file1.txt"`)
		assert.Contains(t, output, `"file2.txt"`)
	})

	t.Run("Summary counts", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{Path: "modified1.txt", Type: manifest.DriftReasonModified},
			{Path: "modified2.txt", Type: manifest.DriftReasonModified},
			{Path: "outdated.txt", Type: manifest.DriftReasonOutdated},
			{Path: "missing.txt", Type: manifest.DriftReasonNew},
			{Path: "orphan.txt", Type: manifest.DriftReasonDeleted},
		}

		output, err := formatDriftAsJSON(drifted)
		require.NoError(t, err)
		assert.Contains(t, output, `"modified": 2`)
		assert.Contains(t, output, `"outdated": 1`)
		assert.Contains(t, output, `"new": 1`)
		assert.Contains(t, output, `"deleted": 1`)
	})
}
