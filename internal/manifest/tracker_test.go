package manifest

import (
	"errors"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTracker(t *testing.T) {
	t.Parallel()

	tracker := NewTracker()

	assert.NotNil(t, tracker)
	assert.NotNil(t, tracker.files)
	assert.Empty(t, tracker.files)
}

func TestTracker_Add(t *testing.T) {
	t.Parallel()

	t.Run("Single Entry", func(t *testing.T) {
		t.Parallel()

		tracker := NewTracker()
		entry := FileEntry{
			Path:   "app.json",
			Source: "project",
			Hash:   "abc123",
		}

		tracker.Add(entry)

		assert.Len(t, tracker.files, 1)
		assert.Equal(t, entry, tracker.files["app.json"])
	})

	t.Run("Multiple Entries", func(t *testing.T) {
		t.Parallel()

		tracker := NewTracker()
		entry1 := FileEntry{Path: "app.json", Source: "project", Hash: "abc123"}
		entry2 := FileEntry{Path: "config.yaml", Source: "app:web", Hash: "def456"}

		tracker.Add(entry1)
		tracker.Add(entry2)

		assert.Len(t, tracker.files, 2)
		assert.Equal(t, entry1, tracker.files["app.json"])
		assert.Equal(t, entry2, tracker.files["config.yaml"])
	})

	t.Run("Overwrites Existing Entry", func(t *testing.T) {
		t.Parallel()

		tracker := NewTracker()
		original := FileEntry{Path: "app.json", Source: "project", Hash: "old"}
		updated := FileEntry{Path: "app.json", Source: "project", Hash: "new"}

		tracker.Add(original)
		tracker.Add(updated)

		// Should only have one entry since it's the same path.
		assert.Len(t, tracker.files, 1)
		assert.Equal(t, "new", tracker.files["app.json"].Hash)
	})
}

func TestTracker_Save(t *testing.T) {
	t.Parallel()

	t.Run("Empty Tracker", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		tracker := NewTracker()

		err := tracker.Save(fs)
		require.NoError(t, err)

		exists, err := afero.Exists(fs, Path)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("With Entries", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		tracker := NewTracker()
		tracker.Add(FileEntry{
			Path:   "app.json",
			Source: "project",
			Hash:   "abc123",
		})

		err := tracker.Save(fs)
		require.NoError(t, err)

		manifest, err := Load(fs)
		require.NoError(t, err)

		assert.NotEmpty(t, manifest.Version)
		assert.False(t, manifest.GeneratedAt.IsZero())
		assert.Len(t, manifest.Files, 1)
		assert.Equal(t, "abc123", manifest.Files["app.json"].Hash)
	})

	t.Run("Marshal Error", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		tracker := NewTracker()

		tracker.marshaller = func(_ any, _, _ string) ([]byte, error) {
			return nil, errors.New("marshal error")
		}

		err := tracker.Save(fs)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "marshal error")
	})

	//t.Run("Creates Nested Directories", func(t *testing.T) {
	//	t.Parallel()
	//
	//	fs := afero.NewMemMapFs()
	//	tracker := NewTracker()
	//
	//	err := tracker.Save(fs)
	//	require.NoError(t, err)
	//
	//	exists, err := afero.Exists(fs)
	//	require.NoError(t, err)
	//	assert.True(t, exists)
	//})

	t.Run("Read-Only Filesystem", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewReadOnlyFs(afero.NewMemMapFs())
		tracker := NewTracker()

		err := tracker.Save(fs)
		assert.Error(t, err)
	})
}

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("Valid Manifest", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		tracker := NewTracker()
		tracker.Add(FileEntry{Path: "test.txt", Source: "project", Hash: "xyz"})

		err := tracker.Save(fs)
		require.NoError(t, err)

		manifest, err := Load(fs)
		require.NoError(t, err)

		assert.NotNil(t, manifest)
		assert.Len(t, manifest.Files, 1)
		assert.Equal(t, "xyz", manifest.Files["test.txt"].Hash)
	})

	t.Run("File Does Not Exist", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		_, err := Load(fs)
		assert.Error(t, err)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		err := afero.WriteFile(fs, Path, []byte("{invalid json"), 0644)
		require.NoError(t, err)

		_, err = Load(fs)
		assert.Error(t, err)
	})
}

func TestManifestTimestamp(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	tracker := NewTracker()

	before := time.Now()
	err := tracker.Save(fs)
	require.NoError(t, err)
	after := time.Now()

	manifest, err := Load(fs)
	require.NoError(t, err)

	// GeneratedAt should be between before and after.
	assert.True(t, manifest.GeneratedAt.After(before) || manifest.GeneratedAt.Equal(before))
	assert.True(t, manifest.GeneratedAt.Before(after) || manifest.GeneratedAt.Equal(after))
}
