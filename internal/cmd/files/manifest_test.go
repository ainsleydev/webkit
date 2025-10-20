package files

import (
	"errors"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/mocks"
)

func TestManifest(t *testing.T) {
	t.Parallel()

	t.Run("Exists Error", func(t *testing.T) {
		t.Parallel()

		mock := mocks.NewMockFS(gomock.NewController(t))
		mock.EXPECT().Stat(gomock.Any()).Return(nil, errors.New("stat error"))

		input := setup(t, mock, &appdef.Definition{})

		err := Manifest(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "stat error")
	})

	t.Run("Already Exists", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{})

		err := manifest.NewTracker().Save(input.FS)
		require.NoError(t, err)

		err = Manifest(t.Context(), input)
		assert.NoError(t, err)
	})

	t.Run("Save Error", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), &appdef.Definition{})

		err := Manifest(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "operation not permitted")
	})

	t.Run("Creates Manifest", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{})

		err := Manifest(t.Context(), input)
		require.NoError(t, err)

		t.Log("File Exists")
		{
			exists, err := afero.Exists(input.FS, manifest.Path)
			require.NoError(t, err)
			assert.True(t, exists, "manifest file should be created")
		}

		t.Log("Verify Manifest")
		{
			loadedManifest, err := manifest.Load(input.FS)
			require.NoError(t, err)
			assert.NotNil(t, loadedManifest)
			assert.NotEmpty(t, loadedManifest.Version)
			assert.False(t, loadedManifest.GeneratedAt.IsZero())
		}
	})
}
