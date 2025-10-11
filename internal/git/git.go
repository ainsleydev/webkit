package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ainsleydev/webkit/internal/executil"
)

type Client struct {
	Runner executil.Runner
}

// New creates a git client with the provided command runner.
// Validates git is available to fail fast rather than on first operation.
func New(runner executil.Runner) (*Client, error) {
	if !executil.Exists("git") {
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
