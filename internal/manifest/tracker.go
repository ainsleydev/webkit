package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/version"
)

// Tracker maintains a collection of generated files and their metadata.
// Used to track which files were created by webkit and from which source.
type Tracker struct {
	files            map[string]FileEntry
	previousManifest *Manifest
	mtx              *sync.Mutex
	marshaller       func(v any, prefix, indent string) ([]byte, error)
}

// NewTracker creates a tracker with an initialized file map.
// If a previous manifest is provided, it will be used to preserve timestamps
// for files that haven't changed.
func NewTracker() *Tracker {
	return &Tracker{
		files:      make(map[string]FileEntry),
		mtx:        &sync.Mutex{},
		marshaller: json.MarshalIndent,
	}
}

// WithPreviousManifest sets the previous manifest for timestamp preservation.
func (t *Tracker) WithPreviousManifest(previous *Manifest) *Tracker {
	t.previousManifest = previous
	return t
}

// Path defines the filepath where the manifest resides.
var Path = filepath.Join(".webkit", "manifest.json")

// Add stores a file entry in the tracker, keyed by its path.
// If an entry with the same path exists, it will be overwritten.
// If the previous manifest contains the same file with the same hash,
// the GeneratedAt timestamp will be preserved.
func (t *Tracker) Add(entry FileEntry) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	// Check if we have a previous manifest and if this file existed before
	if t.previousManifest != nil {
		if previousEntry, exists := t.previousManifest.Files[entry.Path]; exists {
			// If the content hash hasn't changed, preserve the previous timestamp
			if previousEntry.Hash == entry.Hash {
				entry.GeneratedAt = previousEntry.GeneratedAt
			}
		}
	}

	t.files[entry.Path] = entry
}

// Save writes the tracker's files to a manifest JSON file.
// Creates parent directories if they don't exist.
// If a previous manifest was provided and no files have changed,
// the manifest's GeneratedAt timestamp will be preserved.
func (t *Tracker) Save(fs afero.Fs) error {
	generatedAt := time.Now()

	// If we have a previous manifest, check if any files actually changed
	if t.previousManifest != nil && t.hasFilesChanged() == false {
		// No changes detected, preserve the previous manifest's timestamp
		generatedAt = t.previousManifest.GeneratedAt
	}

	manifest := Manifest{
		Version:     version.Version,
		GeneratedAt: generatedAt,
		Files:       t.files,
	}

	data, err := t.marshaller(manifest, "", "\t")
	if err != nil {
		return err
	}

	err = fs.MkdirAll(filepath.Dir(Path), os.ModePerm)
	if err != nil {
		return err
	}

	return afero.WriteFile(fs, Path, data, 0o644)
}

// hasFilesChanged checks if any files have different timestamps
// compared to the previous manifest, indicating actual changes.
func (t *Tracker) hasFilesChanged() bool {
	if t.previousManifest == nil {
		return true
	}

	// Check if file count differs
	if len(t.files) != len(t.previousManifest.Files) {
		return true
	}

	// Check if any file has a different GeneratedAt timestamp
	// (which would have been updated by Add if the hash changed)
	for path, newEntry := range t.files {
		previousEntry, exists := t.previousManifest.Files[path]
		if !exists {
			return true
		}
		// If the timestamp was preserved, they'll be equal
		// If it was updated, they'll be different
		if !newEntry.GeneratedAt.Equal(previousEntry.GeneratedAt) {
			return true
		}
	}

	return false
}

// ErrNoManifest is returned by Load() when there hasen't been
// a manifest generated yet.
var ErrNoManifest = fmt.Errorf("no manifest found")

// Load reads a manifest JSON file and deserializes it.
// Returns an error if the file doesn't exist or contains invalid JSON.
func Load(fs afero.Fs) (*Manifest, error) {
	data, err := afero.ReadFile(fs, Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoManifest
		}
		return nil, err
	}

	var m Manifest
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}
