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
			"name": "test-app",
			"repo": "test/repo",
			"email": "test@example.com",
			"types": ["web"]
		}`), 0644)
		require.NoError(t, err)

		input := CommandInput{
			FS:       fs,
			Manifest: manifest.NewTracker(),
		}

		def := input.AppDef()
		assert.NotNil(t, def)
		assert.Equal(t, "test-app", def.Name)
		assert.Equal(t, "test/repo", def.Repo)
	})

	t.Run("Caching", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "app.json", []byte(`{
			"name": "cached-app",
			"repo": "cached/repo",
			"email": "cached@example.com",
			"types": ["web"]
		}`), 0644)
		require.NoError(t, err)

		input := CommandInput{
			FS:       fs,
			Manifest: manifest.NewTracker(),
		}

		def1 := input.AppDef()
		require.NotNil(t, def1)

		err = afero.WriteFile(fs, "app.json", []byte(`{
			"name": "modified-app",
			"repo": "modified/repo",
			"email": "modified@example.com",
			"types": ["api"]
		}`), 0644)
		require.NoError(t, err)

		def2 := input.AppDef()
		assert.Same(t, def1, def2)
		assert.Equal(t, "cached-app", def2.Name)
	})
}
