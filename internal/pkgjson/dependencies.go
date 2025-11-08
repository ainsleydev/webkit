package pkgjson

import (
	"fmt"
	"strings"
)

// DependencyMatcher is a function that determines whether a dependency should be updated.
type DependencyMatcher func(name string) bool

// UpdateResult contains information about what dependencies were updated.
type UpdateResult struct {
	Updated     []string
	OldVersions map[string]string
}

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

// HasDependency checks if a package.json has a specific dependency in any dependency type.
func HasDependency(pkg *PackageJSON, name string) bool {
	if _, ok := pkg.Dependencies[name]; ok {
		return true
	}
	if _, ok := pkg.DevDependencies[name]; ok {
		return true
	}
	if _, ok := pkg.PeerDependencies[name]; ok {
		return true
	}
	return false
}

// HasAnyDependency checks if a package.json has any dependencies matching the matcher.
func HasAnyDependency(pkg *PackageJSON, matcher DependencyMatcher) bool {
	for name := range pkg.Dependencies {
		if matcher(name) {
			return true
		}
	}
	for name := range pkg.DevDependencies {
		if matcher(name) {
			return true
		}
	}
	for name := range pkg.PeerDependencies {
		if matcher(name) {
			return true
		}
	}
	return false
}

// IsDevDependency checks if a package is in devDependencies.
func IsDevDependency(pkg *PackageJSON, name string) bool {
	_, ok := pkg.DevDependencies[name]
	return ok
}

// PayloadMatcher returns a matcher that matches payload and @payloadcms/* packages.
func PayloadMatcher() DependencyMatcher {
	return func(name string) bool {
		return name == "payload" || strings.HasPrefix(name, "@payloadcms/")
	}
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

// MergeDependencies merges dependencies from source into target.
// Existing dependencies in target are preserved unless overwrite is true.
func MergeDependencies(target, source map[string]string, overwrite bool) {
	for name, version := range source {
		if _, exists := target[name]; !exists || overwrite {
			target[name] = version
		}
	}
}
