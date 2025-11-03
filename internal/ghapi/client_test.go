//go:build !race

package ghapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("Creates client with token", func(t *testing.T) {
		t.Parallel()

		client := NewClient("test-token")
		assert.NotNil(t, client)

		dc, ok := client.(*DefaultClient)
		assert.True(t, ok)
		assert.NotNil(t, dc.client)
	})

	t.Run("Creates client without token", func(t *testing.T) {
		t.Parallel()

		client := NewClient("")
		assert.NotNil(t, client)

		dc, ok := client.(*DefaultClient)
		assert.True(t, ok)
		assert.NotNil(t, dc.client)
	})
}

func TestGetLatestSHATag(t *testing.T) {
	t.Parallel()

	t.Run("Returns empty string for non-existent package", func(t *testing.T) {
		t.Parallel()

		client := NewClient("")
		tag := client.GetLatestSHATag(
			context.Background(),
			"non-existent-owner",
			"non-existent-repo",
			"non-existent-app",
		)

		assert.Equal(t, "", tag)
	})

	t.Run("Returns empty string with invalid inputs", func(t *testing.T) {
		t.Parallel()

		client := NewClient("")
		tag := client.GetLatestSHATag(
			context.Background(),
			"",
			"",
			"",
		)

		assert.Equal(t, "", tag)
	})
}

// mockGHClient is a test client that returns predefined responses.
type mockGHClient struct {
	tags []string
}

func (m *mockGHClient) GetLatestSHATag(ctx context.Context, owner, repo, appName string) string {
	var shaTags []string
	for _, tag := range m.tags {
		if len(tag) > 4 && tag[:4] == "sha-" {
			shaTags = append(shaTags, tag)
		}
	}

	if len(shaTags) == 0 {
		return ""
	}

	// Simple sorting - return last alphabetically.
	// In real implementation, this is sorted properly.
	latest := shaTags[0]
	for _, tag := range shaTags {
		if tag > latest {
			latest = tag
		}
	}
	return latest
}

func TestMockClient(t *testing.T) {
	t.Parallel()

	t.Run("Mock returns latest SHA tag", func(t *testing.T) {
		t.Parallel()

		mock := &mockGHClient{
			tags: []string{
				"latest",
				"sha-abc123",
				"sha-def456",
				"v1.0.0",
			},
		}

		tag := mock.GetLatestSHATag(context.Background(), "owner", "repo", "app")
		assert.Equal(t, "sha-def456", tag)
	})

	t.Run("Mock returns empty for no SHA tags", func(t *testing.T) {
		t.Parallel()

		mock := &mockGHClient{
			tags: []string{"latest", "v1.0.0"},
		}

		tag := mock.GetLatestSHATag(context.Background(), "owner", "repo", "app")
		assert.Equal(t, "", tag)
	})

	t.Run("Mock returns empty for empty tags", func(t *testing.T) {
		t.Parallel()

		mock := &mockGHClient{
			tags: []string{},
		}

		tag := mock.GetLatestSHATag(context.Background(), "owner", "repo", "app")
		assert.Equal(t, "", tag)
	})
}
