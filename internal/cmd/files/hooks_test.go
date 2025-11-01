package files

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestHooks(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
		})

		got := Hooks(t.Context(), input)
		assert.NoError(t, got)

		file, err := afero.ReadFile(input.FS, "lefthook.yml")
		assert.NoError(t, err)
		assert.NotEmpty(t, file)
		assert.Contains(t, string(file), "test-project")
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), &appdef.Definition{})

		got := Hooks(t.Context(), input)
		assert.Error(t, got)
	})
}
