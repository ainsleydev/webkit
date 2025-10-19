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

func TestCreateGitSettings(t *testing.T) {
	t.Parallel()

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

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		err := GitSettings(t.Context(), cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: appDef,
		})
		assert.NoError(t, err)

		for path := range gitSettingsTemplates {
			file, err := afero.ReadFile(fs, path)
			assert.NoError(t, err)
			assert.NotEmpty(t, file)
		}

		got, err := afero.ReadFile(fs, ".github/settings.yaml")
		assert.NoError(t, err)
		assert.NotEmpty(t, got)
		assert.NoError(t, testutil.ValidateYAML(t, got))
	})

	t.Run("Validates dependabot.yaml schema", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := GitSettings(t.Context(), cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: appDef,
		})
		assert.NoError(t, err)

		schema, err := testutil.SchemaFromURL(t, "https://www.schemastore.org/dependabot-2.0.json")
		require.NoError(t, err)

		dep, err := afero.ReadFile(fs, ".github/dependabot.yaml")
		require.NoError(t, err)

		err = schema.ValidateYAML(dep)
		assert.NoError(t, err, "Dependabot file conforms to schema")
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input := cmdtools.CommandInput{
			FS:          afero.NewReadOnlyFs(afero.NewMemMapFs()),
			AppDefCache: &appdef.Definition{},
		}

		got := GitSettings(t.Context(), input)
		assert.Error(t, got)
	})
}
