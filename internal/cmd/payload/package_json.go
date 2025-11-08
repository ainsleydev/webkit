package payload

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type (
	// PackageJSON represents a package.json file structure.
	PackageJSON struct {
		Dependencies     map[string]string `json:"dependencies,omitempty"`
		DevDependencies  map[string]string `json:"devDependencies,omitempty"`
		PeerDependencies map[string]string `json:"peerDependencies,omitempty"`
		// We'll store the raw JSON to preserve formatting and other fields.
		raw map[string]interface{}
	}

	// BumpResult contains information about what was bumped.
	BumpResult struct {
		Path        string
		Bumped      []string
		OldVersions map[string]string
		NewVersion  string
	}
)

// ReadPackageJSON reads and parses a package.json file.
func ReadPackageJSON(fs afero.Fs, path string) (*PackageJSON, error) {
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, errors.Wrap(err, "reading package.json")
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, errors.Wrap(err, "parsing package.json")
	}

	pkg := &PackageJSON{
		Dependencies:     make(map[string]string),
		DevDependencies:  make(map[string]string),
		PeerDependencies: make(map[string]string),
		raw:              raw,
	}

	// Extract dependencies.
	if deps, ok := raw["dependencies"].(map[string]interface{}); ok {
		for k, v := range deps {
			if str, ok := v.(string); ok {
				pkg.Dependencies[k] = str
			}
		}
	}

	// Extract devDependencies.
	if deps, ok := raw["devDependencies"].(map[string]interface{}); ok {
		for k, v := range deps {
			if str, ok := v.(string); ok {
				pkg.DevDependencies[k] = str
			}
		}
	}

	// Extract peerDependencies.
	if deps, ok := raw["peerDependencies"].(map[string]interface{}); ok {
		for k, v := range deps {
			if str, ok := v.(string); ok {
				pkg.PeerDependencies[k] = str
			}
		}
	}

	return pkg, nil
}

// WritePackageJSON writes a PackageJSON back to disk with proper formatting.
func WritePackageJSON(fs afero.Fs, path string, pkg *PackageJSON) error {
	// Update the raw map with our changes.
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

// BumpPayloadDependencies updates all Payload-related dependencies to the specified version.
// Returns information about what was changed.
func BumpPayloadDependencies(pkg *PackageJSON, version string) *BumpResult {
	result := &BumpResult{
		Bumped:      []string{},
		OldVersions: make(map[string]string),
		NewVersion:  version,
	}

	// Bump dependencies.
	for name, oldVer := range pkg.Dependencies {
		if isPayloadDependency(name) {
			result.Bumped = append(result.Bumped, name)
			result.OldVersions[name] = oldVer
			pkg.Dependencies[name] = formatVersion(version, false)
		}
	}

	// Bump devDependencies (use exact version without ^).
	for name, oldVer := range pkg.DevDependencies {
		if isPayloadDependency(name) {
			result.Bumped = append(result.Bumped, name)
			result.OldVersions[name] = oldVer
			pkg.DevDependencies[name] = formatVersion(version, true)
		}
	}

	// Bump peerDependencies.
	for name, oldVer := range pkg.PeerDependencies {
		if isPayloadDependency(name) {
			result.Bumped = append(result.Bumped, name)
			result.OldVersions[name] = oldVer
			pkg.PeerDependencies[name] = formatVersion(version, false)
		}
	}

	return result
}

// isPayloadDependency checks if a package name is a Payload CMS dependency.
func isPayloadDependency(name string) bool {
	return name == "payload" || strings.HasPrefix(name, "@payloadcms/")
}

// formatVersion formats a version string with or without the ^ prefix.
// exactMatch=true returns the version as-is (for devDependencies).
// exactMatch=false adds ^ prefix (for dependencies and peerDependencies).
func formatVersion(version string, exactMatch bool) string {
	if exactMatch {
		return version
	}
	return fmt.Sprintf("^%s", version)
}

// HasPayloadDependencies checks if a package.json has any Payload CMS dependencies.
func HasPayloadDependencies(pkg *PackageJSON) bool {
	for name := range pkg.Dependencies {
		if isPayloadDependency(name) {
			return true
		}
	}
	for name := range pkg.DevDependencies {
		if isPayloadDependency(name) {
			return true
		}
	}
	for name := range pkg.PeerDependencies {
		if isPayloadDependency(name) {
			return true
		}
	}
	return false
}
