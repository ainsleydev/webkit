package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type (
	Manifest struct {
		Version     string               `json:"version"` // Webkit CLI version
		GeneratedAt time.Time            `json:"generated_at"`
		Files       map[string]FileEntry `json:"files"` // Filepath to Entry
	}
	FileEntry struct {
		Path         string    `json:"path"`
		Generator    string    `json:"generator"` // e.g. "cicd.BackupWorkflow", "files.PackageJSON"
		Source       string    `json:"source"`    // What in app.json caused this? e.g., "resource:postgres-prod"
		Hash         string    `json:"hash"`      // SHA256 of content for drift detection
		ScaffoldMode bool      `json:"scaffolded"`
		GeneratedAt  time.Time `json:"generated_at"`
	}
)

// HashContent generates a SHA256 hash of the provided data.
// Used to detect if file contents have changed since generation.
func HashContent(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
