package secrets

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/mocks"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestScaffold(t *testing.T) {
	t.Parallel()

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input, _ := setup(t, &appdef.Definition{})
		input.FS = afero.NewReadOnlyFs(afero.NewMemMapFs())

		err := Scaffold(t.Context(), input)
		assert.Error(t, err)
	})

	t.Run("Scaffold Error", func(t *testing.T) {
		t.Skip()
		ctrl := gomock.NewController(t)

		fsMock := mocks.NewMockFS(ctrl)

		// BasePathFs will call MkdirAll on "resources"
		fsMock.EXPECT().
			MkdirAll(gomock.Eq("resources"), gomock.Any()).
			Return(nil).
			AnyTimes()

		// .sops.yaml should succeed
		fsMock.EXPECT().
			Stat(gomock.Eq("resources/.sops.yaml")).
			Return(nil, nil).
			AnyTimes()
		fsMock.EXPECT().
			OpenFile(gomock.Eq("resources/.sops.yaml"), gomock.Any(), gomock.Any()).
			Return(nil, nil).
			AnyTimes()

		// Fail for secrets/*.yaml
		fsMock.EXPECT().
			Stat(gomock.Cond(func(path string) bool {
				return strings.HasPrefix(path, "resources/secrets/")
			})).
			Return(nil, fmt.Errorf("gen.Bytes error")).
			AnyTimes()
		fsMock.EXPECT().
			OpenFile(gomock.Cond(func(path string) bool {
				return strings.HasPrefix(path, "resources/secrets/")
			}), gomock.Any(), gomock.Any()).
			Return(nil, fmt.Errorf("gen.Bytes error")).
			AnyTimes()

		input, _ := setup(t, &appdef.Definition{})
		input.FS = fsMock

		err := Scaffold(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "mkdir error")
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		input, _ := setup(t, &appdef.Definition{})

		err := Scaffold(t.Context(), input)
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
				path := "resources/secrets/" + enviro.String() + ".yaml"

				exists, err := afero.Exists(input.FS, path)
				assert.NoError(t, err)
				assert.True(t, exists)

				file, err := afero.ReadFile(input.FS, path)
				assert.NoError(t, err)
				assert.Empty(t, string(file))
			}
		}
	})
}
