package pkgjson

import (
	"bytes"
	"encoding/json"

	"github.com/perimeterx/marshmallow"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Exists checks if a package.json file exists at the given path.
func Exists(fs afero.Fs, path string) bool {
	exists, err := afero.Exists(fs, path)
	if err != nil {
		return false
	}
	return exists
}

// Read reads and parses a package.json file from the filesystem.
// Uses marshmallow to preserve all fields including unknown ones.
func Read(fs afero.Fs, path string) (*PackageJSON, error) {
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, errors.Wrap(err, "reading package.json")
	}

	pkg := &PackageJSON{
		Dependencies:     make(map[string]string),
		DevDependencies:  make(map[string]string),
		PeerDependencies: make(map[string]string),
	}

	// Marshmallow unmarshals into struct AND returns complete map with all fields
	raw, err := marshmallow.Unmarshal(data, pkg)
	if err != nil {
		return nil, errors.Wrap(err, "parsing package.json")
	}

	pkg.raw = raw
	return pkg, nil
}

// Write writes a PackageJSON back to disk with proper formatting.
// Merges struct fields back into the raw map to preserve unknown fields.
// Note: Field order from the original file is not preserved due to Go map iteration being random.
func Write(fs afero.Fs, path string, pkg *PackageJSON) error {
	// Start with raw map (contains all fields including unknown ones)
	output := pkg.raw
	if output == nil {
		output = make(map[string]any)
	}

	// Marshal struct to get updated field values, then merge into output
	structData, err := json.Marshal(pkg)
	if err != nil {
		return errors.Wrap(err, "marshalling struct")
	}

	var structMap map[string]any
	if err = json.Unmarshal(structData, &structMap); err != nil {
		return errors.Wrap(err, "unmarshalling struct map")
	}

	// Merge struct fields into output (updates modified fields, preserves unknown fields)
	for key, value := range structMap {
		output[key] = value
	}

	// Marshal final map with indentation and without HTML escaping
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(output); err != nil {
		return errors.Wrap(err, "marshalling package.json")
	}

	data := buffer.Bytes()
	// Encoder adds a newline, but we'll add our own after removing it
	data = bytes.TrimSuffix(data, []byte("\n"))

	// Add trailing newline (standard convention)
	data = append(data, '\n')

	if err := afero.WriteFile(fs, path, data, 0o644); err != nil {
		return errors.Wrap(err, "writing package.json")
	}

	return nil
}
