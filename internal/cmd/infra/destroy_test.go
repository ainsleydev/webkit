package infra

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/infra"
	"github.com/ainsleydev/webkit/internal/infra/mocks"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestDestroy(t *testing.T) {
	t.SkipNow()

	t.Run("Init Error", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, teardown := setup(t, &appdef.Definition{}, mock, true)
		defer teardown()

		err := Destroy(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "init error")
	})

	t.Run("Destroy Error", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Destroy(gomock.Any(), env.Production).
			Return(infra.DestroyOutput{
				Output: "Error: Failed to destroy resource\nTerraform failed",
			}, errors.New("terraform destroy failed"))
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, buf, teardown := setupWithPrinter(t, &appdef.Definition{}, mock, false)
		defer teardown()

		err := Destroy(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "executing terraform destroy")
		assert.Contains(t, buf.String(), "Failed to destroy resource")
	})

	t.Run("Success", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Destroy(gomock.Any(), env.Production).
			Return(infra.DestroyOutput{
				Output: "Destroy complete! Resources: 2 destroyed.",
			}, nil)
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, buf, teardown := setupWithPrinter(t, &appdef.Definition{}, mock, false)
		defer teardown()

		err := Destroy(t.Context(), input)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Destroy complete")
		assert.Contains(t, buf.String(), "Destroy succeeded, see console output")
	})

	t.Run("Filters Unmanaged Apps And Resources", func(t *testing.T) {
		falseVal := false
		def := &appdef.Definition{
			Project: appdef.Project{Name: "test-project"},
			Apps: []appdef.App{
				{Name: "managed-app", TerraformManaged: nil},
				{Name: "unmanaged-app", TerraformManaged: &falseVal},
			},
			Resources: []appdef.Resource{
				{Name: "managed-db", TerraformManaged: nil},
				{Name: "unmanaged-cache", TerraformManaged: &falseVal},
			},
		}

		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Destroy(gomock.Any(), env.Production).
			Return(infra.DestroyOutput{
				Output: "Destroy complete! Resources: 1 destroyed.",
			}, nil)
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, buf, teardown := setupWithPrinter(t, def, mock, false)
		defer teardown()

		err := Destroy(t.Context(), input)
		assert.NoError(t, err)

		output := buf.String()
		// Verify skipped items are displayed
		assert.Contains(t, output, "The following items are not managed by Terraform:")
		assert.Contains(t, output, "unmanaged-app")
		assert.Contains(t, output, "unmanaged-cache")
		assert.NotContains(t, output, "managed-app")
		assert.NotContains(t, output, "managed-db")
	})

	t.Run("No Output When All Managed", func(t *testing.T) {
		def := &appdef.Definition{
			Project: appdef.Project{Name: "test-project"},
			Apps: []appdef.App{
				{Name: "app1", TerraformManaged: nil},
				{Name: "app2", TerraformManaged: nil},
			},
			Resources: []appdef.Resource{
				{Name: "db", TerraformManaged: nil},
			},
		}

		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Destroy(gomock.Any(), env.Production).
			Return(infra.DestroyOutput{
				Output: "Destroy complete! Resources: 2 destroyed.",
			}, nil)
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, buf, teardown := setupWithPrinter(t, def, mock, false)
		defer teardown()

		err := Destroy(t.Context(), input)
		assert.NoError(t, err)

		output := buf.String()
		// Should not show skipped items message when all are managed
		assert.NotContains(t, output, "The following items are not managed by Terraform:")
	})
}
