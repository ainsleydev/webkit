package payload

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockGHClient struct {
	files map[string][]byte
}

func (m *mockGHClient) GetLatestSHATag(_ context.Context, _, _, _ string) (string, error) {
	return "", nil
}

func (m *mockGHClient) GetLatestRelease(_ context.Context, _, _ string) (string, error) {
	return "3.0.0", nil
}

func (m *mockGHClient) GetFileContent(_ context.Context, _, _, path, ref string) ([]byte, error) {
	key := path + "@" + ref
	return m.files[key], nil
}

func TestFetchPayloadDependencies(t *testing.T) {
	t.Parallel()

	t.Run("Fetches and parses dependencies", func(t *testing.T) {
		t.Parallel()

		mock := &mockGHClient{
			files: map[string][]byte{
				"package.json@v3.0.0": []byte(`{
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
				}`),
			},
		}

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

		mock := &mockGHClient{
			files: map[string][]byte{
				"package.json@v3.0.0": []byte(`{
					"name": "payload",
					"version": "3.0.0"
				}`),
			},
		}

		deps, err := fetchPayloadDependencies(context.Background(), mock, "3.0.0")
		require.NoError(t, err)
		require.NotNil(t, deps)

		assert.Empty(t, deps.Dependencies)
		assert.Empty(t, deps.DevDependencies)
		assert.Empty(t, deps.AllDeps)
	})

	t.Run("Returns error for invalid JSON", func(t *testing.T) {
		t.Parallel()

		mock := &mockGHClient{
			files: map[string][]byte{
				"package.json@v3.0.0": []byte(`{invalid json}`),
			},
		}

		_, err := fetchPayloadDependencies(context.Background(), mock, "3.0.0")
		assert.Error(t, err)
	})
}
