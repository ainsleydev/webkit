//go:build integration

package git_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/git"
	"github.com/ainsleydev/webkit/internal/util/executil"
)

func TestIntegration_Clone(t *testing.T) {
	t.Skip()

	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()

	client, err := git.New(executil.DefaultRunner())
	require.NoError(t, err)

	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "webkit")

	t.Log("Clone Repo")
	{
		cfg := git.CloneConfig{
			URL:       "https://github.com/ainsleydev/webkit.git",
			LocalPath: repoPath,
			Ref:       "main",
			Depth:     1,
		}

		err = client.Clone(ctx, cfg)
		require.NoError(t, err)
	}

	t.Logf("Verify the repository was cloned at %s", repoPath)
	{
		assert.True(t, git.IsRepository(repoPath))
	}

	t.Log("Verify .git directory exists")
	{
		gitDir := filepath.Join(repoPath, ".git")
		_, err = os.Stat(gitDir)
		require.NoError(t, err)
	}

	t.Log("Verify some WebKit files exist")
	{
		readmePath := filepath.Join(repoPath, "README.md")
		_, err = os.Stat(readmePath)
		assert.NoError(t, err, "README.md should exist")
	}

	if testing.Verbose() {
		logTopLevelContents(t, repoPath)
	}
}

func TestIntegration_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, err := git.New(executil.DefaultRunner())
	require.NoError(t, err)

	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "webkit")

	t.Log("Clone repository initially")
	{
		cfg := git.CloneConfig{
			URL:       "https://github.com/ainsleydev/webkit.git",
			LocalPath: repoPath,
			Ref:       "main",
			Depth:     1,
		}

		err = client.Clone(ctx, cfg)
		require.NoError(t, err)
	}

	t.Log("Update the repository")
	{
		err = client.Update(ctx, repoPath, "main")
		require.NoError(t, err)
	}

	t.Log("Verify repository is still valid")
	{
		assert.True(t, git.IsRepository(repoPath))
	}

	if testing.Verbose() {
		logTopLevelContents(t, repoPath)
	}
}

func TestIntegration_CloneOrUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, err := git.New(executil.DefaultRunner())
	require.NoError(t, err)

	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "webkit")

	cfg := git.CloneConfig{
		URL:       "https://github.com/ainsleydev/webkit.git",
		LocalPath: repoPath,
		Ref:       "main",
		Depth:     1,
	}

	t.Log("First call should clone")
	{
		err = client.CloneOrUpdate(ctx, cfg)
		require.NoError(t, err)
		assert.True(t, git.IsRepository(repoPath))
	}

	t.Log("Second call should update existing repository")
	{
		err = client.CloneOrUpdate(ctx, cfg)
		require.NoError(t, err)
		assert.True(t, git.IsRepository(repoPath))
	}

	if testing.Verbose() {
		logTopLevelContents(t, repoPath)
	}
}

func logTopLevelContents(t *testing.T, path string) {
	t.Helper()

	entries, err := os.ReadDir(path)
	if err != nil {
		t.Logf("Error reading directory: %v", err)
		return
	}

	t.Logf("Contents of %s:", path)
	for _, entry := range entries {
		if entry.IsDir() {
			t.Logf("  %s/", entry.Name())
		} else {
			t.Logf("  %s", entry.Name())
		}
	}
}
