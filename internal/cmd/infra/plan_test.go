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

func TestPlan(t *testing.T) {
	t.SkipNow()

	t.Run("Init Error", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, teardown := setup(t, &appdef.Definition{}, mock, true)
		defer teardown()

		err := Plan(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "init error")
	})

	t.Run("Plan Error", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Plan(gomock.Any(), env.Production).
			Return(infra.PlanOutput{}, errors.New("plan error"))
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, teardown := setup(t, &appdef.Definition{}, mock, false)
		defer teardown()

		err := Plan(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "plan error")
	})

	t.Run("Success", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Plan(gomock.Any(), env.Production).
			Return(infra.PlanOutput{
				HasChanges: true,
				Output:     "plan output for test",
				Plan:       nil,
			}, nil)
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, buf, teardown := setupWithPrinter(t, &appdef.Definition{}, mock, false)
		defer teardown()

		err := Plan(t.Context(), input)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "plan output for test")
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
			Plan(gomock.Any(), env.Production).
			Return(infra.PlanOutput{
				HasChanges: true,
				Output:     "plan output",
				Plan:       nil,
			}, nil)
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, buf, teardown := setupWithPrinter(t, def, mock, false)
		defer teardown()

		err := Plan(t.Context(), input)
		assert.NoError(t, err)

		output := buf.String()
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
			Plan(gomock.Any(), env.Production).
			Return(infra.PlanOutput{
				HasChanges: true,
				Output:     "plan output",
				Plan:       nil,
			}, nil)
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, buf, teardown := setupWithPrinter(t, def, mock, false)
		defer teardown()

		err := Plan(t.Context(), input)
		assert.NoError(t, err)

		output := buf.String()
		// Should not show skipped items message when all are managed
		assert.NotContains(t, output, "The following items are not managed by Terraform:")
	})
}
