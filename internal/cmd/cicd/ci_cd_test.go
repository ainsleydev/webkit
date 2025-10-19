package cicd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/util/executil"
	"github.com/ainsleydev/webkit/internal/util/testutil"
)

func TestCreatePRWorkflow(t *testing.T) {
	t.Parallel()

	if !executil.Exists("action-validator") {
		t.Skip("action-validator CLI not found in PATH; skipping integration test")
	}

	t.Run("PRs", func(t *testing.T) {
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

				fs := afero.NewMemMapFs()
				def := &appdef.Definition{Apps: []appdef.App{test.input}}
				require.NoError(t, def.ApplyDefaults())

				err := CreatePRWorkflow(t.Context(), cmdtools.CommandInput{
					FS:          fs,
					AppDefCache: def,
				})
				require.NoError(t, err)

				file, err := afero.ReadFile(fs, test.want)
				require.NoError(t, err)

				t.Log("YAML is valid")
				{
					err = testutil.ValidateYAML(t, file)
					assert.NoError(t, err)
				}

				content := string(file)
				fmt.Print(string(content))

				t.Log("Github Action is validated")
				{
					err = testutil.ValidateGithubAction(t, file)
					require.NoError(t, err)
				}

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
}
