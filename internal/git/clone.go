package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CloneConfig configures repository cloning
type CloneConfig struct {
	// URL of the git repository
	URL string

	// Branch, tag, or commit to checkout
	Ref string

	// LocalPath where to clone the repository
	LocalPath string

	// Depth for shallow clone (0 = full clone)
	Depth int
}

// Clone clones a git repository to the specified path
func Clone(cfg CloneConfig) error {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(cfg.LocalPath), 0755); err != nil {
		return fmt.Errorf("creating parent dir: %w", err)
	}

	args := []string{"clone"}

	if cfg.Depth > 0 {
		args = append(args, "--depth", fmt.Sprintf("%d", cfg.Depth))
	}

	if cfg.Ref != "" {
		args = append(args, "--branch", cfg.Ref)
	}

	args = append(args, cfg.URL, cfg.LocalPath)

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}

// Update updates an existing git repository
func Update(repoPath string, ref string) error {
	// Fetch
	fetchCmd := exec.Command("git", "fetch", "origin", ref)
	fetchCmd.Dir = repoPath
	fetchCmd.Stdout = os.Stdout
	fetchCmd.Stderr = os.Stderr

	if err := fetchCmd.Run(); err != nil {
		return fmt.Errorf("git fetch failed: %w", err)
	}

	// Reset to match remote
	resetCmd := exec.Command("git", "reset", "--hard", "origin/"+ref)
	resetCmd.Dir = repoPath
	resetCmd.Stdout = os.Stdout
	resetCmd.Stderr = os.Stderr

	if err := resetCmd.Run(); err != nil {
		return fmt.Errorf("git reset failed: %w", err)
	}

	return nil
}

// IsCloned checks if a directory is a git repository
func IsCloned(path string) bool {
	gitDir := filepath.Join(path, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}
