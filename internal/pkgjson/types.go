package pkgjson

// PackageJSON represents a package.json file structure.
// This type preserves all fields during read/write operations while providing
// strongly-typed access to common fields.
//
// The raw map stores the complete JSON structure, ensuring that fields
// not explicitly defined in the struct are preserved when writing back to disk.
type PackageJSON struct {
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
	Browser          any               `json:"browser,omitempty"`
	Bin              any               `json:"bin,omitempty"`
	Man              any               `json:"man,omitempty"`
	Directories      any               `json:"directories,omitempty"`
	Config           any               `json:"config,omitempty"`
	Pnpm             any               `json:"pnpm,omitempty"`
	Overrides        any               `json:"overrides,omitempty"`
	Resolutions      any               `json:"resolutions,omitempty"`

	// raw stores the complete JSON as a map to preserve all fields.
	// This is particularly important for preserving fields that are not
	// explicitly defined in the struct above.
	raw map[string]any
}
