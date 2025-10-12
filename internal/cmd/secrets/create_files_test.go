package secrets

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/util/testutil"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestCreateSecretFiles(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		input, _ := setup(t, &appdef.Definition{})

		err := CreateFiles(t.Context(), input)
		assert.NoError(t, err)

		t.Log(".sops.yaml Created")
		{
			exists, err := afero.Exists(input.FS, "resources/.sops.yaml")
			assert.NoError(t, err)
			assert.True(t, exists)

			content, err := afero.ReadFile(input.FS, "resources/.sops.yaml")
			require.NoError(t, err)
			assert.Contains(t, string(content), "creation_rules")
			assert.Contains(t, string(content), "secrets/.*\\.yaml$")
			assert.Contains(t, string(content), "age1")
		}

		t.Log("Secret Files Created")
		{
			for _, enviro := range env.All {
				path := "resources/secrets/" + enviro + ".yaml"

				exists, err := afero.Exists(input.FS, path)
				assert.NoError(t, err)
				assert.True(t, exists)

				file, err := afero.ReadFile(input.FS, path)
				assert.NoError(t, err)
				assert.Empty(t, string(file))
			}
		}
	})

	t.Run("SOPS Config Error", func(t *testing.T) {
		t.Parallel()

		input, _ := setup(t, &appdef.Definition{})
		input.FS = &testutil.AferoErrCreateFs{Fs: afero.NewMemMapFs()}

		got := CreateFiles(t.Context(), input)
		assert.Error(t, got)
	})
}
