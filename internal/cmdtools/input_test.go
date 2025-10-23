package cmdtools

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/manifest"
)

func TestCommandInput_AppDef_Success(t *testing.T) {
	t.Parallel()

	// Create in-memory filesystem with a valid app.json
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

	// First call should read and cache the definition
	def1 := input.AppDef()
	assert.NotNil(t, def1)
	assert.Equal(t, "test-app", def1.Name)
	assert.Equal(t, "test/repo", def1.Repo)

	// Second call should return cached value
	def2 := input.AppDef()
	assert.Same(t, def1, def2, "AppDef should return cached instance")
}

func TestCommandInput_AppDef_Caching(t *testing.T) {
	t.Parallel()

	// Create in-memory filesystem with a valid app.json
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

	// First call
	def1 := input.AppDef()
	require.NotNil(t, def1)

	// Modify the file on disk
	err = afero.WriteFile(fs, "app.json", []byte(`{
		"name": "modified-app",
		"repo": "modified/repo",
		"email": "modified@example.com",
		"types": ["api"]
	}`), 0644)
	require.NoError(t, err)

	// Second call should still return the cached value, not the modified one
	def2 := input.AppDef()
	assert.Equal(t, "cached-app", def2.Name, "AppDef should use cached value")
	assert.Equal(t, "cached/repo", def2.Repo, "AppDef should use cached value")
}

// Note: Testing the error path (missing app.json) is complex because AppDef()
// calls os.Exit(1) directly. This would require a subprocess testing pattern
// which is beyond the scope of these unit tests. The error behavior should be
// verified through integration tests or manual testing.
