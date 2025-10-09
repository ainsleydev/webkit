package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/cmdutil"
)

func TestCloneConfigValidate(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   CloneConfig
		wantErr bool
	}{
		"Valid Config": {
			input: CloneConfig{
				URL:       "https://github.com/user/repo.git",
				LocalPath: "/tmp/repo",
				Ref:       "main",
				Depth:     0,
			},
			wantErr: false,
		},
		"Missing URL": {
			input: CloneConfig{
				URL:       "",
				LocalPath: "/tmp/repo",
				Ref:       "main",
				Depth:     0,
			},
			wantErr: true,
		},
		"Missing LocalPath": {
			input: CloneConfig{
				URL:       "https://github.com/user/repo.git",
				LocalPath: "",
				Ref:       "main",
				Depth:     0,
			},
			wantErr: true,
		},
		"Missing URL And LocalPath": {
			input: CloneConfig{
				URL:       "",
				LocalPath: "",
				Ref:       "main",
				Depth:     0,
			},
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := test.input.Validate()
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

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

func TestCloneOrUpdate(t *testing.T) {
	t.Parallel()

	t.Run("Clone When Repository Does Not Exist", func(t *testing.T) {
		t.Parallel()

		client, mock := setupClient(t)
		localPath := t.TempDir() + "/repo"
		cfg := CloneConfig{
			URL:       "https://example.com/repo.git",
			LocalPath: localPath,
			Ref:       "main",
		}

		mock.AddStub("git clone", cmdutil.Result{Output: "cloned"})

		err := client.CloneOrUpdate(t.Context(), cfg)
		assert.NoError(t, err)
	})

	t.Run("Update When Repository Exists", func(t *testing.T) {
		t.Parallel()

		client, mock := setupClient(t)
		localPath := t.TempDir() + "/repo"
		touchGitDir(t, localPath)

		mock.AddStub("git fetch", cmdutil.Result{Output: "fetched"})
		mock.AddStub("git reset", cmdutil.Result{Output: "reset"})

		cfg := CloneConfig{
			URL:       "https://example.com/repo.git",
			LocalPath: localPath,
			Ref:       "main",
		}

		err := client.CloneOrUpdate(t.Context(), cfg)
		assert.NoError(t, err)
	})

	t.Run("Validation Error", func(t *testing.T) {
		t.Parallel()

		client, _ := setupClient(t)
		localPath := t.TempDir() + "/repo"

		cfg := CloneConfig{
			URL:       "",
			LocalPath: localPath,
		}

		err := client.CloneOrUpdate(t.Context(), cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "URL is required")
	})

	t.Run("Runner Clone Error", func(t *testing.T) {
		t.Parallel()

		client, mock := setupClient(t)
		localPath := t.TempDir() + "/repo"

		cfg := CloneConfig{
			URL:       "https://example.com/repo.git",
			LocalPath: localPath,
			Ref:       "main",
		}

		mock.AddStub("git clone", cmdutil.Result{Err: assert.AnError})

		err := client.CloneOrUpdate(t.Context(), cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "git clone failed")
	})

	t.Run("Runner Update Error", func(t *testing.T) {
		t.Parallel()

		client, mock := setupClient(t)
		localPath := t.TempDir() + "/repo"
		touchGitDir(t, localPath)

		mock.AddStub("git fetch", cmdutil.Result{Err: assert.AnError})

		cfg := CloneConfig{
			URL:       "https://example.com/repo.git",
			LocalPath: localPath,
			Ref:       "main",
		}

		err := client.CloneOrUpdate(t.Context(), cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "git fetch failed")
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
