package pkgjson

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
)

type (
	// UpdateResult contains information about what dependencies were updated.
	UpdateResult struct {
		Updated     []string
		Skipped     []string // Dependencies skipped because updating would downgrade them.
		OldVersions map[string]string
	}
	// DependencyMatcher is a function that determines whether
	// a dependency should be updated.
	DependencyMatcher func(name string) bool
	// UpdateOptions configures the behaviour of UpdateDependencies.
	UpdateOptions struct {
		// AllowDowngrades permits updating to a lower version than currently installed.
		AllowDowngrades bool
	}
)

// UpdateDependencies updates dependencies that match the provided matcher function.
// It updates dependencies across all dependency types (dependencies, devDependencies, peerDependencies).
// By default, downgrades are prevented. Pass UpdateOptions with AllowDowngrades to override.
//
// The versionFormatter function determines how to format the version string for each dependency.
// For example, regular dependencies might use "^1.0.0" while devDependencies use "1.0.0".
func UpdateDependencies(
	pkg *PackageJSON,
	matcher DependencyMatcher,
	versionFormatter func(name, version string) string,
	opts ...UpdateOptions,
) *UpdateResult {
	var options UpdateOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	result := &UpdateResult{
		Updated:     []string{},
		Skipped:     []string{},
		OldVersions: make(map[string]string),
	}

	updateDep := func(deps map[string]string, name, oldVer string) {
		newVer := versionFormatter(name, oldVer)
		if !options.AllowDowngrades && IsDowngrade(oldVer, newVer) {
			result.Skipped = append(result.Skipped, name)
			return
		}
		result.Updated = append(result.Updated, name)
		result.OldVersions[name] = oldVer
		deps[name] = newVer
	}

	// Update regular dependencies.
	for name, oldVer := range pkg.Dependencies {
		if matcher(name) {
			updateDep(pkg.Dependencies, name, oldVer)
		}
	}

	// Update devDependencies.
	for name, oldVer := range pkg.DevDependencies {
		if matcher(name) {
			updateDep(pkg.DevDependencies, name, oldVer)
		}
	}

	// Update peerDependencies.
	for name, oldVer := range pkg.PeerDependencies {
		if matcher(name) {
			updateDep(pkg.PeerDependencies, name, oldVer)
		}
	}

	return result
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

// StripVersionPrefix removes common version prefixes (^, ~, >=, >, <=, <, =)
// from a version string, returning the raw semver portion.
func StripVersionPrefix(version string) string {
	version = strings.TrimSpace(version)
	for _, prefix := range []string{">=", "<=", "^", "~", ">", "<", "="} {
		if strings.HasPrefix(version, prefix) {
			return strings.TrimSpace(version[len(prefix):])
		}
	}
	return version
}

// IsDowngrade returns true if updating from oldVersion to newVersion would
// result in a version downgrade. If either version cannot be parsed as
// valid semver, returns false (allowing the update to proceed).
func IsDowngrade(oldVersion, newVersion string) bool {
	oldParsed, err := semver.NewVersion(StripVersionPrefix(oldVersion))
	if err != nil {
		return false
	}
	newParsed, err := semver.NewVersion(StripVersionPrefix(newVersion))
	if err != nil {
		return false
	}
	return newParsed.LessThan(oldParsed)
}
