package wordpress

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ainsleydev/webkit/pkg/util/httputil"
)

// Client is a WordPress API client.
type Client struct {
	client  *http.Client
	baseURL string
}

// New creates a new WordPress API client.
func New(baseURL string) *Client {
	return &Client{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

// Get sends a GET request to the specified WordPress URL and returns the response body.
// The base URL is prepended to the URL, for example:
// https://wordpress/wp-json/wp/v2/posts/21
func (c *Client) Get(url string) ([]byte, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/%s", c.baseURL, strings.TrimLeft(url, "/")))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if !httputil.Is2xx(resp.StatusCode) {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	return body, nil
}

// GetAndUnmarshal performs an HTTP GET request to the specified WordPress URL,
// unmarshal the response body into the provided struct type, and returns any error.
func (c *Client) GetAndUnmarshal(url string, v any) error {
	body, err := c.Get(url)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}
