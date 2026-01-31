package pkgjson

import (
	"sort"
)

type (
	// PackageJSON represents a package.json file structure.
	// This type preserves all fields during read/write operations while providing
	// strongly-typed access to common fields.
	//
	// Unknown fields are automatically captured by marshmallow during unmarshaling
	// and merged back during marshaling, ensuring complete preservation.
	PackageJSON struct {
		// Strongly-typed fields for common package.json properties.
		Name             string            `json:"name,omitempty"`
		Description      string            `json:"description,omitempty"`
		License          string            `json:"license,omitempty"`
		Private          any               `json:"private,omitempty"` // Can be string or bool
		Type             string            `json:"type,omitempty"`
		Version          string            `json:"version,omitempty"`
		Scripts          map[string]any    `json:"scripts,omitempty"`
		Dependencies     map[string]string `json:"dependencies,omitempty"`
		DevDependencies  map[string]string `json:"devDependencies,omitempty"`
		PeerDependencies map[string]string `json:"peerDependencies,omitempty"`
		PackageManager   string            `json:"packageManager,omitempty"`
		Engines          map[string]any    `json:"engines,omitempty"`
		Workspaces       any               `json:"workspaces,omitempty"`
		Repository       any               `json:"repository,omitempty"`
		Keywords         []string          `json:"keywords,omitempty"`
		Author           any               `json:"author,omitempty"` // Can be string or object
		Contributors     any               `json:"contributors,omitempty"`
		Maintainers      any               `json:"maintainers,omitempty"`
		Homepage         string            `json:"homepage,omitempty"`
		Bugs             any               `json:"bugs,omitempty"`
		Funding          any               `json:"funding,omitempty"`
		Files            []string          `json:"files,omitempty"`
		Main             string            `json:"main,omitempty"`
		Module           string            `json:"module,omitempty"`
		Exports          any               `json:"exports,omitempty"`
		Imports          any               `json:"imports,omitempty"`
		Browser          any               `json:"browser,omitempty"`
		Bin              any               `json:"bin,omitempty"`
		Man              any               `json:"man,omitempty"`
		Directories      any               `json:"directories,omitempty"`
		Config           any               `json:"config,omitempty"`
		Pnpm             any               `json:"pnpm,omitempty"`
		Overrides        any               `json:"overrides,omitempty"`
		Resolutions      any               `json:"resolutions,omitempty"`

		// Raw stores unknown fields
		raw map[string]any
	}
	// PnpmConfig represents pnpm-specific configuration in package.json.
	PnpmConfig struct {
		OnlyBuiltDependencies []string `json:"onlyBuiltDependencies,omitempty"`
	}
	// Author represents a package.json author or contributor.
	// Can be used for author, contributors, and maintainers fields.
	Author struct {
		Name  string `json:"name"`
		Email string `json:"email,omitempty"`
		URL   string `json:"url,omitempty"`
	}
)

// HasDependency checks if the package has a specific dependency
// in any dependency type.
func (pkg *PackageJSON) HasDependency(name string) bool {
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

// HasAnyDependency checks if the package has any dependencies
// matching the matcher.
func (pkg *PackageJSON) HasAnyDependency(matcher DependencyMatcher) bool {
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
func (pkg *PackageJSON) IsDevDependency(name string) bool {
	_, ok := pkg.DevDependencies[name]
	return ok
}

// GetDependencyVersion returns the version string for a dependency across all dependency types.
// Returns an empty string if the dependency is not found in any type.
func (pkg *PackageJSON) GetDependencyVersion(name string) string {
	if ver := pkg.Dependencies[name]; ver != "" {
		return ver
	}
	if ver := pkg.DevDependencies[name]; ver != "" {
		return ver
	}
	if ver := pkg.PeerDependencies[name]; ver != "" {
		return ver
	}
	return ""
}

// sortDependencies sorts all dependency maps alphabetically in-place.
func (pkg *PackageJSON) sortDependencies() {
	pkg.Dependencies = sortMap(pkg.Dependencies)
	pkg.DevDependencies = sortMap(pkg.DevDependencies)
	pkg.PeerDependencies = sortMap(pkg.PeerDependencies)
}

// sortMap returns a new map with the same entries but sorted keys.
// Go 1.12+ preserves insertion order when marshaling to JSON.
func sortMap(m map[string]string) map[string]string {
	if len(m) == 0 {
		return m
	}

	// Get sorted keys
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build new map with sorted insertion order
	sorted := make(map[string]string, len(m))
	for _, k := range keys {
		sorted[k] = m[k]
	}

	return sorted
}
