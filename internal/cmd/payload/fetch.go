package payload

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ainsleydev/webkit/internal/ghapi"
)

const (
	// packageJSONPath is the path to Payload's package.json in the repository.
	packageJSONPath = "package.json"
)

// PayloadDependencies contains all dependencies from Payload's package.json.
// This is used to determine which dependencies should be bumped when updating.
type PayloadDependencies struct {
	Dependencies    map[string]string
	DevDependencies map[string]string
	AllDeps         map[string]string // Combined for easier lookup
}

// FetchPayloadDependencies fetches Payload CMS's package.json from GitHub
// and extracts all its dependencies.
//
// This allows us to bump ALL dependencies that Payload uses, not just
// payload and @payloadcms/* packages. For example, if Payload depends on
// lexical@0.28.0, we can update the user's lexical to match.
func FetchPayloadDependencies(ctx context.Context, client ghapi.Client, version string) (*PayloadDependencies, error) {
	// Fetch package.json content from GitHub.
	// Use the version tag (with v prefix) as the ref.
	ref := "v" + version
	content, err := client.GetFileContent(ctx, payloadOwner, payloadRepo, packageJSONPath, ref)
	if err != nil {
		return nil, errors.Wrap(err, "fetching package.json from GitHub")
	}

	if content == nil {
		return nil, errors.New("package.json content is nil")
	}

	// Parse the package.json.
	var pkgJSON struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(content, &pkgJSON); err != nil {
		return nil, errors.Wrap(err, "parsing package.json")
	}

	// Combine all dependencies for easier lookup.
	allDeps := make(map[string]string)
	for name, version := range pkgJSON.Dependencies {
		allDeps[name] = version
	}
	for name, version := range pkgJSON.DevDependencies {
		allDeps[name] = version
	}

	return &PayloadDependencies{
		Dependencies:    pkgJSON.Dependencies,
		DevDependencies: pkgJSON.DevDependencies,
		AllDeps:         allDeps,
	}, nil
}
