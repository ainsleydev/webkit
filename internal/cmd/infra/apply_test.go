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

func TestApply(t *testing.T) {
	t.SkipNow()

	t.Run("Init Error", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, teardown := setup(t, &appdef.Definition{}, mock, true)
		defer teardown()

		err := Apply(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "init error")
	})

	t.Run("Apply Error", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Apply(gomock.Any(), env.Production).
			Return(infra.ApplyOutput{
				Output: "Error: Failed to provision resource\nTerraform failed",
			}, errors.New("terraform apply failed"))
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, buf, teardown := setupWithPrinter(t, &appdef.Definition{}, mock, false)
		defer teardown()

		err := Apply(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "executing terraform apply")
		assert.Contains(t, buf.String(), "Failed to provision resource")
	})

	t.Run("Success", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Apply(gomock.Any(), env.Production).
			Return(infra.ApplyOutput{
				Output: "Apply complete! Resources: 2 added, 0 changed, 0 destroyed.",
			}, nil)
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, buf, teardown := setupWithPrinter(t, &appdef.Definition{}, mock, false)
		defer teardown()

		err := Apply(t.Context(), input)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Apply complete")
		assert.Contains(t, buf.String(), "Apply succeeded, see console output")
	})

	t.Run("No Changes To Apply", func(t *testing.T) {
		mock := mockinfra.NewMockManager(gomock.NewController(t))
		mock.EXPECT().
			Apply(gomock.Any(), env.Production).
			Return(infra.ApplyOutput{
				Output: "Apply complete! Resources: 0 added, 0 changed, 0 destroyed.",
			}, nil)
		mock.EXPECT().
			Cleanup().
			Times(1)

		input, buf, teardown := setupWithPrinter(t, &appdef.Definition{}, mock, false)
		defer teardown()

		err := Apply(t.Context(), input)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "0 added, 0 changed, 0 destroyed")
	})
}
