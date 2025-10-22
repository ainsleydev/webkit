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
}
