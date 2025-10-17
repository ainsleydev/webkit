package infra

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestWriteBackendConfig(t *testing.T) {
	infraDir := "/infra"

	t.Run("Failure", func(t *testing.T) {
		tf := setup(t, &appdef.Definition{})
		defer tf.Cleanup()

		tf.fs = afero.NewReadOnlyFs(afero.NewMemMapFs())

		_, err := tf.writeS3Backend(infraDir, env.Production)
		assert.Error(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		tf := setup(t, &appdef.Definition{
			Project: appdef.Project{Name: "test-project"},
		})
		defer tf.Cleanup()

		tf.fs = afero.NewMemMapFs()

		gotPath, err := tf.writeS3Backend(infraDir, env.Production)
		require.NoError(t, err)

		t.Log("backend.hcl")
		{
			assert.Equal(t, filepath.Join(infraDir, backendHclFileName), gotPath)

			exists, err := afero.Exists(tf.fs, gotPath)
			assert.NoError(t, err)
			assert.True(t, exists, "backend.hcl file should exist")

			content, err := afero.ReadFile(tf.fs, gotPath)
			assert.NoError(t, err)
			contentStr := string(content)
			assert.Contains(t, contentStr, "bucket")
			assert.Contains(t, contentStr, tf.env.BackBlazeBucket)
			assert.Contains(t, contentStr, "test-project/production/terraform.tfstate")
			assert.Contains(t, contentStr, tf.env.BackBlazeKeyID)
			assert.Contains(t, contentStr, tf.env.BackBlazeApplicationKey)
		}

		t.Log("backend.tf")
		{
			backendTfPath := filepath.Join(infraDir, backendTfFileName)
			exists, err := afero.Exists(tf.fs, backendTfPath)
			assert.NoError(t, err)
			assert.True(t, exists, "backend.tf file should exist")

			content, err := afero.ReadFile(tf.fs, backendTfPath)
			assert.NoError(t, err)
			contentStr := string(content)
			assert.Contains(t, contentStr, "terraform")
			assert.Contains(t, contentStr, `backend "s3"`)
		}
	})
}
