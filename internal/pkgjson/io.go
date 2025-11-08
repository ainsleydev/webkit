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

	// Extract all strongly-typed fields from raw map.
	if name, ok := raw["name"].(string); ok {
		pkg.Name = name
	}
	if version, ok := raw["version"].(string); ok {
		pkg.Version = version
	}
	if description, ok := raw["description"].(string); ok {
		pkg.Description = description
	}
	if license, ok := raw["license"].(string); ok {
		pkg.License = license
	}
	if pkgType, ok := raw["type"].(string); ok {
		pkg.Type = pkgType
	}
	if packageManager, ok := raw["packageManager"].(string); ok {
		pkg.PackageManager = packageManager
	}
	if homepage, ok := raw["homepage"].(string); ok {
		pkg.Homepage = homepage
	}
	if main, ok := raw["main"].(string); ok {
		pkg.Main = main
	}
	if module, ok := raw["module"].(string); ok {
		pkg.Module = module
	}

	// Extract any-type fields.
	if private, ok := raw["private"]; ok {
		pkg.Private = private
	}
	if scripts, ok := raw["scripts"]; ok {
		pkg.Scripts, _ = scripts.(map[string]any)
	}
	if engines, ok := raw["engines"]; ok {
		pkg.Engines, _ = engines.(map[string]any)
	}
	if workspaces, ok := raw["workspaces"]; ok {
		pkg.Workspaces = workspaces
	}
	if repository, ok := raw["repository"]; ok {
		pkg.Repository = repository
	}
	if author, ok := raw["author"]; ok {
		pkg.Author = author
	}
	if contributors, ok := raw["contributors"]; ok {
		pkg.Contributors = contributors
	}
	if maintainers, ok := raw["maintainers"]; ok {
		pkg.Maintainers = maintainers
	}
	if bugs, ok := raw["bugs"]; ok {
		pkg.Bugs = bugs
	}
	if funding, ok := raw["funding"]; ok {
		pkg.Funding = funding
	}
	if browser, ok := raw["browser"]; ok {
		pkg.Browser = browser
	}
	if bin, ok := raw["bin"]; ok {
		pkg.Bin = bin
	}
	if man, ok := raw["man"]; ok {
		pkg.Man = man
	}
	if directories, ok := raw["directories"]; ok {
		pkg.Directories = directories
	}
	if config, ok := raw["config"]; ok {
		pkg.Config = config
	}
	if pnpm, ok := raw["pnpm"]; ok {
		pkg.Pnpm = pnpm
	}
	if overrides, ok := raw["overrides"]; ok {
		pkg.Overrides = overrides
	}
	if resolutions, ok := raw["resolutions"]; ok {
		pkg.Resolutions = resolutions
	}

	// Extract string arrays.
	if keywords, ok := raw["keywords"].([]interface{}); ok {
		pkg.Keywords = make([]string, 0, len(keywords))
		for _, kw := range keywords {
			if str, ok := kw.(string); ok {
				pkg.Keywords = append(pkg.Keywords, str)
			}
		}
	}
	if files, ok := raw["files"].([]interface{}); ok {
		pkg.Files = make([]string, 0, len(files))
		for _, f := range files {
			if str, ok := f.(string); ok {
				pkg.Files = append(pkg.Files, str)
			}
		}
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
	// Update the raw map with all strongly-typed fields.
	if pkg.Name != "" {
		pkg.raw["name"] = pkg.Name
	}
	if pkg.Version != "" {
		pkg.raw["version"] = pkg.Version
	}
	if pkg.Description != "" {
		pkg.raw["description"] = pkg.Description
	}
	if pkg.License != "" {
		pkg.raw["license"] = pkg.License
	}
	if pkg.Type != "" {
		pkg.raw["type"] = pkg.Type
	}
	if pkg.PackageManager != "" {
		pkg.raw["packageManager"] = pkg.PackageManager
	}
	if pkg.Homepage != "" {
		pkg.raw["homepage"] = pkg.Homepage
	}
	if pkg.Main != "" {
		pkg.raw["main"] = pkg.Main
	}
	if pkg.Module != "" {
		pkg.raw["module"] = pkg.Module
	}

	// Update any-type fields.
	if pkg.Private != nil {
		pkg.raw["private"] = pkg.Private
	}
	if pkg.Scripts != nil {
		pkg.raw["scripts"] = pkg.Scripts
	}
	if pkg.Engines != nil {
		pkg.raw["engines"] = pkg.Engines
	}
	if pkg.Workspaces != nil {
		pkg.raw["workspaces"] = pkg.Workspaces
	}
	if pkg.Repository != nil {
		pkg.raw["repository"] = pkg.Repository
	}
	if pkg.Author != nil {
		pkg.raw["author"] = pkg.Author
	}
	if pkg.Contributors != nil {
		pkg.raw["contributors"] = pkg.Contributors
	}
	if pkg.Maintainers != nil {
		pkg.raw["maintainers"] = pkg.Maintainers
	}
	if pkg.Bugs != nil {
		pkg.raw["bugs"] = pkg.Bugs
	}
	if pkg.Funding != nil {
		pkg.raw["funding"] = pkg.Funding
	}
	if pkg.Browser != nil {
		pkg.raw["browser"] = pkg.Browser
	}
	if pkg.Bin != nil {
		pkg.raw["bin"] = pkg.Bin
	}
	if pkg.Man != nil {
		pkg.raw["man"] = pkg.Man
	}
	if pkg.Directories != nil {
		pkg.raw["directories"] = pkg.Directories
	}
	if pkg.Config != nil {
		pkg.raw["config"] = pkg.Config
	}
	if pkg.Pnpm != nil {
		pkg.raw["pnpm"] = pkg.Pnpm
	}
	if pkg.Overrides != nil {
		pkg.raw["overrides"] = pkg.Overrides
	}
	if pkg.Resolutions != nil {
		pkg.raw["resolutions"] = pkg.Resolutions
	}

	// Update string arrays.
	if pkg.Keywords != nil {
		pkg.raw["keywords"] = pkg.Keywords
	}
	if pkg.Files != nil {
		pkg.raw["files"] = pkg.Files
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
