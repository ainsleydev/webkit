package env

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/mocks"
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestSync(t *testing.T) {
	ctx := t.Context()
	ageIdentity, err := age.NewIdentity()
	require.NoError(t, err)

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

		appDef := &appdef.Definition{
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Production: appdef.EnvVar{
						"HELLO": {Value: "world", Source: appdef.EnvSourceSOPS},
					},
				},
			},
		}

		input, _ := setup(t, appDef)
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

		input, _ := setup(t, appDef)
		input.FS = fsMock
		input.SOPSCache = mocks.NewMockEncrypterDecrypter(ctrl)

		err = Sync(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "mkdir error")
	})

	t.Run("Success", func(t *testing.T) {
		t.Setenv(age.KeyEnvVar, ageIdentity.String())
		defer os.Unsetenv(age.KeyEnvVar)

		input, _ := setup(t, appDef)

		err = Sync(ctx, input)
		assert.NoError(t, err)

		for _, app := range appDef.Apps {
			for _, envName := range environmentsWithDotEnv {
				fileName := ".env"
				if envName != env.Development {
					fileName = ".env." + envName.String()
				}
				path := filepath.Join(app.Path, fileName)

				exists, err := afero.Exists(input.FS, path)
				assert.NoError(t, err)
				assert.Truef(t, exists, "%s should exist", path)

				content, err := afero.ReadFile(input.FS, path)
				assert.NoError(t, err)
				assert.NotEmptyf(t, content, "%s should not be empty", path)

				// Only check Production env for merged vars
				if envName == env.Production {
					for k, v := range app.Env.Production {
						assert.Contains(t, string(content), k)
						assert.Contains(t, string(content), v.Value)
					}
				}
			}
		}
	})

	t.Run("Marshal Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		orig := dotEnvMarshaller
		defer func() { dotEnvMarshaller = orig }()
		dotEnvMarshaller = func(envMap map[string]string) (string, error) {
			return "", fmt.Errorf("marshal error")
		}

		input, _ := setup(t, appDef)
		input.SOPSCache = mocks.NewMockEncrypterDecrypter(ctrl)

		err = Sync(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "marshal error")
	})
}
