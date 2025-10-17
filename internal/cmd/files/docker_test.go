package files

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
)

func TestDockerIgnore(t *testing.T) {
	t.Parallel()

	t.Run("No Apps", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{},
		}

		input := cmdtools.CommandInput{
			FS:          afero.NewMemMapFs(),
			AppDefCache: appDef,
		}

		err := DockerIgnore(t.Context(), input)
		assert.NoError(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "cms", Path: "./apps/cms"},
				{Name: "api", Path: "./apps/api"},
			},
		}

		fs := afero.NewMemMapFs()

		input := cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: appDef,
		}

		err := DockerIgnore(t.Context(), input)
		require.NoError(t, err)

		for _, app := range appDef.Apps {
			path := app.Path + "/.dockerignore"
			exists, err := afero.Exists(fs, path)
			require.NoError(t, err)
			assert.True(t, exists, "expected %s to exist", path)
		}
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "api", Path: "./apps/api"},
			},
		}

		input := cmdtools.CommandInput{
			FS:          afero.NewReadOnlyFs(afero.NewMemMapFs()),
			AppDefCache: appDef,
		}

		err := DockerIgnore(t.Context(), input)
		assert.Error(t, err)
	})
}
