package files

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/util/testutil"
)

func TestCreatePackageJson(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		fs := afero.NewMemMapFs()

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

		err := PackageJSON(t.Context(), cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: appDef,
		})
		assert.NoError(t, err)

		t.Log("File Exists")
		{
			exists, err := afero.Exists(fs, "package.json")
			assert.NoError(t, err)
			assert.True(t, exists)
		}

		t.Log("Conforms to Schema")
		{
			schema, err := testutil.SchemaFromURL(t, "https://www.schemastore.org/package.json")
			require.NoError(t, err)

			file, err := afero.ReadFile(fs, "package.json")
			require.NoError(t, err)

			err = schema.ValidateJSON(file)
			assert.NoError(t, err, "Package.json file conforms to schema")
		}
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input := cmdtools.CommandInput{
			FS:          afero.NewReadOnlyFs(afero.NewMemMapFs()),
			AppDefCache: &appdef.Definition{},
		}

		got := PackageJSON(t.Context(), input)
		assert.Error(t, got)
	})
}
