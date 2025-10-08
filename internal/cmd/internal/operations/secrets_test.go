package operations

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestCreateSecretFiles(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		err := CreateSecretFiles(t.Context(), cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: &appdef.Definition{},
		})
		assert.NoError(t, err)

		t.Log(".sops.yaml Created")
		{
			exists, err := afero.Exists(fs, "resources/.sops.yaml")
			assert.NoError(t, err)
			assert.True(t, exists)

			content, err := afero.ReadFile(fs, "resources/.sops.yaml")
			require.NoError(t, err)
			assert.Contains(t, string(content), "creation_rules")
			assert.Contains(t, string(content), "secrets/.*\\.yaml$")
			assert.Contains(t, string(content), "age1")
		}

		t.Log("Secret Files Created")
		{
			environments := []string{env.Development, env.Staging, env.Production}
			for _, enviro := range environments {
				path := "resources/secrets/" + enviro + ".yaml"

				exists, err := afero.Exists(fs, path)
				assert.NoError(t, err, "File should exist: "+path)
				assert.True(t, exists, "File should exist: "+path)
			}
		}
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		got := CreateSecretFiles(t.Context(), cmdtools.CommandInput{
			FS:          &errCreateFs{Fs: afero.NewMemMapFs()},
			AppDefCache: &appdef.Definition{},
		})
		assert.Error(t, got)
	})
}
