package cicd

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestPR(t *testing.T) {
	t.Parallel()

	t.Run("Creates Drift Workflow", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{})

		err := PR(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "drift.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		// Verify it's the drift workflow
		content := string(file)
		assert.Contains(t, content, "WebKit Drift Detection")
		assert.Contains(t, content, "drift-detection:")
	})

	t.Run("No Apps - Still Creates Drift", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := PR(t.Context(), input)
		require.NoError(t, err)

		// Should still create drift workflow
		exists, err := afero.Exists(input.FS, filepath.Join(workflowsPath, "drift.yaml"))
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Creates App Workflows", func(t *testing.T) {
		t.Parallel()

		tt := map[string]struct {
			input appdef.App
			want  string
		}{
			"Javascript": {
				input: appdef.App{
					Name:  "cms",
					Title: "CMS",
					Path:  "./cms",
					Type:  appdef.AppTypePayload,
				},
				want: ".github/workflows/pr-cms.yaml",
			},
			"Go": {
				input: appdef.App{
					Name:  "web",
					Title: "Web",
					Path:  "./web",
					Type:  appdef.AppTypeGoLang,
				},
				want: ".github/workflows/pr-web.yaml",
			},
		}

		for name, test := range tt {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				appDef := &appdef.Definition{Apps: []appdef.App{test.input}}
				input := setup(t, afero.NewMemMapFs(), appDef)
				require.NoError(t, appDef.ApplyDefaults())

				err := PR(t.Context(), input)
				require.NoError(t, err)

				// Check drift workflow exists
				driftFile, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "drift.yaml"))
				require.NoError(t, err)
				assert.Contains(t, string(driftFile), "WebKit Drift Detection")

				// Check app workflow
				file, err := afero.ReadFile(input.FS, test.want)
				require.NoError(t, err)

				err = validateGithubYaml(t, file, false)
				assert.NoError(t, err)

				t.Log("Commands are in order")
				{
					content := string(file)

					// Get positions for each command in the canonical order
					var positions []int
					for _, cmd := range appdef.Commands {
						pos := strings.Index(content, "name: "+strings.Title(cmd.String()))
						if pos != -1 {
							positions = append(positions, pos)
						}
					}

					// Verify positions are in ascending order
					for i := 0; i < len(positions)-1; i++ {
						assert.Less(t, positions[i], positions[i+1],
							"commands should appear in order defined by appdef.Commands")
					}
				}
			})
		}
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:  "web",
					Title: "Web",
					Path:  "./web",
					Type:  appdef.AppTypeGoLang,
				},
			},
		}

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), appDef)

		err := PR(t.Context(), input)
		assert.Error(t, err)
	})
}
