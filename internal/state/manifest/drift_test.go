package manifest

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDriftReason_String(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input DriftReason
		want  string
	}{
		"Modified": {DriftReasonModified, "modified"},
		"Deleted":  {DriftReasonDeleted, "deleted"},
		"Outdated": {DriftReasonOutdated, "outdated"},
		"New":      {DriftReasonNew, "new"},
		"Unknown":  {DriftReason(999), "unknown(999)"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.input.String()
			assert.Equal(t, test.want, got)
			assert.IsType(t, "", got)
		})
	}
}

func TestDriftReason_FilterEntries(t *testing.T) {
	t.Parallel()

	entries := []DriftEntry{
		{Path: "a.txt", Type: DriftReasonModified, Source: "src1", Generator: "genA"},
		{Path: "b.txt", Type: DriftReasonDeleted, Source: "src2", Generator: "genB"},
		{Path: "c.txt", Type: DriftReasonOutdated, Source: "src3", Generator: "genC"},
		{Path: "d.txt", Type: DriftReasonNew, Source: "src4", Generator: "genD"},
		{Path: "e.txt", Type: DriftReasonModified, Source: "src5", Generator: "genE"},
	}

	tt := map[string]struct {
		input DriftReason
		want  []string
	}{
		"Modified": {DriftReasonModified, []string{"a.txt", "e.txt"}},
		"Deleted":  {DriftReasonDeleted, []string{"b.txt"}},
		"Outdated": {DriftReasonOutdated, []string{"c.txt"}},
		"New":      {DriftReasonNew, []string{"d.txt"}},
		"Unknown":  {DriftReason(999), []string{}},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.input.FilterEntries(entries)

			var gotPaths []string
			for _, e := range got {
				gotPaths = append(gotPaths, e.Path)
			}

			assert.ElementsMatch(t, test.want, gotPaths)
		})
	}
}

