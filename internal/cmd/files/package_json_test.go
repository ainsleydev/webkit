package files

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/util/testutil"
)

func TestPackageJSON(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		appDef := &appdef.Definition{
			Project: appdef.Project{Name: "my-website"},
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "cms",
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := PackageJSON(t.Context(), input)
		assert.NoError(t, err)

		t.Log("File Exists")
		{
			exists, err := afero.Exists(input.FS, "package.json")
			assert.NoError(t, err)
			assert.True(t, exists)
		}

		t.Log("Conforms to Schema")
		{
			schema, err := testutil.SchemaFromURL(t, "https://www.schemastore.org/package.json")
			require.NoError(t, err)

			file, err := afero.ReadFile(input.FS, "package.json")
			require.NoError(t, err)

			err = schema.ValidateJSON(file)
			assert.NoError(t, err, "Package.json file conforms to schema")
		}
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), &appdef.Definition{})

		got := PackageJSON(t.Context(), input)
		assert.Error(t, got)
	})
}
