package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/manifest"
)

func TestDrift(t *testing.T) {
	t.Run("No Drift, No Manifest", func(t *testing.T) {
		input, buf := setupWithPrinter(t, afero.NewMemMapFs(), &appdef.Definition{})

		err := drift(t.Context(), input)
		assert.NoError(t, err)
		assert.Empty(t, buf.String()) // No output when no manifest exists
	})

	t.Run("No Drift, Empty Manifest", func(t *testing.T) {
		input, buf := setupWithPrinter(t, afero.NewMemMapFs(), &appdef.Definition{})

		tracker := manifest.NewTracker()
		err := tracker.Save(input.FS)
		require.NoError(t, err)

		err = drift(t.Context(), input)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "No drift detected")
	})

	t.Run("No Drift, Matching Files", func(t *testing.T) {
		input, buf := setupWithPrinter(t, afero.NewMemMapFs(), &appdef.Definition{})

		tracker := manifest.NewTracker()
		tracker.Add(manifest.FileEntry{
			Path:   ".github/workflows/deploy.yml",
			Source: "template",
			Hash:   manifest.HashContent([]byte("content")),
		})
		err := tracker.Save(input.FS)
		require.NoError(t, err)

		err = afero.WriteFile(input.FS, ".github/workflows/deploy.yml", []byte("content"), 0644)
		require.NoError(t, err)

		err = drift(t.Context(), input)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "No drift detected")
	})

	t.Run("Drift, File Changed", func(t *testing.T) {
		input, buf := setupWithPrinter(t, afero.NewMemMapFs(), &appdef.Definition{})

		tracker := manifest.NewTracker()
		tracker.Add(manifest.FileEntry{
			Path:   ".github/workflows/deploy.yml",
			Hash:   "abc123",
			Source: "template",
		})
		err := tracker.Save(input.FS)
		require.NoError(t, err)

		err = afero.WriteFile(input.FS, ".github/workflows/deploy.yml", []byte("modified content"), 0644)
		require.NoError(t, err)

		err = drift(t.Context(), input)
		assert.Error(t, err)
		assert.Contains(t, buf.String(), "Drift found")
		assert.Contains(t, buf.String(), "Action Required: Run webkit update")
		assert.Contains(t, buf.String(), ".github/workflows/deploy.yml")
	})

	t.Run("Drift - File Deleted", func(t *testing.T) {
		input, buf := setupWithPrinter(t, afero.NewMemMapFs(), &appdef.Definition{})

		tracker := manifest.NewTracker()
		tracker.Add(manifest.FileEntry{
			Path:   ".github/workflows/deploy.yml",
			Hash:   "abc123",
			Source: "template",
		})
		err := tracker.Save(input.FS)
		require.NoError(t, err)

		// Don't create the file, it's been deleted.

		err = drift(t.Context(), input)
		assert.Error(t, err)
		assert.Contains(t, buf.String(), "Drift found")
		assert.Contains(t, buf.String(), ".github/workflows/deploy.yml")
	})
}