func TestDetectDrift(t *testing.T) {
	t.Parallel()

	t.Run("No Drift", func(t *testing.T) {
		t.Parallel()

		actualFS := afero.NewMemMapFs()
		expectedFS := afero.NewMemMapFs()

		content := []byte("test content")

		// Both have same file
		require.NoError(t, afero.WriteFile(actualFS, "file.txt", content, 0o644))
		require.NoError(t, afero.WriteFile(expectedFS, "file.txt", content, 0o644))

		// Both manifests match
		setupManifest(t, actualFS, "file.txt", content)
		setupManifest(t, expectedFS, "file.txt", content)

		drift, err := DetectDrift(actualFS, expectedFS)
		require.NoError(t, err)
		assert.Empty(t, drift)
	})

	t.Run("Error - No Actual Manifest", func(t *testing.T) {
		t.Parallel()

		actualFS := afero.NewMemMapFs()
		expectedFS := afero.NewMemMapFs()

		// No manifest on actual filesystem
		// Setup expected manifest
		content := []byte("test")
		require.NoError(t, afero.WriteFile(expectedFS, "file.txt", content, 0o644))
		setupManifest(t, expectedFS, "file.txt", content)

		drift, err := DetectDrift(actualFS, expectedFS)

		assert.Error(t, err, "should error when actual manifest missing")
		assert.Nil(t, drift)
	})

	t.Run("Error - No Expected Manifest", func(t *testing.T) {
		t.Parallel()

		actualFS := afero.NewMemMapFs()
		expectedFS := afero.NewMemMapFs()

		content := []byte("test")
		require.NoError(t, afero.WriteFile(actualFS, "file.txt", content, 0o644))
		setupManifest(t, actualFS, "file.txt", content)

		drift, err := DetectDrift(actualFS, expectedFS)
		assert.Error(t, err, "should error when expected manifest missing")
		assert.Nil(t, drift)
	})

	t.Run("FS Error (ReadFile)", func(t *testing.T) {
		t.Parallel()

		actualFS := afero.NewMemMapFs()
		expectedFS := afero.NewMemMapFs()

		content := []byte("test")

		// File exists on actual
		require.NoError(t, afero.WriteFile(actualFS, "file.txt", content, 0o644))

		// Manifest says file should exist on expected, but it doesn't
		setupManifest(t, actualFS, "file.txt", content)
		setupManifest(t, expectedFS, "file.txt", content)

		drift, err := DetectDrift(actualFS, expectedFS)
		require.NoError(t, err)
		assert.Empty(t, drift, "should skip files that can't be read from expected")
	})

	t.Run("Manual Modification", func(t *testing.T) {
		t.Parallel()

		actualFS := afero.NewMemMapFs()
		expectedFS := afero.NewMemMapFs()

		originalContent := []byte("original")
		modifiedContent := []byte("user modified")

		// User manually modified the file
		require.NoError(t, afero.WriteFile(actualFS, "file.txt", modifiedContent, 0o644))

		// Expected hasn't changed
		require.NoError(t, afero.WriteFile(expectedFS, "file.txt", originalContent, 0o644))

		// Old manifest shows original
		setupManifest(t, actualFS, "file.txt", originalContent)
		setupManifest(t, expectedFS, "file.txt", originalContent)

		drift, err := DetectDrift(actualFS, expectedFS)

		require.NoError(t, err)
		require.Len(t, drift, 1)
		assert.Equal(t, DriftReasonModified, drift[0].Type)
		assert.Equal(t, "file.txt", drift[0].Path)
	})

	t.Run("Outdated From app.json Change", func(t *testing.T) {
		t.Parallel()

		actualFS := afero.NewMemMapFs()
		expectedFS := afero.NewMemMapFs()

		oldContent := []byte("DATABASE_URL=old")
		newContent := []byte("DATABASE_URL=new")

		// File on disk matches what was last generated
		require.NoError(t, afero.WriteFile(actualFS, ".env", oldContent, 0o644))

		// But app.json changed, so expected is different
		require.NoError(t, afero.WriteFile(expectedFS, ".env", newContent, 0o644))

		// Old manifest shows we generated the old content
		setupManifest(t, actualFS, ".env", oldContent)
		setupManifest(t, expectedFS, ".env", newContent)

		drift, err := DetectDrift(actualFS, expectedFS)

		require.NoError(t, err)
		require.Len(t, drift, 1)
		assert.Equal(t, DriftReasonOutdated, drift[0].Type)
		assert.Equal(t, ".env", drift[0].Path)
	})

	t.Run("Missing File", func(t *testing.T) {
		t.Parallel()

		actualFS := afero.NewMemMapFs()
		expectedFS := afero.NewMemMapFs()

		content := []byte("should exist")

		// Expected has file
		require.NoError(t, afero.WriteFile(expectedFS, "new.txt", content, 0o644))

		// Actual doesn't
		setupManifest(t, actualFS)
		setupManifest(t, expectedFS, "new.txt", content)

		drift, err := DetectDrift(actualFS, expectedFS)

		require.NoError(t, err)
		require.Len(t, drift, 1)
		assert.Equal(t, DriftReasonNew, drift[0].Type)
		assert.Equal(t, "new.txt", drift[0].Path)
	})

	t.Run("Orphaned File", func(t *testing.T) {
		t.Parallel()

		actualFS := afero.NewMemMapFs()
		expectedFS := afero.NewMemMapFs()

		oldContent := []byte("old")

		// File exists on actual
		require.NoError(t, afero.WriteFile(actualFS, "orphaned.txt", oldContent, 0o644))

		// Old manifest tracked it
		setupManifest(t, actualFS, "orphaned.txt", oldContent)

		// New manifest doesn't (removed from app.json)
		setupManifest(t, expectedFS)

		drift, err := DetectDrift(actualFS, expectedFS)

		require.NoError(t, err)
		require.Len(t, drift, 1)
		assert.Equal(t, DriftReasonDeleted, drift[0].Type)
		assert.Equal(t, "orphaned.txt", drift[0].Path)
	})

	t.Run("Skips Scaffold Files", func(t *testing.T) {
		t.Parallel()

		actualFS := afero.NewMemMapFs()
		expectedFS := afero.NewMemMapFs()

		// User modified scaffold file
		require.NoError(t, afero.WriteFile(actualFS, ".env", []byte("USER=modified"), 0o644))
		require.NoError(t, afero.WriteFile(expectedFS, ".env", []byte("ORIGINAL=true"), 0o644))

		setupManifestWithScaffold(t, actualFS, ".env", []byte("ORIGINAL=true"), true)
		setupManifestWithScaffold(t, expectedFS, ".env", []byte("ORIGINAL=true"), true)

		drift, err := DetectDrift(actualFS, expectedFS)

		require.NoError(t, err)
		assert.Empty(t, drift, "scaffold files should be ignored")
	})
}

func setupManifest(t *testing.T, fs afero.Fs, pathsAndContent ...interface{}) {
	t.Helper()

	tracker := NewTracker()

	for i := 0; i < len(pathsAndContent); i += 2 {
		path := pathsAndContent[i].(string)
		content := pathsAndContent[i+1].([]byte)

		tracker.Add(FileEntry{
			Path:         path,
			Source:       "test",
			Generator:    "test:gen",
			ScaffoldMode: false,
			Hash:         HashContent(content),
		})
	}

	require.NoError(t, tracker.Save(fs))
}

func setupManifestWithScaffold(t *testing.T, fs afero.Fs, path string, content []byte, scaffold bool) {
	t.Helper()

	tracker := NewTracker()
	tracker.Add(FileEntry{
		Path:         path,
		Source:       "test",
		ScaffoldMode: scaffold,
		Hash:         HashContent(content),
	})

	require.NoError(t, tracker.Save(fs))
}
