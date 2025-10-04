package cmd

// PrettierConfig represents a .prettierrc configuration.
// Based on: https://json.schemastore.org/prettierrc
// Generated reference: prettier_base.go (not used directly)
type (
	PrettierConfig struct {
		UseTabs       bool               `json:"useTabs,omitempty"`
		SingleQuote   bool               `json:"singleQuote,omitempty"`
		TrailingComma string             `json:"trailingComma"` // "all" | "es5" | "none"
		PrintWidth    int                `json:"printWidth"`
		TabWidth      int                `json:"tabWidth"`
		Semi          bool               `json:"semi"`
		Plugins       []string           `json:"plugins,omitempty"`
		Overrides     []PrettierOverride `json:"overrides,omitempty"`
	}
	PrettierOverride struct {
		Files   []string               `json:"files"`
		Options map[string]interface{} `json:"options"`
	}
)
