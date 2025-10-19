package manifest

import (
	"time"
)

type (
	Manifest struct {
		Version     string               `json:"version"` // Webkit CLI version
		GeneratedAt time.Time            `json:"generated_at"`
		Files       map[string]FileEntry `json:"files"` // Filepath to Entry
	}
	FileEntry struct {
		Path        string    `json:"path"`
		Generator   string    `json:"generator"` // e.g. "cicd.BackupWorkflow", "files.PackageJSON"
		Source      string    `json:"source"`    // What in app.json caused this? e.g., "resource:postgres-prod"
		Hash        string    `json:"hash"`      // SHA256 of content for drift detection
		GeneratedAt time.Time `json:"generated_at"`
	}
)
