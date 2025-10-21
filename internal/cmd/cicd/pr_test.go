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

	t.Run("Creates Workflow", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:  "cms",
					Title: "CMS",
					Path:  "./cms",
					Type:  appdef.AppTypePayload,
				},
				{
					Name:  "web",
					Title: "Web",
					Path:  "./web",
					Type:  appdef.AppTypeGoLang,
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := PR(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "pr.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		content := string(file)

		t.Log("Drift")
		{

			assert.Contains(t, content, "WebKit Drift Detection")
			assert.Contains(t, content, "drift-detection:")
		}

		t.Log("Apps")
		{
			for _, app := range appDef.Apps {
				jobName := strings.ToLower(app.Name)
				assert.Contains(t, content, jobName+":", "workflow should contain job for app "+app.Name)

				switch app.Type {
				case appdef.AppTypeGoLang:
					assert.Contains(t, content, "Set up Go", "Go app should have Go setup")
				case appdef.AppTypePayload:
					assert.Contains(t, content, "Install pnpm", "JS app should have pnpm setup")
					assert.Contains(t, content, "Set up Node", "JS app should have Node setup")
				}
			}
		}

		t.Log("Commands")
		{
			// Get positions for each command in the canonical order.
			var positions []int
			for _, cmd := range appdef.Commands {
				pos := strings.Index(content, "name: "+strings.Title(cmd.String()))
				if pos != -1 {
					positions = append(positions, pos)
				}
			}

			// Verify positions are in ascending order.
			for i := 0; i < len(positions)-1; i++ {
				assert.Less(t, positions[i], positions[i+1],
					"commands should appear in order defined by appdef.Commands")
			}
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
