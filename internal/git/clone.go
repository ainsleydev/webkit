package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ainsleydev/webkit/internal/cmdutil"
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

// Clone clones a git repository to the specified path
func (c Client) Clone(ctx context.Context, cfg CloneConfig) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(cfg.LocalPath), os.ModePerm); err != nil {
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

	cmd := cmdutil.NewCommand("git", args...)
	res := c.Runner.Run(ctx, cmd)
	if res.Err != nil {
		return fmt.Errorf("git clone failed: %w", res.Err)
	}

	return nil
}

func (c Client) CloneOrUpdate(ctx context.Context, cfg CloneConfig) error {
	if IsRepository(cfg.LocalPath) {
		ref := cfg.Ref
		if ref == "" {
			ref = "main" // or "master", depending on your needs
		}
		return c.Update(ctx, cfg.LocalPath, ref)
	}
	return c.Clone(ctx, cfg)
}

func (c Client) Update(ctx context.Context, repoPath, ref string) error {
	if !IsRepository(repoPath) {
		return fmt.Errorf("not a git repository: %s", repoPath)
	}

	fetchCmd := cmdutil.NewCommand("git", "fetch", "origin", ref)
	fetchCmd.Dir = repoPath
	if res := c.Runner.Run(ctx, fetchCmd); res.Err != nil {
		return fmt.Errorf("git fetch failed: %w", res.Err)
	}

	resetCmd := cmdutil.NewCommand("git", "reset", "--hard", "origin/"+ref)
	resetCmd.Dir = repoPath
	if res := c.Runner.Run(ctx, resetCmd); res.Err != nil {
		return fmt.Errorf("git reset failed: %w", res.Err)
	}

	return nil
}
