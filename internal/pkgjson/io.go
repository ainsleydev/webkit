package pkgjson

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Exists checks if a package.json file exists at the given path.
func Exists(fs afero.Fs, path string) (bool, error) {
	exists, err := afero.Exists(fs, path)
	if err != nil {
		return false, errors.Wrap(err, "checking if package.json exists")
	}
	return exists, nil
}

// Read reads and parses a package.json file from the filesystem.
// It preserves all fields in the raw map to ensure no data loss when writing back.
func Read(fs afero.Fs, path string) (*PackageJSON, error) {
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, errors.Wrap(err, "reading package.json")
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, errors.Wrap(err, "parsing package.json")
	}

	pkg := &PackageJSON{
		Dependencies:     make(map[string]string),
		DevDependencies:  make(map[string]string),
		PeerDependencies: make(map[string]string),
		raw:              raw,
	}

	// Extract strongly-typed fields.
	if name, ok := raw["name"].(string); ok {
		pkg.Name = name
	}
	if version, ok := raw["version"].(string); ok {
		pkg.Version = version
	}
	if description, ok := raw["description"].(string); ok {
		pkg.Description = description
	}

	// Extract dependencies.
	extractDependencies(raw, "dependencies", pkg.Dependencies)
	extractDependencies(raw, "devDependencies", pkg.DevDependencies)
	extractDependencies(raw, "peerDependencies", pkg.PeerDependencies)

	return pkg, nil
}

// Write writes a PackageJSON back to disk with proper formatting.
// It preserves all fields from the original raw map and updates modified fields.
func Write(fs afero.Fs, path string, pkg *PackageJSON) error {
	// Update the raw map with our changes.
	if pkg.Name != "" {
		pkg.raw["name"] = pkg.Name
	}
	if pkg.Version != "" {
		pkg.raw["version"] = pkg.Version
	}
	if pkg.Description != "" {
		pkg.raw["description"] = pkg.Description
	}

	// Update dependency maps.
	if len(pkg.Dependencies) > 0 {
		pkg.raw["dependencies"] = pkg.Dependencies
	}
	if len(pkg.DevDependencies) > 0 {
		pkg.raw["devDependencies"] = pkg.DevDependencies
	}
	if len(pkg.PeerDependencies) > 0 {
		pkg.raw["peerDependencies"] = pkg.PeerDependencies
	}

	// Marshal with indentation to match standard package.json formatting.
	data, err := json.MarshalIndent(pkg.raw, "", "\t")
	if err != nil {
		return errors.Wrap(err, "marshalling package.json")
	}

	// Add trailing newline (standard convention).
	data = append(data, '\n')

	if err := afero.WriteFile(fs, path, data, 0o644); err != nil {
		return errors.Wrap(err, "writing package.json")
	}

	return nil
}

// extractDependencies extracts dependencies from the raw map into the target map.
func extractDependencies(raw map[string]any, key string, target map[string]string) {
	if deps, ok := raw[key].(map[string]any); ok {
		for k, v := range deps {
			if str, ok := v.(string); ok {
				target[k] = str
			}
		}
	}
}
