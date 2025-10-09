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

func IsRepository(path string) bool {
	gitDir := filepath.Join(path, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}
