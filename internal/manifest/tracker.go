package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/version"
)

type Tracker struct {
	files map[string]FileEntry
}

func NewTracker() *Tracker {
	return &Tracker{
		files: make(map[string]FileEntry),
	}
}

func (t *Tracker) Add(entry FileEntry) {
	t.files[entry.Path] = entry
}

func (t *Tracker) Save(fs afero.Fs, path string) error {
	manifest := Manifest{
		Version:     version.Version,
		GeneratedAt: time.Now(),
		Files:       t.files,
	}

	data, err := json.MarshalIndent(manifest, "", "\t")
	if err != nil {
		return err
	}

	err = fs.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	return afero.WriteFile(fs, path, data, 0644)
}

func Load(fs afero.Fs, path string) (*Manifest, error) {
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, err
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func hashContent(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
