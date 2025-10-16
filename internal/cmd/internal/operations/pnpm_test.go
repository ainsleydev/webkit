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
	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func TestCreateTurboJson(t *testing.T) {
	t.Parallel()

	t.Run("No Apps", func(t *testing.T) {
		t.Context()

		appDef := &appdef.Definition{
			Apps: []appdef.App{},
		}

		input := cmdtools.CommandInput{
			FS:          afero.NewMemMapFs(),
			AppDefCache: appDef,
		}

		got := CreateTurboJson(t.Context(), input)
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

		input := cmdtools.CommandInput{
			FS:          afero.NewMemMapFs(),
			AppDefCache: appDef,
		}

		got := CreateTurboJson(t.Context(), input)
		assert.NoError(t, got)
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "cms", Type: appdef.AppTypePayload, Path: "./apps/cms"},
			},
		}

		ctrl := gomock.NewController(t)
		fsMock := mocks.NewMockFS(ctrl)
		fsMock.EXPECT().
			MkdirAll(gomock.Any(), gomock.Any()).
			Return(fmt.Errorf("mkdir error"))

		input := cmdtools.CommandInput{
			FS:          fsMock,
			AppDefCache: appDef,
		}

		got := CreateTurboJson(t.Context(), input)
		assert.Error(t, got)
	})

	t.Run("Creates Successfully", func(t *testing.T) {
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

		fs := afero.NewMemMapFs()
		input := cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: appDef,
		}

		err := CreateTurboJson(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(fs, "turbo.json")
		require.NoError(t, err)
		assert.NotEmpty(t, file)
		assert.Contains(t, string(file), "https://turborepo.com/schema.json")
	})
}
