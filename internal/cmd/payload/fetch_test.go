package payload

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/mocks"
)

func TestFetchPayloadDependencies(t *testing.T) {
	t.Parallel()

	t.Run("Fetches and parses dependencies", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mock := mocks.NewGHClient(ctrl)

		packageJSON := []byte(`{
			"name": "payload",
			"version": "3.0.0",
			"dependencies": {
				"@lexical/headless": "0.28.0",
				"lexical": "0.28.0",
				"react": "^18.0.0"
			},
			"devDependencies": {
				"typescript": "^5.0.0"
			}
		}`)

		mock.EXPECT().
			GetFileContent(gomock.Any(), "payloadcms", "payload", "package.json", "v3.0.0").
			Return(packageJSON, nil).
			Times(1)

		deps, err := fetchPayloadDependencies(context.Background(), mock, "3.0.0")
		require.NoError(t, err)
		require.NotNil(t, deps)

		assert.Equal(t, "0.28.0", deps.Dependencies["lexical"])
		assert.Equal(t, "^18.0.0", deps.Dependencies["react"])
		assert.Equal(t, "^5.0.0", deps.DevDependencies["typescript"])

		// Verify AllDeps contains everything.
		assert.Equal(t, 4, len(deps.AllDeps))
		assert.Contains(t, deps.AllDeps, "lexical")
		assert.Contains(t, deps.AllDeps, "typescript")
	})

	t.Run("Handles empty dependencies", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mock := mocks.NewGHClient(ctrl)

		packageJSON := []byte(`{
			"name": "payload",
			"version": "3.0.0"
		}`)

		mock.EXPECT().
			GetFileContent(gomock.Any(), "payloadcms", "payload", "package.json", "v3.0.0").
			Return(packageJSON, nil).
			Times(1)

		deps, err := fetchPayloadDependencies(context.Background(), mock, "3.0.0")
		require.NoError(t, err)
		require.NotNil(t, deps)

		assert.Empty(t, deps.Dependencies)
		assert.Empty(t, deps.DevDependencies)
		assert.Empty(t, deps.AllDeps)
	})

	t.Run("Returns error for invalid JSON", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mock := mocks.NewGHClient(ctrl)

		invalidJSON := []byte(`{invalid json}`)

		mock.EXPECT().
			GetFileContent(gomock.Any(), "payloadcms", "payload", "package.json", "v3.0.0").
			Return(invalidJSON, nil).
			Times(1)

		_, err := fetchPayloadDependencies(context.Background(), mock, "3.0.0")
		assert.Error(t, err)
	})
}
