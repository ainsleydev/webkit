package cicd

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestDriftDetection(t *testing.T) {
	t.Parallel()

	t.Run("Creates Workflow", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{})

		got := DriftDetection(t.Context(), input)
		assert.NoError(t, got)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "drift.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), &appdef.Definition{})

		got := DriftDetection(t.Context(), input)
		assert.Error(t, got)
	})
}
