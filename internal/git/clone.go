package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ainsleydev/webkit/internal/util/executil"
)

// CloneConfig configures repository cloning
type CloneConfig struct {
	URL       string // URL of the git repository.
	Ref       string // Branch, tag, or commit to checkout.
	LocalPath string // LocalPath where to clone the repository
	Depth     int    // Depth for shallow clone (0 = full clone)
}

// Validate config before using it
func (cfg CloneConfig) Validate() error {
	if cfg.URL == "" {
		return fmt.Errorf("URL is required")
	}
	if cfg.LocalPath == "" {
		return fmt.Errorf("LocalPath is required")
	}
	return nil
}

// Clone creates a new local copy of a remote repository.
// Parent directories are created automatically to avoid manual setup.
func (c Client) Clone(ctx context.Context, cfg CloneConfig) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	// Create parent directory to avoid "no such file or directory" errors
	if err := os.MkdirAll(filepath.Dir(cfg.LocalPath), os.ModePerm); err != nil {
		return fmt.Errorf("creating parent dir: %w", err)
	}

	args := []string{"clone"}

	// Shallow clone reduces download time and disk usage for large repos
	if cfg.Depth > 0 {
		args = append(args, "--depth", fmt.Sprintf("%d", cfg.Depth))
	}

	// Allows cloning specific branches/tags without fetching all refs
	if cfg.Ref != "" {
		args = append(args, "--branch", cfg.Ref)
	}

	args = append(args, cfg.URL, cfg.LocalPath)

	cmd := executil.NewCommand("git", args...)
	_, err := c.Runner.Run(ctx, cmd)
	if err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}

// CloneOrUpdate ensures a repository exists and is up to date.
// Useful for idempotent operations where initial state
// is not important.
func (c Client) CloneOrUpdate(ctx context.Context, cfg CloneConfig) error {
	if IsRepository(cfg.LocalPath) {
		ref := cfg.Ref
		if ref == "" {
			// Defaults to main to match conventions.
			ref = "main"
		}
		return c.Update(ctx, cfg.LocalPath, ref)
	}
	return c.Clone(ctx, cfg)
}

// Update synchronizes a local repository with the remote ref.
// Uses reset --hard to discard local changes, ensuring a
// clean state matching remote.
func (c Client) Update(ctx context.Context, repoPath, ref string) error {
	if !IsRepository(repoPath) {
		return fmt.Errorf("not a git repository: %s", repoPath)
	}

	// Fetch updates the remote tracking branch without modifying working directory
	fetchCmd := executil.NewCommand("git", "fetch", "origin", ref)
	fetchCmd.Dir = repoPath
	_, err := c.Runner.Run(ctx, fetchCmd)
	if err != nil {
		return fmt.Errorf("git fetch failed: %w", err)
	}

	// Hard reset discards local changes to ensure consistency with remote
	resetCmd := executil.NewCommand("git", "reset", "--hard", "origin/"+ref)
	resetCmd.Dir = repoPath
	_, err = c.Runner.Run(ctx, resetCmd)
	if err != nil {
		return fmt.Errorf("git reset failed: %w", err)
	}

	return nil
}
