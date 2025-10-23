package cmdtools

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/manifest"
)

// Note: Testing the error path (missing app.json) is not possible because AppDef()
// calls os.Exit(1) directly. This would require a subprocess testing pattern which
// is beyond the scope of these unit tests.
func TestCommandInput_AppDef(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "app.json", []byte(`{
			"project": {
				"name": "test-app",
				"repo": {
					"owner": "test",
					"name": "repo"
				}
			}
		}`), 0o644)
		require.NoError(t, err)

		input := CommandInput{
			FS:       fs,
			Manifest: manifest.NewTracker(),
		}

		def := input.AppDef()
		assert.NotNil(t, def)
		assert.Equal(t, "test-app", def.Project.Name)
		assert.Equal(t, "test", def.Project.Repo.Owner)
	})

	t.Run("Caching", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "app.json", []byte(`{
			"project": {
				"name": "cached-app",
				"repo": {
					"owner": "cached",
					"name": "repo"
				}
			}
		}`), 0o644)
		require.NoError(t, err)

		input := CommandInput{
			FS:       fs,
			Manifest: manifest.NewTracker(),
		}

		def1 := input.AppDef()
		require.NotNil(t, def1)

		err = afero.WriteFile(fs, "app.json", []byte(`{
			"project": {
				"name": "modified-app",
				"repo": {
					"owner": "modified",
					"name": "repo"
				}
			}
		}`), 0o644)
		require.NoError(t, err)

		def2 := input.AppDef()
		assert.Same(t, def1, def2)
		assert.Equal(t, "cached-app", def2.Project.Name)
	})
}
