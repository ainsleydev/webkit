package payloadcms

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/goccy/go-json"
)

// Client represents a Payload CMS client.
// For more information, see https://payloadcms.com/docs/api.
type Client struct {
	client  *http.Client
	baseURL string
	apiKey  string

	Collections CollectionService
	Media       MediaService
}

// New creates a new Payload CMS client.
func New(baseURL, apiKey string) *Client {
	c := &Client{
		client:  http.DefaultClient,
		baseURL: baseURL,
		apiKey:  apiKey,
	}
	c.Collections = CollectionService{Client: c}
	c.Media = MediaService{Client: c}
	return c
}

// Do sends an HTTP request and returns the response body as a byte slice.
// It returns an error if the request fails or if the response status code is not in the 2xx range.
func (c *Client) Do(ctx context.Context, method, path string, body io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, strings.TrimPrefix(path, "/"))
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "users API-Key "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, buf)
	}

	return buf, nil
}

// DoAndUnmarshal sends an HTTP request and unmarshal the response body into the given value.
func (c *Client) DoAndUnmarshal(ctx context.Context, method, url string, body io.Reader, out any) error {
	buf, err := c.Do(ctx, method, url, body)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, out)
}
