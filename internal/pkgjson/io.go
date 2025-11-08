package pkgjson

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Exists checks if a package.json file exists at the given path.
func Exists(fs afero.Fs, path string) bool {
	exists, _ := afero.Exists(fs, path)
	return exists
}

// Read parses package.json from disk and preserves unknown fields in raw.
func Read(fs afero.Fs, path string) (*PackageJSON, error) {
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, errors.Wrap(err, "reading package.json")
	}

	pkg := &PackageJSON{
		Dependencies:     make(map[string]string),
		DevDependencies:  make(map[string]string),
		PeerDependencies: make(map[string]string),
		raw:              make(map[string]any),
	}

	// Unmarshal into struct (known fields)
	if err = json.Unmarshal(data, pkg); err != nil {
		return nil, errors.Wrap(err, "parsing package.json struct")
	}

	// Unmarshal into map to capture all fields (including unknown)
	if err = json.Unmarshal(data, &pkg.raw); err != nil {
		return nil, errors.Wrap(err, "parsing package.json raw")
	}

	return pkg, nil
}

// Write saves PackageJSON to disk with npm-standard field order,
// pretty-printed JSON & no HTML escaping.
func Write(fs afero.Fs, path string, pkg *PackageJSON) error {
	// Encode JSON with indentation and no HTML escaping
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(pkg); err != nil {
		return errors.Wrap(err, "encoding package.json")
	}

	// Write file
	if err := afero.WriteFile(fs, path, buf.Bytes(), 0o644); err != nil {
		return errors.Wrap(err, "writing package.json")
	}

	return nil
}
