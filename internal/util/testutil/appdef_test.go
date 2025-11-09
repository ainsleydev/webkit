package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestNewAppDefWithDefaults(t *testing.T) {
	t.Parallel()

	t.Run("Applies defaults successfully", func(t *testing.T) {
		t.Parallel()

		def := &appdef.Definition{
			Project: appdef.Project{
				Name:        "test-project",
				Description: "Test description",
			},
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "./apps/web",
				},
			},
		}

		result := NewAppDefWithDefaults(t, def)

		assert.NotNil(t, result)
		assert.Equal(t, "test-project", result.Project.Name)
		assert.Len(t, result.Apps, 1)

		// Verify defaults were applied (SvelteKit apps get default commands).
		assert.NotEmpty(t, result.Apps[0].Commands)
	})

	t.Run("Handles empty project", func(t *testing.T) {
		t.Parallel()

		def := &appdef.Definition{
			Project: appdef.Project{
				Name: "empty-project",
			},
		}

		result := NewAppDefWithDefaults(t, def)

		assert.NotNil(t, result)
		assert.Equal(t, "empty-project", result.Project.Name)
	})
}
