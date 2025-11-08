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
// Fields are written in npm-standard order: name, description, license, private, type, version, etc.
func Write(fs afero.Fs, path string, pkg *PackageJSON) error {
	// Update raw map with current struct values
	if pkg.raw == nil {
		pkg.raw = make(map[string]any)
	}

	// Marshal struct to get updated field values
	structData, err := json.Marshal(pkg)
	if err != nil {
		return errors.Wrap(err, "marshalling struct")
	}

	var structMap map[string]any
	if err = json.Unmarshal(structData, &structMap); err != nil {
		return errors.Wrap(err, "unmarshalling struct map")
	}

	// Merge struct fields into raw map (preserves unknown fields)
	for key, value := range structMap {
		pkg.raw[key] = value
	}

	// Marshal with field ordering preserved
	data, err := marshalOrdered(pkg.raw)
	if err != nil {
		return errors.Wrap(err, "marshalling package.json")
	}

	// Add trailing newline (standard convention)
	data = append(data, '\n')

	if err := afero.WriteFile(fs, path, data, 0o644); err != nil {
		return errors.Wrap(err, "writing package.json")
	}

	return nil
}

// marshalOrdered marshals a package.json map with fields in standard npm order.
// Standard order: name, description, license, private, type, version, scripts, dependencies, etc.
func marshalOrdered(m map[string]any) ([]byte, error) {
	// Define the standard field order for package.json
	fieldOrder := []string{
		"name",
		"description",
		"license",
		"private",
		"type",
		"version",
		"scripts",
		"dependencies",
		"devDependencies",
		"peerDependencies",
		"packageManager",
		"engines",
		"workspaces",
		"repository",
		"keywords",
		"author",
		"contributors",
		"maintainers",
		"homepage",
		"bugs",
		"funding",
		"files",
		"main",
		"module",
		"browser",
		"bin",
		"man",
		"directories",
		"config",
		"pnpm",
		"overrides",
		"resolutions",
	}

	buffer := &bytes.Buffer{}
	buffer.WriteString("{\n")

	written := make(map[string]bool)
	first := true

	// Write fields in defined order
	for _, key := range fieldOrder {
		if value, exists := m[key]; exists {
			if !first {
				buffer.WriteString(",\n")
			}
			first = false

			// Marshal the key
			keyJSON, err := json.Marshal(key)
			if err != nil {
				return nil, err
			}

			// Marshal the value without HTML escaping
			var valueJSON []byte
			valueBuf := &bytes.Buffer{}
			encoder := json.NewEncoder(valueBuf)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(value); err != nil {
				return nil, err
			}
			valueJSON = bytes.TrimSuffix(valueBuf.Bytes(), []byte("\n"))

			// Write field with proper indentation
			buffer.WriteString("\t")
			buffer.Write(keyJSON)
			buffer.WriteString(": ")

			// Handle multiline values (objects and arrays)
			if bytes.Contains(valueJSON, []byte("\n")) {
				// Add indentation to each line
				lines := bytes.Split(valueJSON, []byte("\n"))
				for i, line := range lines {
					if i > 0 {
						buffer.WriteString("\n\t")
					}
					buffer.Write(line)
				}
			} else {
				buffer.Write(valueJSON)
			}

			written[key] = true
		}
	}

	// Write any remaining fields not in the standard order
	for key, value := range m {
		if !written[key] {
			if !first {
				buffer.WriteString(",\n")
			}
			first = false

			keyJSON, err := json.Marshal(key)
			if err != nil {
				return nil, err
			}

			var valueJSON []byte
			valueBuf := &bytes.Buffer{}
			encoder := json.NewEncoder(valueBuf)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(value); err != nil {
				return nil, err
			}
			valueJSON = bytes.TrimSuffix(valueBuf.Bytes(), []byte("\n"))

			buffer.WriteString("\t")
			buffer.Write(keyJSON)
			buffer.WriteString(": ")

			if bytes.Contains(valueJSON, []byte("\n")) {
				lines := bytes.Split(valueJSON, []byte("\n"))
				for i, line := range lines {
					if i > 0 {
						buffer.WriteString("\n\t")
					}
					buffer.Write(line)
				}
			} else {
				buffer.Write(valueJSON)
			}
		}
	}

	buffer.WriteString("\n}")
	return buffer.Bytes(), nil
}
