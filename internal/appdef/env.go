package appdef

import (
	"fmt"
	"strings"
)

type (
	// Environment contains environment-specific variable configurations.
	Environment struct {
		Dev        EnvVar `json:"dev,omitempty"`
		Staging    EnvVar `json:"staging,omitempty"`
		Production EnvVar `json:"production,omitempty"`
	}
	// EnvVar is a map of variable names to their configurations.
	EnvVar map[string]EnvValue
	// EnvValue represents a single environment variable configuration
	EnvValue struct {
		Source EnvSource `json:"source"`          // See below
		Value  string    `json:"value,omitempty"` // Used for "value" and "resource" sources
		Path   string    `json:"path,omitempty"`  // Used for "sops" source (format: "file:key")
	}
)

// EnvSource defines the type of application being run.
type EnvSource string

const (
	// EnvSourceValue is a static string value.
	// Example: "https://api.example.com"
	EnvSourceValue EnvSource = "value"

	// EnvSourceResource references a Terraform resource output.
	// Example: "db.connection_url"
	EnvSourceResource EnvSource = "resource"

	// EnvSourceSOPS is an encrypted secret stored in a SOPS file.
	// Example: "secrets/production.yaml:API_KEY"
	EnvSourceSOPS EnvSource = "sops"
)

// String implements fmt.Stringer on the EnvSource.
func (e EnvSource) String() string {
	return string(e)
}

// SOPSPath represents a parsed SOPS file path and key
type SOPSPath struct {
	// Path to the SOPS encrypted file (e.g., "secrets/production.yaml")
	File string
	// Key within the SOPS file (e.g., "API_KEY")
	Key string
}

// ParseSOPSPath splits the SOPS path into file and key parts.
//
// Example:
// "secrets/production.yaml:PAYLOAD_SECRET"
// ↓
// SOPSPath{File: "secrets/production.yaml", Key: "PAYLOAD_SECRET"}
func (e EnvValue) ParseSOPSPath() (SOPSPath, error) {
	if e.Source != EnvSourceSOPS {
		return SOPSPath{}, fmt.Errorf("not a SOPS source")
	}

	parts := strings.Split(e.Path, ":")
	if len(parts) != 2 {
		return SOPSPath{}, fmt.Errorf("invalid SOPS path format: %s (expected format: file:key)", e.Path)
	}

	return SOPSPath{
		File: parts[0],
		Key:  parts[1],
	}, nil
}
