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

	// GetLatestRelease returns the latest stable release tag for a repository.
	// Excludes draft and pre-release versions.
	GetLatestRelease(ctx context.Context, owner, repo string) (string, error)

	// GetFileContent fetches the content of a file from a repository at a specific ref (tag/branch/commit).
	// Returns the decoded file content as bytes.
	GetFileContent(ctx context.Context, owner, repo, path, ref string) ([]byte, error)
}

// DefaultClient implements the Client interface using the official go-github library.
type DefaultClient struct {
	client *github.Client
}

// New creates a new GitHub API client with authentication.
func New(token string) Client {
	enforce.NotEqual(token, "", "github token cannot be empty")

	return &DefaultClient{
		client: github.NewClient(nil).WithAuthToken(token),
	}
}

// NewWithoutAuth creates a new unauthenticated GitHub API client.
// Useful for accessing public repositories in tests or rate-limited scenarios.
func NewWithoutAuth() Client {
	return &DefaultClient{
		client: github.NewClient(nil),
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

// GetLatestRelease fetches the latest stable release for a repository.
// It excludes draft and pre-release versions.
//
// Returns an error if no stable releases are found or if the API request fails.
func (c *DefaultClient) GetLatestRelease(ctx context.Context, owner, repo string) (string, error) {
	// Fetch all releases for the repository.
	releases, _, err := c.client.Repositories.ListReleases(
		ctx,
		owner,
		repo,
		&github.ListOptions{PerPage: 50},
	)
	if err != nil {
		return "", err
	}

	if len(releases) == 0 {
		return "", errors.New("no releases found")
	}

	// Find the first stable release (not draft, not pre-release).
	for _, release := range releases {
		if !release.GetDraft() && !release.GetPrerelease() {
			tag := release.GetTagName()
			// Remove 'v' prefix if present.
			return strings.TrimPrefix(tag, "v"), nil
		}
	}

	return "", errors.New("no stable releases found")
}

// GetFileContent fetches a file's content from a GitHub repository.
// The ref parameter can be a tag, branch name, or commit SHA.
//
// Returns an error if the file doesn't exist or cannot be fetched.
func (c *DefaultClient) GetFileContent(ctx context.Context, owner, repo, path, ref string) ([]byte, error) {
	fileContent, _, _, err := c.client.Repositories.GetContents(
		ctx,
		owner,
		repo,
		path,
		&github.RepositoryContentGetOptions{
			Ref: ref,
		},
	)
	if err != nil {
		return nil, err
	}

	if fileContent == nil {
		return nil, errors.New("file content is nil")
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, err
	}

	return []byte(content), nil
}
