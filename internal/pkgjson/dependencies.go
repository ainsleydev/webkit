package pkgjson

import (
	"fmt"
)

type (
	// UpdateResult contains information about what dependencies were updated.
	UpdateResult struct {
		Updated     []string
		OldVersions map[string]string
	}
	// DependencyMatcher is a function that determines whether
	// a dependency should be updated.
	DependencyMatcher func(name string) bool
)

// UpdateDependencies updates dependencies that match the provided matcher function.
// It updates dependencies across all dependency types (dependencies, devDependencies, peerDependencies).
//
// The versionFormatter function determines how to format the version string for each dependency.
// For example, regular dependencies might use "^1.0.0" while devDependencies use "1.0.0".
func UpdateDependencies(
	pkg *PackageJSON,
	matcher DependencyMatcher,
	versionFormatter func(name, version string) string,
) *UpdateResult {
	result := &UpdateResult{
		Updated:     []string{},
		OldVersions: make(map[string]string),
	}

	// Update regular dependencies.
	for name, oldVer := range pkg.Dependencies {
		if matcher(name) {
			result.Updated = append(result.Updated, name)
			result.OldVersions[name] = oldVer
			pkg.Dependencies[name] = versionFormatter(name, oldVer)
		}
	}

	// Update devDependencies.
	for name, oldVer := range pkg.DevDependencies {
		if matcher(name) {
			result.Updated = append(result.Updated, name)
			result.OldVersions[name] = oldVer
			pkg.DevDependencies[name] = versionFormatter(name, oldVer)
		}
	}

	// Update peerDependencies.
	for name, oldVer := range pkg.PeerDependencies {
		if matcher(name) {
			result.Updated = append(result.Updated, name)
			result.OldVersions[name] = oldVer
			pkg.PeerDependencies[name] = versionFormatter(name, oldVer)
		}
	}

	return result
}

// SetMatcher returns a matcher that matches any dependency in the provided set.
func SetMatcher(deps map[string]string) DependencyMatcher {
	return func(name string) bool {
		_, ok := deps[name]
		return ok
	}
}

// FormatVersion formats a version string with optional caret prefix.
// If useExactVersion is true, returns the version as-is.
// Otherwise, adds a caret (^) prefix.
func FormatVersion(version string, useExactVersion bool) string {
	if useExactVersion {
		return version
	}
	return fmt.Sprintf("^%s", version)
}
