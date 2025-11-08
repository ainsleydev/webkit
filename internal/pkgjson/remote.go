package pkgjson

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// FetchFromRemote fetches and parses a package.json from a remote URL.
func FetchFromRemote(ctx context.Context, url string) (*PackageJSON, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "fetching remote package.json")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading response body")
	}

	pkg := &PackageJSON{
		Dependencies:     make(map[string]string),
		DevDependencies:  make(map[string]string),
		PeerDependencies: make(map[string]string),
		raw:              make(map[string]any),
	}

	// Unmarshal into struct (known fields)
	if err = json.Unmarshal(data, pkg); err != nil {
		return nil, errors.Wrap(err, "parsing package.json struct")
	}

	// Unmarshal into map to capture all fields (including unknown)
	if err = json.Unmarshal(data, &pkg.raw); err != nil {
		return nil, errors.Wrap(err, "parsing package.json raw")
	}

	return pkg, nil
}
