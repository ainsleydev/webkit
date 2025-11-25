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
	"github.com/ainsleydev/webkit/internal/state/manifest"
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
			Stat(gomock.Any()).
			Return(nil, fmt.Errorf("stat error")).
			AnyTimes()

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

	//t.Run("Format Error", func(t *testing.T) {
	//	t.Parallel()
	//
	//	fs := afero.NewMemMapFs()
	//	appDef := &appdef.Definition{
	//		Project: appdef.Project{
	//			Name: "test",
	//			Repo: appdef.GitHubRepo{Owner: "test", Name: "test"},
	//		},
	//	}
	//
	//	input := setup(t, fs, appDef)
	//	err := update(t.Context(), input)
	//	require.NoError(t, err)
	//
	//	input.Command = &cli.Command{
	//		Flags: []cli.Flag{
	//			&cli.StringFlag{
	//				Name:  "format",
	//				Value: "wrong",
	//			},
	//		},
	//	}
	//
	//	err = drift(t.Context(), input)
	//	//fmt.Println(err.Error())
	//	assert.Error(t, err)
	//})

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

	t.Run("No Drift - With Custom AGENTS.md Content", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test",
				Repo: appdef.GitHubRepo{Owner: "test", Name: "test"},
			},
		}

		customContent := "Custom Project Rules\n\nThese are project-specific guidelines."
		err := afero.WriteFile(fs, "docs/AGENTS.md", []byte(customContent), 0o644)
		require.NoError(t, err)

		input := setup(t, fs, appDef)
		err = update(t.Context(), input)
		require.NoError(t, err)

		agentsContent, err := afero.ReadFile(fs, "AGENTS.md")
		require.NoError(t, err)
		assert.Contains(t, string(agentsContent), customContent)

		input, buf := setupWithPrinter(t, fs, appDef)
		err = drift(t.Context(), input)

		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "No drift detected")
		assert.Contains(t, buf.String(), "all files are up to date")
	})

	t.Run("Drift - Modified Custom AGENTS.md Content", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test",
				Repo: appdef.GitHubRepo{Owner: "test", Name: "test"},
			},
		}

		customContent := "Custom Project Rules\n\nThese are project-specific guidelines."
		err := afero.WriteFile(fs, "docs/AGENTS.md", []byte(customContent), 0o644)
		require.NoError(t, err)

		input := setup(t, fs, appDef)
		err = update(t.Context(), input)
		require.NoError(t, err)

		newCustomContent := "Updated Rules\n\nThese rules have been updated."
		err = afero.WriteFile(fs, "docs/AGENTS.md", []byte(newCustomContent), 0o644)
		require.NoError(t, err)

		input, buf := setupWithPrinter(t, fs, appDef)
		err = drift(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, buf.String(), "Outdated files detected")
		assert.Contains(t, buf.String(), "AGENTS.md")
		assert.Contains(t, buf.String(), "Template or configuration changed")
	})

	t.Run("No Drift - With Custom README.md Content", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test",
				Repo: appdef.GitHubRepo{Owner: "test", Name: "test"},
			},
		}

		customContent := "## Custom Section\n\nCustom README content goes here."
		err := afero.WriteFile(fs, "docs/README.md", []byte(customContent), 0o644)
		require.NoError(t, err)

		input := setup(t, fs, appDef)
		err = update(t.Context(), input)
		require.NoError(t, err)

		readmeContent, err := afero.ReadFile(fs, "README.md")
		require.NoError(t, err)
		assert.Contains(t, string(readmeContent), customContent)

		input, buf := setupWithPrinter(t, fs, appDef)
		err = drift(t.Context(), input)

		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "No drift detected")
		assert.Contains(t, buf.String(), "all files are up to date")
	})

	t.Run("Drift - Modified Custom README.md Content", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test",
				Repo: appdef.GitHubRepo{Owner: "test", Name: "test"},
			},
		}

		customContent := "## Custom Section\n\nOriginal README content."
		err := afero.WriteFile(fs, "docs/README.md", []byte(customContent), 0o644)
		require.NoError(t, err)

		input := setup(t, fs, appDef)
		err = update(t.Context(), input)
		require.NoError(t, err)

		newCustomContent := "## Updated Section\n\nThis content has been updated."
		err = afero.WriteFile(fs, "docs/README.md", []byte(newCustomContent), 0o644)
		require.NoError(t, err)

		input, buf := setupWithPrinter(t, fs, appDef)
		err = drift(t.Context(), input)

		assert.Error(t, err)
		assert.Contains(t, buf.String(), "Outdated files detected")
		assert.Contains(t, buf.String(), "README.md")
		assert.Contains(t, buf.String(), "Template or configuration changed")
	})

	t.Run("No Drift - With outputs.json", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test",
				Repo: appdef.GitHubRepo{Owner: "test", Name: "test"},
			},
		}

		// Create outputs.json with monitoring data (simulating Terraform output)
		outputsJSON := `{
			"peekaping": {
				"endpoint": "https://uptime.example.com",
				"project_tag": "test-project-123"
			},
			"monitors": [
				{"id": "mon123", "name": "HTTP - example.com", "type": "http"}
			],
			"slack": {"channel_name": "alerts", "channel_id": "C123"}
		}`
		err := fs.MkdirAll(".webkit", 0o755)
		require.NoError(t, err)
		err = afero.WriteFile(fs, ".webkit/outputs.json", []byte(outputsJSON), 0o644)
		require.NoError(t, err)

		// Run update to generate README.md with status badges
		input := setup(t, fs, appDef)
		err = update(t.Context(), input)
		require.NoError(t, err)

		// Verify README contains status section
		readmeContent, err := afero.ReadFile(fs, "README.md")
		require.NoError(t, err)
		assert.Contains(t, string(readmeContent), "Status")

		// Check drift - should detect no drift because outputs.json is copied
		input, buf := setupWithPrinter(t, fs, appDef)
		err = drift(t.Context(), input)

		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "No drift detected")
		assert.Contains(t, buf.String(), "all files are up to date")
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
		assert.Contains(t, output, "Template or configuration changed")
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

	t.Run("Includes generator and source info", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{
				Path:      "README.md",
				Type:      manifest.DriftReasonOutdated,
				Source:    "project",
				Generator: "readme",
			},
			{
				Path:      ".github/workflows/pr.yaml",
				Type:      manifest.DriftReasonModified,
				Source:    "app:cms",
				Generator: "cicd",
			},
		}

		output := formatDriftAsText(drifted)
		assert.Contains(t, output, "README.md")
		assert.Contains(t, output, "generated by readme from project config")
		assert.Contains(t, output, ".github/workflows/pr.yaml")
		assert.Contains(t, output, "generated by cicd from app:cms config")
	})
}

