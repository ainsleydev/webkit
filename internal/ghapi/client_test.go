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
	releases []mockRelease
}

type mockVersion struct {
	tags      []string
	createdAt time.Time
}

type mockRelease struct {
	tagName    string
	draft      bool
	preRelease bool
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

func (m *mockClient) GetLatestRelease(_ context.Context, _, _ string) (string, error) {
	// Find the first stable release (not draft, not pre-release).
	for _, release := range m.releases {
		if !release.draft && !release.preRelease {
			// Remove 'v' prefix if present.
			if len(release.tagName) > 0 && release.tagName[0] == 'v' {
				return release.tagName[1:], nil
			}
			return release.tagName, nil
		}
	}
	return "", nil
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

func TestGetLatestRelease(t *testing.T) {
	t.Parallel()

	t.Run("Returns latest stable release", func(t *testing.T) {
		t.Parallel()

		mock := &mockClient{
			releases: []mockRelease{
				{tagName: "v3.2.0", draft: false, preRelease: false},
				{tagName: "v3.1.0", draft: false, preRelease: false},
			},
		}

		version, err := mock.GetLatestRelease(context.Background(), "payloadcms", "payload")
		assert.NoError(t, err)
		assert.Equal(t, "3.2.0", version)
	})

	t.Run("Skips draft releases", func(t *testing.T) {
		t.Parallel()

		mock := &mockClient{
			releases: []mockRelease{
				{tagName: "v3.3.0", draft: true, preRelease: false},
				{tagName: "v3.2.0", draft: false, preRelease: false},
			},
		}

		version, err := mock.GetLatestRelease(context.Background(), "payloadcms", "payload")
		assert.NoError(t, err)
		assert.Equal(t, "3.2.0", version)
	})

	t.Run("Skips pre-release versions", func(t *testing.T) {
		t.Parallel()

		mock := &mockClient{
			releases: []mockRelease{
				{tagName: "v3.3.0-beta.1", draft: false, preRelease: true},
				{tagName: "v3.2.0", draft: false, preRelease: false},
			},
		}

		version, err := mock.GetLatestRelease(context.Background(), "payloadcms", "payload")
		assert.NoError(t, err)
		assert.Equal(t, "3.2.0", version)
	})

	t.Run("Handles version without v prefix", func(t *testing.T) {
		t.Parallel()

		mock := &mockClient{
			releases: []mockRelease{
				{tagName: "3.2.0", draft: false, preRelease: false},
			},
		}

		version, err := mock.GetLatestRelease(context.Background(), "payloadcms", "payload")
		assert.NoError(t, err)
		assert.Equal(t, "3.2.0", version)
	})

	t.Run("Returns empty string if no stable releases", func(t *testing.T) {
		t.Parallel()

		mock := &mockClient{
			releases: []mockRelease{
				{tagName: "v3.3.0", draft: true, preRelease: false},
				{tagName: "v3.2.0-beta", draft: false, preRelease: true},
			},
		}

		version, err := mock.GetLatestRelease(context.Background(), "payloadcms", "payload")
		assert.NoError(t, err)
		assert.Equal(t, "", version)
	})

	t.Run("Returns empty string for empty releases", func(t *testing.T) {
		t.Parallel()

		mock := &mockClient{
			releases: []mockRelease{},
		}

		version, err := mock.GetLatestRelease(context.Background(), "payloadcms", "payload")
		assert.NoError(t, err)
		assert.Equal(t, "", version)
	})
}
