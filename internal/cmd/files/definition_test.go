package files

import (
	"context"
	"errors"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/state/manifest"
	"github.com/ainsleydev/webkit/internal/version"
)

func TestDefinition(t *testing.T) {
	t.Parallel()

	appDef := &appdef.Definition{
		WebkitVersion: "0.0.1", // Old version
		Project: appdef.Project{
			Name: "test-project",
			Repo: appdef.GitHubRepo{Owner: "test", Name: "test"},
		},
	}

	t.Run("Created", func(t *testing.T) {
		input := setup(t, afero.NewMemMapFs(), appDef)

		err := Definition(context.Background(), input)
		require.NoError(t, err)

		exists, err := afero.Exists(input.FS, appdef.JsonFileName)
		require.NoError(t, err)
		assert.True(t, exists)

		updated, err := appdef.Read(input.FS)
		require.NoError(t, err)

		t.Log("File Updated")
		{
			assert.Equal(t, version.Version, updated.WebkitVersion)
			assert.Equal(t, "test-project", updated.Project.Name)
		}
	})

	t.Run("Marshal Failure", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		tracker := manifest.NewTracker()

		orig := identMarshaller
		defer func() { identMarshaller = orig }()
		identMarshaller = func(_ any, _, _ string) ([]byte, error) {
			return nil, errors.New("marshal error")
		}

		input := setup(t, fs, appDef)
		input.Manifest = tracker

		err := Definition(context.Background(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "marshaling definition")
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), &appdef.Definition{})

		err := Definition(t.Context(), input)
		assert.Error(t, err)
	})
}