func TestFormatDriftAsMarkdown(t *testing.T) {
	t.Parallel()

	t.Run("No drift", func(t *testing.T) {
		t.Parallel()

		output := formatDriftAsMarkdown([]manifest.DriftEntry{})
		assert.Contains(t, output, "Drift Detection")
		assert.Contains(t, output, "No drift detected")
		assert.Contains(t, output, "all files are up to date")
	})

	t.Run("Modified files - singular", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{Path: "file1.txt", Type: manifest.DriftReasonModified},
		}

		output := formatDriftAsMarkdown(drifted)
		assert.Contains(t, output, "Drift Detection")
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
		assert.Contains(t, output, "Template or configuration changed")
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
		assert.Contains(t, output, "Manual modifications detected")
		assert.Contains(t, output, "Outdated files detected")
		assert.Contains(t, output, "Missing files detected")
		assert.Contains(t, output, "Orphaned files detected")
	})

	t.Run("Includes generator and source info", func(t *testing.T) {
		t.Parallel()

		drifted := []manifest.DriftEntry{
			{
				Path:      "README.md",
				Type:      manifest.DriftReasonOutdated,
				Source:    "project",
				Generator: "readme",
			},
			{
				Path:      ".github/workflows/pr.yaml",
				Type:      manifest.DriftReasonModified,
				Source:    "app:cms",
				Generator: "cicd",
			},
		}

		output := formatDriftAsMarkdown(drifted)
		assert.Contains(t, output, "`README.md`")
		assert.Contains(t, output, "generated by **readme** from **project** config")
		assert.Contains(t, output, "`.github/workflows/pr.yaml`")
		assert.Contains(t, output, "generated by **cicd** from **app:cms** config")
	})
}

