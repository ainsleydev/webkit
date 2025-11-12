package files

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestSchema(t *testing.T) {
	t.Parallel()

	t.Run("Creates Schema", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{})

		err := Schema(t.Context(), input)
		require.NoError(t, err)

		t.Log("File Exists")
		{
			exists, err := afero.Exists(input.FS, ".webkit/schema.json")
			require.NoError(t, err)
			assert.True(t, exists, "schema file should be created")
		}
	})
}
