package files

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func TestTurboJSON(t *testing.T) {
	t.Parallel()

	t.Run("No Apps", func(t *testing.T) {
		t.Context()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{
			Apps: []appdef.App{},
		})

		got := TurboJSON(t.Context(), input)
		assert.NoError(t, got)
	})

	t.Run("No NPM Apps", func(t *testing.T) {
		t.Context()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:    "cms",
					Type:    appdef.AppTypePayload,
					Path:    "./apps/cms",
					UsesNPM: ptr.BoolPtr(false),
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		got := TurboJSON(t.Context(), input)
		assert.NoError(t, got)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "./apps/cms",
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := TurboJSON(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, "turbo.json")
		require.NoError(t, err)
		assert.NotEmpty(t, file)
		assert.Contains(t, string(file), "https://turborepo.com/schema.json")
		assert.NotContains(t, string(file), scaffold.WebKitNotice)
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "cms", Type: appdef.AppTypePayload, Path: "./apps/cms"},
			},
		}

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), appDef)

		got := TurboJSON(t.Context(), input)
		assert.Error(t, got)
	})
}
