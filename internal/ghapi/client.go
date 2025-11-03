package ghapi

import (
	"context"
	"sort"
	"strings"

	"github.com/google/go-github/v76/github"
)

//go:generate go tool go.uber.org/mock/mockgen -source=client.go -destination ../mocks/ghapi.go -package=mocks -mock_names=Client=GHClient

// Client provides methods for interacting with the GitHub API.
type Client interface {
	// GetLatestSHATag returns the most recent sha-* tag for a given container image.
	//
	// Returns empty string if no sha tags are found or if the query fails.
	GetLatestSHATag(ctx context.Context, owner, repo, appName string) string
}

// DefaultClient implements the Client interface using the official go-github library.
type DefaultClient struct {
	client *github.Client
}

// New creates a new GitHub API client.
// If token is empty, the client will be unauthenticated (rate limited).
func New(token string) Client {
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
	// For container packages, we need to check both:
	// 1. version.Name - the primary version identifier (often the tag)
	// 2. version.ContainerMetadata.Tags - additional tags (if available)
	shaTagsMap := make(map[string]bool)
	for _, version := range versions {
		// Check version.Name first (primary tag).
		if version.Name != nil && strings.HasPrefix(*version.Name, "sha-") {
			shaTagsMap[*version.Name] = true
		}

		// Also check ContainerMetadata.Tag if available.
		if version.ContainerMetadata != nil &&
			version.ContainerMetadata.Tag != nil &&
			version.ContainerMetadata.Tag.Name != nil {
			if strings.HasPrefix(*version.ContainerMetadata.Tag.Name, "sha-") {
				shaTagsMap[*version.ContainerMetadata.Tag.Name] = true
			}
		}
	}

	if len(shaTagsMap) == 0 {
		return ""
	}

	// Convert map to slice and sort.
	shaTags := make([]string, 0, len(shaTagsMap))
	for tag := range shaTagsMap {
		shaTags = append(shaTags, tag)
	}

	// Sort alphabetically and return the last one (most recent).
	// SHA tags are in format sha-<40-char-hash>, so alphabetical sorting
	// isn't perfect but works reasonably well for recent tags.
	sort.Strings(shaTags)

	return shaTags[len(shaTags)-1]
}
