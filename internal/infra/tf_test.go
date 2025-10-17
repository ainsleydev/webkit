package infra

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/util/executil"
)

func TestNewTerraform(t *testing.T) {
	if !executil.Exists("terraform") {
		t.Skip("terraform not found in PATH")
	}

	t.Run("Success", func(t *testing.T) {
		tf, err := NewTerraform(t.Context())
		require.NoError(t, err)
		assert.NotNil(t, tf)
		assert.NotEmpty(t, tf.path)
		assert.Contains(t, tf.path, "terraform")
	})

	t.Run("TerraformNotInPath", func(t *testing.T) {
		t.Setenv("PATH", "/nonexistent")

		_, err := NewTerraform(t.Context())
		assert.Error(t, err)
	})
}

func TestTerraform_Init(t *testing.T) {
	t.Skip()

	if !executil.Exists("terraform") {
		t.Skip("terraform not found in PATH")
	}

	t.Run("Success", func(t *testing.T) {
		tf, err := NewTerraform(t.Context())
		require.NoError(t, err)
		defer tf.Cleanup()

		err = tf.Init(t.Context())
		require.NoError(t, err)

		// Verify temp directory was created
		assert.NotEmpty(t, tf.tmpDir)
		assert.DirExists(t, tf.tmpDir)

		// Verify base directory exists
		baseDir := filepath.Join(tf.tmpDir, "base")
		assert.DirExists(t, baseDir)

		// Verify terraform was initialized (.terraform directory)
		terraformDir := filepath.Join(baseDir, ".terraform")
		assert.DirExists(t, terraformDir)

		// Verify tf executor was created
		assert.NotNil(t, tf.tf)
	})

	t.Run("InitTwice", func(t *testing.T) {
		tf, err := NewTerraform(t.Context())
		require.NoError(t, err)
		defer tf.Cleanup()

		err = tf.Init(t.Context())
		require.NoError(t, err)

		firstTmpDir := tf.tmpDir

		// Init again should create new temp dir
		err = tf.Init(t.Context())
		require.NoError(t, err)

		assert.NotEqual(t, firstTmpDir, tf.tmpDir)
	})
}

func TestTerraform_Cleanup(t *testing.T) {
	t.Skip()

	if !executil.Exists("terraform") {
		t.Skip("terraform not found in PATH")
	}

	t.Run("Removes Dir", func(t *testing.T) {
		tf, err := NewTerraform(t.Context())
		require.NoError(t, err)

		err = tf.Init(t.Context())
		require.NoError(t, err)

		tmpDir := tf.tmpDir
		assert.DirExists(t, tmpDir)

		tf.Cleanup()

		assert.NoDirExists(t, tmpDir)
	})

	t.Run("NoOpIfNotInitialized", func(t *testing.T) {
		tf := &Terraform{
			path: "/usr/bin/terraform",
		}

		assert.NotPanics(t, func() {
			tf.Cleanup()
		})
	})
}
