package files

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestDockerIgnore(t *testing.T) {
	t.Parallel()

	t.Run("No Apps", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{
			Apps: []appdef.App{},
		})

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

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := DockerIgnore(t.Context(), input)
		require.NoError(t, err)

		for _, app := range appDef.Apps {
			path := app.Path + "/.dockerignore"
			exists, err := afero.Exists(input.FS, path)
			require.NoError(t, err)
			assert.True(t, exists, "expected %s to exist", path)
		}
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), &appdef.Definition{
			Apps: []appdef.App{
				{Name: "api", Path: "./apps/api"},
			},
		})

		err := DockerIgnore(t.Context(), input)
		assert.Error(t, err)
	})
}
