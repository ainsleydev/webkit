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
	files      map[string]FileEntry
	mtx        *sync.Mutex
	marshaller func(v any, prefix, indent string) ([]byte, error)
}

// NewTracker creates a tracker with an initialized file map.
func NewTracker() *Tracker {
	return &Tracker{
		files:      make(map[string]FileEntry),
		mtx:        &sync.Mutex{},
		marshaller: json.MarshalIndent,
	}
}

// Path defines the filepath where the manifest resides.
var Path = filepath.Join(".webkit", "manifest.json")

// Add stores a file entry in the tracker, keyed by its path.
// If an entry with the same path exists, it will be overwritten.
func (t *Tracker) Add(entry FileEntry) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.files[entry.Path] = entry
}

// Save writes the tracker's files to a manifest JSON file.
// Creates parent directories if they don't exist.
func (t *Tracker) Save(fs afero.Fs) error {
	manifest := Manifest{
		Version:     version.Version,
		GeneratedAt: time.Now(),
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