func TestCopyUserFiles(t *testing.T) {
	t.Parallel()

	t.Run("Copies AGENTS.md when it exists", func(t *testing.T) {
		t.Parallel()

		srcFS := afero.NewMemMapFs()
		dstFS := afero.NewMemMapFs()

		customContent := "Custom Content\n\nProject-specific rules."
		err := afero.WriteFile(srcFS, "docs/AGENTS.md", []byte(customContent), 0o644)
		require.NoError(t, err)

		err = copyUserFiles(srcFS, dstFS)
		require.NoError(t, err)

		exists, err := afero.Exists(dstFS, "docs/AGENTS.md")
		require.NoError(t, err)
		assert.True(t, exists)

		copiedContent, err := afero.ReadFile(dstFS, "docs/AGENTS.md")
		require.NoError(t, err)
		assert.Equal(t, customContent, string(copiedContent))
	})

	t.Run("Copies AGENTS.md.tmpl when it exists", func(t *testing.T) {
		t.Parallel()

		srcFS := afero.NewMemMapFs()
		dstFS := afero.NewMemMapFs()

		templateContent := "App: {{ .Definition.Project.Name }}"
		err := afero.WriteFile(srcFS, "docs/AGENTS.md.tmpl", []byte(templateContent), 0o644)
		require.NoError(t, err)

		err = copyUserFiles(srcFS, dstFS)
		require.NoError(t, err)

		exists, err := afero.Exists(dstFS, "docs/AGENTS.md.tmpl")
		require.NoError(t, err)
		assert.True(t, exists)

		copiedContent, err := afero.ReadFile(dstFS, "docs/AGENTS.md.tmpl")
		require.NoError(t, err)
		assert.Equal(t, templateContent, string(copiedContent))
	})

	t.Run("No error when files don't exist", func(t *testing.T) {
		t.Parallel()

		srcFS := afero.NewMemMapFs()
		dstFS := afero.NewMemMapFs()

		err := copyUserFiles(srcFS, dstFS)
		assert.NoError(t, err)
	})

	t.Run("Copies both files when both exist", func(t *testing.T) {
		t.Parallel()

		srcFS := afero.NewMemMapFs()
		dstFS := afero.NewMemMapFs()

		err := afero.WriteFile(srcFS, "docs/AGENTS.md", []byte("static content"), 0o644)
		require.NoError(t, err)
		err = afero.WriteFile(srcFS, "docs/AGENTS.md.tmpl", []byte("template content"), 0o644)
		require.NoError(t, err)

		err = copyUserFiles(srcFS, dstFS)
		require.NoError(t, err)

		exists1, err := afero.Exists(dstFS, "docs/AGENTS.md")
		require.NoError(t, err)
		exists2, err := afero.Exists(dstFS, "docs/AGENTS.md.tmpl")
		require.NoError(t, err)
		assert.True(t, exists1)
		assert.True(t, exists2)
	})

	t.Run("Copies README.md when it exists", func(t *testing.T) {
		t.Parallel()

		srcFS := afero.NewMemMapFs()
		dstFS := afero.NewMemMapFs()

		customContent := "Custom README Content\n\nProject-specific information."
		err := afero.WriteFile(srcFS, "docs/README.md", []byte(customContent), 0o644)
		require.NoError(t, err)

		err = copyUserFiles(srcFS, dstFS)
		require.NoError(t, err)

		exists, err := afero.Exists(dstFS, "docs/README.md")
		require.NoError(t, err)
		assert.True(t, exists)

		copiedContent, err := afero.ReadFile(dstFS, "docs/README.md")
		require.NoError(t, err)
		assert.Equal(t, customContent, string(copiedContent))
	})

	t.Run("Copies outputs.json when it exists", func(t *testing.T) {
		t.Parallel()

		srcFS := afero.NewMemMapFs()
		dstFS := afero.NewMemMapFs()

		outputsContent := `{"peekaping":{"endpoint":"https://uptime.example.com"}}`
		err := srcFS.MkdirAll(".webkit", 0o755)
		require.NoError(t, err)
		err = afero.WriteFile(srcFS, ".webkit/outputs.json", []byte(outputsContent), 0o644)
		require.NoError(t, err)

		err = copyUserFiles(srcFS, dstFS)
		require.NoError(t, err)

		exists, err := afero.Exists(dstFS, ".webkit/outputs.json")
		require.NoError(t, err)
		assert.True(t, exists)

		copiedContent, err := afero.ReadFile(dstFS, ".webkit/outputs.json")
		require.NoError(t, err)
		assert.Equal(t, outputsContent, string(copiedContent))
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
