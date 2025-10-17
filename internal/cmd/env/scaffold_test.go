package env

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/mocks"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestScaffold(t *testing.T) {
	t.Parallel()

	appDef := &appdef.Definition{
		WebkitVersion: "",
		Project:       appdef.Project{},
		Shared:        appdef.Shared{},
		Resources:     nil,
		Apps: []appdef.App{
			{Name: "app1", Path: "./app1"},
			{Name: "app2", Path: "app2/nested"},
		},
	}

	t.Run("FS Failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		fsMock := mocks.NewMockFS(ctrl)
		fsMock.EXPECT().
			MkdirAll(gomock.Any(), gomock.Any()).
			Return(fmt.Errorf("mkdir error"))

		input, _ := setup(t, appDef)
		input.FS = fsMock
		input.SOPSCache = mocks.NewMockEncrypterDecrypter(ctrl)

		err := Scaffold(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "mkdir error")
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		input, _ := setup(t, appDef)

		err := Scaffold(t.Context(), input)
		assert.NoError(t, err)

		for _, app := range appDef.Apps {
			for _, envName := range environmentsWithDotEnv {
				if testing.Verbose() {
					t.Logf("Checking .env file for app %s (%s)", app.Name, envName)
				}

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
				assert.Contains(t, string(content), scaffold.WebKitNotice)
			}
		}
	})
}
