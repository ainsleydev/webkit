package payloadcms

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/goccy/go-json"

	"github.com/ainsleydev/webkit/pkg/util/httputil"
)

// Client represents a Payload CMS client.
// For more information, see https://payloadcms.com/docs/api.
type Client struct {
	// Each collection is mounted using its slug value. For example, if a collection's slug is
	// users, all corresponding routes will be mounted on /api/users.
	// For more info, visit: https://payloadcms.com/docs/rest-api/overview#collections
	Collections CollectionService
	// Globals cannot be created or deleted, so there are only two REST endpoints opened:
	// For more info, visit: https://payloadcms.com/docs/rest-api/overview#globals
	Globals GlobalsService
	// Media is a separate service used to upload and manage media files.
	// For more info, visit: https://payloadcms.com/docs/upload/overview
	Media MediaService

	// Private fields
	client  *http.Client
	baseURL string
	apiKey  string
}

// Response is a PayloadAPI API response. This wraps the standard http.Response
// returned from Payload and provides convenient access to things like
// body bytes.
type Response struct {
	*http.Response
	Content []byte
	Message string
	Errors  []Error
}

type Error struct {
	Message string
}

func (e Error) Error() string {
	return e.Message
}

// New creates a new Payload CMS client.
func New(baseURL, apiKey string) *Client {
	c := &Client{
		client:  http.DefaultClient,
		baseURL: baseURL,
		apiKey:  apiKey,
	}
	c.Collections = CollectionServiceOp{Client: c}
	c.Globals = GlobalsService{Client: c}
	c.Media = MediaService{Client: c}
	return c
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.
//
// Errors occur in the eventuality if the http.StatusCode is not 2xx.
func (c Client) Do(ctx context.Context, method, path string, body io.Reader, v any) (Response, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, strings.TrimPrefix(path, "/"))
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "users API-Key "+c.apiKey)

	r, err := c.performRequest(req)
	if err != nil {
		return Response{}, err
	}

	if v == nil {
		return r, nil
	}

	return r, json.Unmarshal(r.Content, v)
}

// DoWithRequest sends an API request using the provided http.Request and returns the API response.
// The API response is JSON decoded and stored in the value pointed to by v, or returned
// as an error if an API error has occurred.
func (c Client) DoWithRequest(_ context.Context, req *http.Request, v any) (Response, error) {
	r, err := c.performRequest(req)
	if err != nil {
		return Response{}, err
	}
	if v == nil {
		return r, nil
	}
	return r, json.Unmarshal(r.Content, v)
}

// Get sends an HTTP GET request and returns the API response.
// The API response is JSON decoded and stored in the value pointed to by v, or returned
// as an error if an API error has occurred.
func (c Client) Get(ctx context.Context, path string, v any) (Response, error) {
	return c.Do(ctx, http.MethodGet, path, nil, v)
}

// Put sends an HTTP PUT request and returns the API response.
// The API response is JSON decoded and stored in the value pointed to by v, or returned
// as an error if an API error has occurred.
func (c Client) Put(ctx context.Context, path string, in any) (Response, error) {
	buf, err := json.Marshal(in)
	if err != nil {
		return Response{}, err
	}
	return c.Do(ctx, http.MethodPut, path, bytes.NewReader(buf), nil)
}

// Post sends an HTTP POST request and returns the API response.
// The API response is JSON decoded and stored in the value pointed to by v, or returned
// as an error if an API error has occurred.
func (c Client) Post(ctx context.Context, path string, in any) (Response, error) {
	buf, err := json.Marshal(in)
	if err != nil {
		return Response{}, err
	}
	return c.Do(ctx, http.MethodPost, path, bytes.NewReader(buf), nil)
}

// Delete sends an HTTP DELETE request and returns the API response.
// The API response is JSON decoded and stored in the value pointed to by v, or returned
// as an error if an API error has occurred.
func (c Client) Delete(ctx context.Context, path string, v any) (Response, error) {
	return c.Do(ctx, http.MethodDelete, path, nil, v)
}

func (c Client) newRequest(ctx context.Context, path, method string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, strings.TrimPrefix(path, "/"))
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "users API-Key "+c.apiKey)

	return req, nil
}

func (c Client) newFormRequest(ctx context.Context, method, path string, body io.Reader, contentType string) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, strings.TrimPrefix(path, "/"))
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// Set the content type to contain the boundary.
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "users API-Key "+c.apiKey)

	return req, nil
}

func (c Client) performRequest(req *http.Request) (Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	r := Response{Response: resp}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	r.Content = buf

	if !httputil.Is2xx(resp.StatusCode) {
		return r, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, buf)
	}

	return r, nil
}
