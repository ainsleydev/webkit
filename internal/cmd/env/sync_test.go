package env

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/util/testutil"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestSync(t *testing.T) {
	t.Parallel()

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
						"FOO": {Value: "bar"},
						"BAZ": {Value: "qux"},
					},
				},
			},
			{
				Name: "app2",
				Path: "app2/nested",
				Env: appdef.Environment{
					Production: appdef.EnvVar{
						"HELLO": {Value: "world"},
					},
				},
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		input, _ := setup(t, appDef)

		err := Sync(t.Context(), input)
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
		t.Parallel()

		orig := dotEnvMarshaller
		defer func() { dotEnvMarshaller = orig }()
		dotEnvMarshaller = func(envMap map[string]string) (string, error) {
			return "", fmt.Errorf("marshal error")
		}

		input, _ := setup(t, appDef)

		err := Sync(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "marshal error")
	})

	t.Run("Write Error", func(t *testing.T) {
		t.Parallel()

		input, _ := setup(t, appDef)
		input.FS = &testutil.AferoErrCreateFs{Fs: afero.NewMemMapFs()}

		err := Sync(t.Context(), input)
		assert.Error(t, err)
	})
}
