package secrets

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestSync(t *testing.T) {
	t.Parallel()

	envFixture := appdef.Environment{
		Production: map[string]appdef.EnvValue{
			"SECRET_KEY": {
				Source: appdef.EnvSourceSOPS,
				Value:  "production",
			},
		},
	}

	t.Run("No Files", func(t *testing.T) {
		t.Parallel()

		got := Sync(t.Context(), cmdtools.CommandInput{
			FS:          afero.NewMemMapFs(),
			AppDefCache: &appdef.Definition{},
		})

		assert.NoError(t, got)
		// TODO: Assert "No secrets" in std out.
	})

	t.Run("Scaffold Error", func(t *testing.T) {
		t.Parallel()

		def := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "app1", Env: envFixture},
			},
		}

		fs := afero.NewMemMapFs()
		input := cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: def,
		}

		got := Sync(t.Context(), input)
		assert.NoError(t, got, "No production.yaml file causes error")

		// TODO: Capture stdout
	})

	t.Run("Sync's Shared", func(t *testing.T) {
		t.Parallel()

		def := &appdef.Definition{
			Shared: appdef.Shared{
				Env: envFixture,
			},
		}

		fs := afero.NewMemMapFs()

		got := Sync(t.Context(), cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: def,
		})

		assert.NoError(t, got)
	})

	t.Run("Sync's Apps", func(t *testing.T) {
		t.Parallel()

		def := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "app1", Env: envFixture},
			},
		}

		fs := afero.NewMemMapFs()
		input := cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: def,
		}

		err := CreateFiles(t.Context(), input)
		assert.NoError(t, err)

		got := Sync(t.Context(), input)
		assert.NoError(t, got)

		file, err := afero.ReadFile(fs, secrets.FilePathFromEnv(env.Production))
		assert.NoError(t, err)
		assert.Contains(t, string(file), `SECRET_KEY: "REPLACE_ME_SECRET_KEY"`)
	})
}
