package operations

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/mocks"
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

	t.Run("Creates Successfully", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		err := CreateGitSettings(t.Context(), cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: appDef,
		})
		assert.NoError(t, err)

		for path := range gitSettingsTemplates {
			file, err := afero.ReadFile(fs, path)
			assert.NoError(t, err)
			assert.NotEmpty(t, file)
		}

		// Check settings.yaml exists
		settingsFile, err := afero.ReadFile(fs, ".github/settings.yaml")
		assert.NoError(t, err)
		assert.NotEmpty(t, settingsFile)
		assert.NoError(t, testutil.ValidateYAML(t, settingsFile))
	})

	t.Run("Validates dependabot.yaml schema", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := CreateGitSettings(t.Context(), cmdtools.CommandInput{
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

		ctrl := gomock.NewController(t)
		fsMock := mocks.NewMockFS(ctrl)
		fsMock.EXPECT().
			MkdirAll(gomock.Any(), gomock.Any()).
			Return(fmt.Errorf("mkdir error"))

		input := cmdtools.CommandInput{
			FS:          fsMock,
			AppDefCache: &appdef.Definition{},
		}

		got := CreateGitSettings(t.Context(), input)
		assert.Error(t, got)
	})
}
