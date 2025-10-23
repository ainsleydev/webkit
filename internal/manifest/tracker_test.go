package manifest

import (
	"errors"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/mocks"
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
		assert.Equal(t, ErrNoManifest, err)
	})

	t.Run("FS Error", func(t *testing.T) {
		t.Parallel()

		fs := mocks.NewMockFS(gomock.NewController(t))
		fs.EXPECT().Open(gomock.Any()).Return(nil, errors.New("open error"))

		_, err := Load(fs)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "open error")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		err := afero.WriteFile(fs, Path, []byte("{invalid json"), 0o644)
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

func TestTracker_WithPreviousManifest(t *testing.T) {
	t.Parallel()

	previousManifest := &Manifest{
		Version:     "v1.0.0",
		GeneratedAt: time.Now(),
		Files: map[string]FileEntry{
			"test.txt": {
				Path:        "test.txt",
				Hash:        "abc123",
				GeneratedAt: time.Now(),
			},
		},
	}

	tracker := NewTracker().WithPreviousManifest(previousManifest)

	assert.NotNil(t, tracker.previousManifest)
	assert.Equal(t, previousManifest, tracker.previousManifest)
}

func TestTracker_Add_PreservesTimestamp(t *testing.T) {
	t.Parallel()

	t.Run("Preserves timestamp when hash unchanged", func(t *testing.T) {
		t.Parallel()

		previousTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		previousManifest := &Manifest{
			Files: map[string]FileEntry{
				"test.txt": {
					Path:        "test.txt",
					Hash:        "abc123",
					GeneratedAt: previousTime,
				},
			},
		}

		tracker := NewTracker().WithPreviousManifest(previousManifest)
		newEntry := FileEntry{
			Path:        "test.txt",
			Hash:        "abc123", // Same hash
			GeneratedAt: time.Now(),
		}

		tracker.Add(newEntry)

		// Timestamp should be preserved from previous manifest
		assert.Equal(t, previousTime, tracker.files["test.txt"].GeneratedAt)
	})

	t.Run("Updates timestamp when hash changed", func(t *testing.T) {
		t.Parallel()

		previousTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		previousManifest := &Manifest{
			Files: map[string]FileEntry{
				"test.txt": {
					Path:        "test.txt",
					Hash:        "abc123",
					GeneratedAt: previousTime,
				},
			},
		}

		tracker := NewTracker().WithPreviousManifest(previousManifest)
		newTime := time.Now()
		newEntry := FileEntry{
			Path:        "test.txt",
			Hash:        "xyz789", // Different hash
			GeneratedAt: newTime,
		}

		tracker.Add(newEntry)

		// Timestamp should be the new one since hash changed
		assert.Equal(t, newTime.Truncate(time.Second), tracker.files["test.txt"].GeneratedAt.Truncate(time.Second))
	})

	t.Run("New file gets new timestamp", func(t *testing.T) {
		t.Parallel()

		previousManifest := &Manifest{
			Files: map[string]FileEntry{},
		}

		tracker := NewTracker().WithPreviousManifest(previousManifest)
		newTime := time.Now()
		newEntry := FileEntry{
			Path:        "new.txt",
			Hash:        "new123",
			GeneratedAt: newTime,
		}

		tracker.Add(newEntry)

		// New file should keep its new timestamp
		assert.Equal(t, newTime.Truncate(time.Second), tracker.files["new.txt"].GeneratedAt.Truncate(time.Second))
	})

	t.Run("Works without previous manifest", func(t *testing.T) {
		t.Parallel()

		tracker := NewTracker()
		newTime := time.Now()
		entry := FileEntry{
			Path:        "test.txt",
			Hash:        "abc123",
			GeneratedAt: newTime,
		}

		tracker.Add(entry)

		// Should just add the entry normally
		assert.Equal(t, newTime.Truncate(time.Second), tracker.files["test.txt"].GeneratedAt.Truncate(time.Second))
	})
}

func TestTracker_Save_PreservesManifestTimestamp(t *testing.T) {
	t.Parallel()

	t.Run("Preserves manifest timestamp when nothing changed", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		previousTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		fileTime := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)

		previousManifest := &Manifest{
			Version:     "v1.0.0",
			GeneratedAt: previousTime,
			Files: map[string]FileEntry{
				"test.txt": {
					Path:        "test.txt",
					Hash:        "abc123",
					GeneratedAt: fileTime,
				},
			},
		}

		tracker := NewTracker().WithPreviousManifest(previousManifest)
		tracker.Add(FileEntry{
			Path:        "test.txt",
			Hash:        "abc123", // Same hash, so timestamp will be preserved
			GeneratedAt: time.Now(),
		})

		err := tracker.Save(fs)
		require.NoError(t, err)

		manifest, err := Load(fs)
		require.NoError(t, err)

		// Manifest timestamp should be preserved
		assert.Equal(t, previousTime, manifest.GeneratedAt)
	})

	t.Run("Updates manifest timestamp when file changed", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		previousTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		previousManifest := &Manifest{
			Version:     "v1.0.0",
			GeneratedAt: previousTime,
			Files: map[string]FileEntry{
				"test.txt": {
					Path:        "test.txt",
					Hash:        "abc123",
					GeneratedAt: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
				},
			},
		}

		tracker := NewTracker().WithPreviousManifest(previousManifest)
		tracker.Add(FileEntry{
			Path:        "test.txt",
			Hash:        "xyz789", // Different hash
			GeneratedAt: time.Now(),
		})

		before := time.Now()
		err := tracker.Save(fs)
		require.NoError(t, err)
		after := time.Now()

		manifest, err := Load(fs)
		require.NoError(t, err)

		// Manifest timestamp should be updated
		assert.True(t, manifest.GeneratedAt.After(previousTime))
		assert.True(t, manifest.GeneratedAt.After(before) || manifest.GeneratedAt.Equal(before))
		assert.True(t, manifest.GeneratedAt.Before(after) || manifest.GeneratedAt.Equal(after))
	})

	t.Run("Updates manifest timestamp when file added", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		previousTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		previousManifest := &Manifest{
			Version:     "v1.0.0",
			GeneratedAt: previousTime,
			Files: map[string]FileEntry{
				"old.txt": {
					Path:        "old.txt",
					Hash:        "old123",
					GeneratedAt: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
				},
			},
		}

		tracker := NewTracker().WithPreviousManifest(previousManifest)
		tracker.Add(FileEntry{
			Path:        "old.txt",
			Hash:        "old123",
			GeneratedAt: time.Now(),
		})
		tracker.Add(FileEntry{
			Path:        "new.txt", // New file
			Hash:        "new123",
			GeneratedAt: time.Now(),
		})

		before := time.Now()
		err := tracker.Save(fs)
		require.NoError(t, err)
		after := time.Now()

		manifest, err := Load(fs)
		require.NoError(t, err)

		// Manifest timestamp should be updated because a new file was added
		assert.True(t, manifest.GeneratedAt.After(previousTime))
		assert.True(t, manifest.GeneratedAt.After(before) || manifest.GeneratedAt.Equal(before))
		assert.True(t, manifest.GeneratedAt.Before(after) || manifest.GeneratedAt.Equal(after))
	})

	t.Run("Updates manifest timestamp when file removed", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		previousTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		previousManifest := &Manifest{
			Version:     "v1.0.0",
			GeneratedAt: previousTime,
			Files: map[string]FileEntry{
				"file1.txt": {
					Path:        "file1.txt",
					Hash:        "hash1",
					GeneratedAt: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
				},
				"file2.txt": {
					Path:        "file2.txt",
					Hash:        "hash2",
					GeneratedAt: time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
				},
			},
		}

		tracker := NewTracker().WithPreviousManifest(previousManifest)
		// Only add one file (removing file2.txt)
		tracker.Add(FileEntry{
			Path:        "file1.txt",
			Hash:        "hash1",
			GeneratedAt: time.Now(),
		})

		before := time.Now()
		err := tracker.Save(fs)
		require.NoError(t, err)
		after := time.Now()

		manifest, err := Load(fs)
		require.NoError(t, err)

		// Manifest timestamp should be updated because a file was removed
		assert.True(t, manifest.GeneratedAt.After(previousTime))
		assert.True(t, manifest.GeneratedAt.After(before) || manifest.GeneratedAt.Equal(before))
		assert.True(t, manifest.GeneratedAt.Before(after) || manifest.GeneratedAt.Equal(after))
	})
}
