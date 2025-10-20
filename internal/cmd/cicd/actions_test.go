package cicd

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestActionTemplates(t *testing.T) {
	t.Parallel()

	t.Run("Creates Templates", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{})

		got := ActionTemplates(t.Context(), input)
		assert.NoError(t, got)

		for _, path := range actionTemplates {
			_, err := afero.ReadFile(input.FS, filepath.Join(actionsPath, path))
			require.NoError(t, err)
		}
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), &appdef.Definition{})

		got := ActionTemplates(t.Context(), input)
		assert.Error(t, got)
	})
}
