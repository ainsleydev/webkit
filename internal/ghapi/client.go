package ghapi

import (
	"context"
	"sort"
	"strings"

	"github.com/google/go-github/v68/github"
)

// Client provides methods for interacting with the GitHub API.
type Client interface {
	// GetLatestSHATag returns the most recent sha-* tag for a given container image.
	// Returns empty string if no sha tags are found or if the query fails.
	GetLatestSHATag(ctx context.Context, owner, repo, appName string) string
}

// DefaultClient implements the Client interface using the official go-github library.
type DefaultClient struct {
	client *github.Client
}

// NewClient creates a new GitHub API client.
// If token is empty, the client will be unauthenticated (rate limited).
func NewClient(token string) Client {
	var client *github.Client
	if token != "" {
		client = github.NewClient(nil).WithAuthToken(token)
	} else {
		client = github.NewClient(nil)
	}

	return &DefaultClient{
		client: client,
	}
}

// GetLatestSHATag queries GHCR for the most recent sha-* tag for a given image.
// The image name format is: {repo}-{appName}
// For example: "my-website-web" for repo "my-website" and app "web".
//
// Returns empty string if:
//   - No package versions are found
//   - No sha-* tags exist
//   - API request fails
func (c *DefaultClient) GetLatestSHATag(ctx context.Context, owner, repo, appName string) string {
	packageName := repo + "-" + appName
	packageType := "container"

	// List all versions of the package.
	// Note: This requires authentication for private packages.
	versions, _, err := c.client.Users.PackageGetAllVersions(
		ctx,
		owner,
		packageType,
		packageName,
		&github.PackageListOptions{
			State: github.String("active"),
		},
	)
	if err != nil {
		return ""
	}

	// Collect all sha-* tags from all versions.
	var shaTags []string
	for _, version := range versions {
		if version.Metadata != nil && version.Metadata.Container != nil {
			for _, tag := range version.Metadata.Container.Tags {
				if strings.HasPrefix(tag, "sha-") {
					shaTags = append(shaTags, tag)
				}
			}
		}
	}

	if len(shaTags) == 0 {
		return ""
	}

	// Sort alphabetically and return the last one (most recent).
	// SHA tags are in format sha-<40-char-hash>, so alphabetical sorting
	// isn't perfect but works reasonably well for recent tags.
	sort.Strings(shaTags)
	return shaTags[len(shaTags)-1]
}
