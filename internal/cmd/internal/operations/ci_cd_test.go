package operations

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/testutil"
)

func TestCreateCICD(t *testing.T) {
	t.Parallel()

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
					Type:  appdef.AppTypePayload,
				},
				want: ".github/workflows/app-pr-cms.yaml",
			},
			"Go": {
				input: appdef.App{
					Name:  "web",
					Title: "Web",
					Type:  appdef.AppTypeGoLang,
				},
				want: ".github/workflows/app-pr-web.yaml",
			},
		}

		for name, test := range tt {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				fs := afero.NewMemMapFs()

				def := &appdef.Definition{
					Apps: []appdef.App{test.input},
				}

				err := CreateCICD(t.Context(), cmdtools.CommandInput{
					FS:          fs,
					AppDefCache: def,
				})
				if err != nil {
					return
				}

				file, err := afero.ReadFile(fs, test.want)
				require.NoError(t, err)

				err = testutil.ValidateYAML(t, file)
				assert.NoError(t, err)

				err = testutil.ValidateGithubAction(t, file)
				require.NoError(t, err)
			})
		}
	})
}
