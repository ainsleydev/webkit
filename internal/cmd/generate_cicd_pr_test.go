package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	cmdtools "github.com/ainsleydev/webkit/internal/cmd/internal"
	"github.com/ainsleydev/webkit/internal/testutil"
)

func Test_CreateCICD(t *testing.T) {
	t.Parallel()

	t.Run("Javascript", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		def := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:  "cms",
					Title: "CMS",
					Type:  appdef.AppTypePayload,
				},
			},
		}

		err := createCICD(t.Context(), cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: def,
		})
		if err != nil {
			return
		}

		file, err := afero.ReadFile(fs, ".github/workflows/pr-cms.yaml")
		require.NoError(t, err)

		err = testutil.ValidateYAML(t, file)
		assert.NoError(t, err)

		err = testutil.ValidateGithubAction(t, file)
		require.NoError(t, err)
	})

	t.Run("Go", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		def := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:  "web",
					Title: "Web",
					Type:  appdef.AppTypeGoLang,
				},
			},
		}

		err := createCICD(t.Context(), cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: def,
		})
		if err != nil {
			return
		}

		file, err := afero.ReadFile(fs, ".github/workflows/pr-web.yaml")
		require.NoError(t, err)

		err = testutil.ValidateYAML(t, file)
		assert.NoError(t, err)

		err = testutil.ValidateGithubAction(t, file)
		require.NoError(t, err)
	})
}
