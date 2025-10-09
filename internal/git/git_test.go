package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/cmdutil"
)

func setupClient(t *testing.T) (*Client, *cmdutil.MemRunner) {
	t.Helper()
	mock := cmdutil.NewMemRunner()
	client, err := New(mock)
	require.NoError(t, err)
	return client, mock
}

//func TestNewGitCommandMissing(t *testing.T) {
//	t.Parallel()
//
//	t.Run("Git Not Found In PATH", func(t *testing.T) {
//		t.Parallel()
//
//		// Patch cmdutil.Exists to simulate git missing
//		origExists := cmdutil.Exists
//		cmdutil.Exists = func(name string) bool { return false }
//		defer func() { cmdutil.Exists = origExists }()
//
//		_, err := New(nil)
//		assert.Error(t, err)
//		assert.Contains(t, err.Error(), "git command not found")
//	})
//}

func TestIsRepository(t *testing.T) {
	t.Parallel()

	t.Run("Git Directory Exists", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		err := os.Mkdir(filepath.Join(tmpDir, ".git"), os.ModePerm)
		require.NoError(t, err)

		got := IsRepository(tmpDir)
		assert.True(t, got)
	})

	t.Run("Git Directory Does Not Exist", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()

		got := IsRepository(tmpDir)
		assert.False(t, got)
	})
}
