package files

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/scaffold"
)

func TestSchema(t *testing.T) {
	t.Parallel()

	t.Run("Creates Schema", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{})

		err := Schema(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, ".webkit/schema.json")
		assert.NoError(t, err)
		assert.NotEmpty(t, file)
		assert.NotContains(t, string(file), scaffold.WebKitNotice)
	})
}
