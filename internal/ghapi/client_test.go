package ghapi

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("Creates client with token", func(t *testing.T) {
		t.Parallel()
		client := New("test-token")
		assert.NotNil(t, client)
	})
}

type mockClient struct {
	versions []mockVersion
}

type mockVersion struct {
	tags      []string
	createdAt time.Time
}

func (m *mockClient) GetLatestSHATag(_ context.Context, _, _, _ string) (string, error) {
	var shaTags []struct {
		tag       string
		createdAt time.Time
	}

	for _, v := range m.versions {
		for _, t := range v.tags {
			if t != "" && len(t) >= 4 && t[:4] == "sha-" {
				shaTags = append(shaTags, struct {
					tag       string
					createdAt time.Time
				}{t, v.createdAt})
			}
		}
	}

	if len(shaTags) == 0 {
		return "", nil
	}

	// Sort by creation time descending
	for i := 0; i < len(shaTags)-1; i++ {
		for j := i + 1; j < len(shaTags); j++ {
			if shaTags[j].createdAt.After(shaTags[i].createdAt) {
				shaTags[i], shaTags[j] = shaTags[j], shaTags[i]
			}
		}
	}

	return shaTags[0].tag, nil
}

func TestGetLatestSHATag(t *testing.T) {
	t.Parallel()

	t.Run("Returns latest SHA tag by creation time", func(t *testing.T) {
		t.Parallel()

		mock := &mockClient{
			versions: []mockVersion{
				{tags: []string{"sha-aaa"}, createdAt: time.Date(2025, 11, 1, 10, 0, 0, 0, time.UTC)},
				{tags: []string{"sha-bbb"}, createdAt: time.Date(2025, 11, 1, 11, 0, 0, 0, time.UTC)},
				{tags: []string{"latest"}, createdAt: time.Date(2025, 11, 1, 12, 0, 0, 0, time.UTC)},
			},
		}

		tag, err := mock.GetLatestSHATag(context.Background(), "owner", "repo", "app")
		assert.NoError(t, err)
		assert.Equal(t, "sha-bbb", tag)
	})

	t.Run("Returns empty string if no SHA tags", func(t *testing.T) {
		t.Parallel()

		mock := &mockClient{
			versions: []mockVersion{
				{tags: []string{"latest"}, createdAt: time.Now()},
				{tags: []string{"v1.0.0"}, createdAt: time.Now()},
			},
		}

		tag, err := mock.GetLatestSHATag(context.Background(), "owner", "repo", "app")
		assert.NoError(t, err)
		assert.Equal(t, "", tag)
	})

	t.Run("Returns empty string for empty versions", func(t *testing.T) {
		t.Parallel()

		mock := &mockClient{
			versions: []mockVersion{},
		}

		tag, err := mock.GetLatestSHATag(context.Background(), "owner", "repo", "app")
		assert.NoError(t, err)
		assert.Equal(t, "", tag)
	})
}
