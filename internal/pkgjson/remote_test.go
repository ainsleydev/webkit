package pkgjson

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchFromRemote(t *testing.T) {
	t.Parallel()

	t.Run("Bad URL", func(t *testing.T) {
		t.Parallel()

		_, err := FetchFromRemote(nil, "https://github.com") // Nil context returns error.
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "creating request")
	})

	t.Run("Fetches and parses remote package.json", func(t *testing.T) {
		t.Parallel()

		// Create a test HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"name": "test-package",
				"version": "1.0.0",
				"dependencies": {
					"react": "^18.0.0",
					"vue": "^3.0.0"
				},
				"devDependencies": {
					"typescript": "^5.0.0"
				}
			}`))
		}))
		defer server.Close()

		pkg, err := FetchFromRemote(t.Context(), server.URL)
		require.NoError(t, err)
		require.NotNil(t, pkg)

		assert.Equal(t, "test-package", pkg.Name)
		assert.Equal(t, "1.0.0", pkg.Version)
		assert.Len(t, pkg.Dependencies, 2)
		assert.Equal(t, "^18.0.0", pkg.Dependencies["react"])
		assert.Equal(t, "^3.0.0", pkg.Dependencies["vue"])
		assert.Len(t, pkg.DevDependencies, 1)
		assert.Equal(t, "^5.0.0", pkg.DevDependencies["typescript"])
	})

	t.Run("Returns error for non-200 status", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		_, err := FetchFromRemote(t.Context(), server.URL)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code")
	})

	t.Run("Returns error for invalid JSON", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{invalid json}`))
		}))
		defer server.Close()

		_, err := FetchFromRemote(t.Context(), server.URL)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parsing package.json")
	})

	t.Run("Handles context cancellation", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(t.Context())
		cancel() // Cancel immediately

		_, err := FetchFromRemote(ctx, "http://example.com")
		assert.Error(t, err)
	})
}
