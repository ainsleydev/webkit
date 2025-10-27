package env

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/mocks"
	"github.com/ainsleydev/webkit/internal/secrets/age"
)

func TestSync(t *testing.T) {
	ctx := t.Context()
	ageIdentity, err := age.NewIdentity()
	require.NoError(t, err)

	// App Definition only with values, not secrets, so we don't
	// have to create SOPS files for unit testing env files.
	appDef := &appdef.Definition{
		WebkitVersion: "",
		Project:       appdef.Project{},
		Shared:        appdef.Shared{},
		Resources:     nil,
		Apps: []appdef.App{
			{
				Name: "app1",
				Path: "./app1",
				Env: appdef.Environment{
					Dev: appdef.EnvVar{
						"FOO": {Value: "bar", Source: appdef.EnvSourceValue},
						"BAZ": {Value: "qux", Source: appdef.EnvSourceValue},
					},
					Production: appdef.EnvVar{
						"FOO": {Value: "bar", Source: appdef.EnvSourceValue},
						"BAZ": {Value: "qux", Source: appdef.EnvSourceValue},
					},
				},
			},
			{
				Name: "app2",
				Path: "app2/nested",
				Env: appdef.Environment{
					Dev: appdef.EnvVar{
						"HELLO": {Value: "world", Source: appdef.EnvSourceValue},
					},
					Production: appdef.EnvVar{
						"HELLO": {Value: "world", Source: appdef.EnvSourceValue},
					},
				},
			},
		},
	}

	t.Run("Decrypt Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockEncrypterDecrypter(ctrl)
		mock.EXPECT().
			Decrypt(gomock.Any()).
			Return(fmt.Errorf("decrypt error"))

		appDefSops := &appdef.Definition{
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Production: appdef.EnvVar{
						"HELLO": {Value: "world", Source: appdef.EnvSourceSOPS},
					},
				},
			},
		}

		input := setup(t, appDefSops)
		input.SOPSCache = mock

		err = Sync(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "decrypt error")
	})

	t.Run("Write Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		fsMock := mocks.NewMockFS(ctrl)
		fsMock.EXPECT().
			MkdirAll(gomock.Any(), gomock.Any()).
			Return(fmt.Errorf("mkdir error"))

		input := setup(t, appDef)
		input.FS = fsMock
		input.SOPSCache = mocks.NewMockEncrypterDecrypter(ctrl)

		err = Sync(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "mkdir error")
	})

	t.Run("Success", func(t *testing.T) {
		t.Setenv(age.KeyEnvVar, ageIdentity.String())

		input := setup(t, appDef)

		err = Sync(ctx, input)
		assert.NoError(t, err)

		t.Log("Dev")
		{
			// App 1
			path := filepath.Join("app1", ".env")
			content, err := afero.ReadFile(input.FS, path)
			require.NoError(t, err)

			got := string(content)

			assert.Contains(t, got, "BAZ=\"qux\"")
			assert.Contains(t, got, "FOO=\"bar\"")

			// App 2
			path = filepath.Join("app2/nested", ".env")
			content, err = afero.ReadFile(input.FS, path)
			require.NoError(t, err)

			got = string(content)
			assert.Contains(t, got, "HELLO=\"world\"")
		}

		t.Log("Production")
		{
			// App 1
			path := filepath.Join("app1", ".env.production")
			content, err := afero.ReadFile(input.FS, path)
			require.NoError(t, err)

			got := string(content)
			assert.Contains(t, got, "BAZ=\"qux\"")
			assert.Contains(t, got, "FOO=\"bar\"")

			// App 2
			path = filepath.Join("app2/nested", ".env.production")
			content, err = afero.ReadFile(input.FS, path)
			require.NoError(t, err)

			got = string(content)
			assert.Contains(t, got, "HELLO=\"world\"")
		}
	})

	t.Run("Marshal Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		orig := dotEnvMarshaller
		defer func() { dotEnvMarshaller = orig }()
		dotEnvMarshaller = func(envMap map[string]string) (string, error) {
			return "", fmt.Errorf("marshal error")
		}

		input := setup(t, appDef)
		input.SOPSCache = mocks.NewMockEncrypterDecrypter(ctrl)

		err = Sync(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "marshal error")
	})
}
