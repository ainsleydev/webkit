package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ainsleydev/webkit/internal/cmdutil"
)

type Client struct {
	Runner cmdutil.Runner
}

// New creates a git client with the provided command runner.
// Validates git is available to fail fast rather than on first operation.
func New(runner cmdutil.Runner) (*Client, error) {
	if !cmdutil.Exists("git") {
		return nil, errors.New("git command not found in $PATH")
	}
	return &Client{Runner: runner}, nil
}

var (
	ErrNotRepository = fmt.Errorf("not a git repository")
	ErrInvalidConfig = fmt.Errorf("invalid configuration")
)

// IsRepository checks for .git directory presence to verify
// repository status.
func IsRepository(path string) bool {
	gitDir := filepath.Join(path, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}
