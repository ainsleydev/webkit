package ghapi

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v76/github"

	"github.com/ainsleydev/webkit/pkg/enforce"
)

//go:generate go tool go.uber.org/mock/mockgen -source=client.go -destination ../mocks/ghapi.go -package=mocks -mock_names=Client=GHClient

// Client provides methods for interacting with the GitHub API.
type Client interface {
	// GetLatestSHATag returns the most recent sha-* tag for a given container image.
	//
	// Returns empty string if no sha tags are found or if the query fails.
	GetLatestSHATag(ctx context.Context, owner, repo, appName string) (string, error)
}

// DefaultClient implements the Client interface using the official go-github library.
type DefaultClient struct {
	client *github.Client
}

// New creates a new GitHub API client.
// If token is empty, the client will be unauthenticated (rate limited).
func New(token string) Client {
	enforce.NotEqual(token, "", "github token cannot be empty")

	return &DefaultClient{
		client: github.NewClient(nil).WithAuthToken(token),
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
func (c *DefaultClient) GetLatestSHATag(ctx context.Context, owner, repo, appName string) (string, error) {
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
		return "", err
	}

	type shaTag struct {
		Tag       string
		CreatedAt time.Time
	}

	var shaTags []shaTag
	for _, version := range versions {
		var meta github.PackageMetadata
		if err = json.Unmarshal(version.Metadata, &meta); err != nil {
			return "", err
		}
		for _, tag := range meta.GetContainer().Tags {
			if tag != "" && strings.HasPrefix(tag, "sha-") {
				shaTags = append(shaTags, shaTag{
					Tag:       tag,
					CreatedAt: version.GetCreatedAt().Time,
				})
			}
		}
	}

	if len(shaTags) == 0 {
		return "", errors.New("no sha-tags found")
	}

	// Sort by creation date descending (newest first)
	sort.Slice(shaTags, func(i, j int) bool {
		return shaTags[i].CreatedAt.After(shaTags[j].CreatedAt)
	})

	return shaTags[0].Tag, nil
}
