package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/cmdutil"
)

func TestClone(t *testing.T) {
	t.Parallel()

	t.Run("Successful Clone", func(t *testing.T) {
		t.Parallel()

		client, mock := setupClient(t)
		mock.AddStub("git clone", cmdutil.Result{Output: "cloned!"})

		cfg := CloneConfig{
			URL:       "https://example.com/repo.git",
			LocalPath: "/tmp/repo",
			Ref:       "main",
			Depth:     1,
		}

		err := client.Clone(t.Context(), cfg)
		assert.NoError(t, err)
	})

	t.Run("Validation Error", func(t *testing.T) {
		t.Parallel()

		client, _ := setupClient(t)

		cfg := CloneConfig{
			URL:       "",
			LocalPath: "/tmp/repo",
		}

		err := client.Clone(t.Context(), cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "URL is required")
	})

	t.Run("Runner Error", func(t *testing.T) {
		t.Parallel()

		client, mock := setupClient(t)
		mock.AddStub("git clone", cmdutil.Result{Err: assert.AnError})

		cfg := CloneConfig{
			URL:       "https://example.com/repo.git",
			LocalPath: "/tmp/repo",
		}

		err := client.Clone(t.Context(), cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "git clone failed")
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("Successful Update", func(t *testing.T) {
		t.Parallel()

		client, mock := setupClient(t)
		mock.AddStub("git fetch", cmdutil.Result{Output: "fetched"})
		mock.AddStub("git reset", cmdutil.Result{Output: "reset"})

		repoPath := t.TempDir()
		touchGitDir(t, repoPath)

		err := client.Update(t.Context(), repoPath, "main")
		assert.NoError(t, err)
	})

	t.Run("Not a Git Repository", func(t *testing.T) {
		t.Parallel()

		client, _ := setupClient(t)
		repoPath := t.TempDir()

		err := client.Update(t.Context(), repoPath, "main")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a git repository")
	})

	t.Run("Runner Fetch Error", func(t *testing.T) {
		t.Parallel()

		client, mock := setupClient(t)
		mock.AddStub("git fetch", cmdutil.Result{Err: assert.AnError})

		repoPath := t.TempDir()
		touchGitDir(t, repoPath)

		err := client.Update(t.Context(), repoPath, "main")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "git fetch failed")
	})

	t.Run("Runner Reset Error", func(t *testing.T) {
		t.Parallel()

		client, mock := setupClient(t)
		mock.AddStub("git fetch", cmdutil.Result{Output: "fetched"})
		mock.AddStub("git reset", cmdutil.Result{Err: assert.AnError})

		repoPath := t.TempDir()
		touchGitDir(t, repoPath)

		err := client.Update(t.Context(), repoPath, "main")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "git reset failed")
	})
}

func touchGitDir(t *testing.T, path string) {
	t.Helper()
	err := os.MkdirAll(filepath.Join(path, ".git"), 0o755)
	require.NoError(t, err)
}
