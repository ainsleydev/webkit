package payloadcms

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
)

// CollectionService handles communication with the collect related methods of
// the Shopify API.
type CollectionService struct {
	Client *Client
}

// ListResponse represents a list of entities that is sent back
// from the Payload CMS.
type ListResponse[T any] struct {
	Docs          []T  `json:"docs"`
	Total         int  `json:"total"`
	TotalDocs     int  `json:"totalDocs"`
	Limit         int  `json:"limit"`
	TotalPages    int  `json:"totalPages"`
	Page          int  `json:"page"`
	PagingCounter int  `json:"pagingCounter"`
	HasPrevPage   bool `json:"hasPrevPage"`
	HasNextPage   bool `json:"hasNextPage"`
	PrevPage      any  `json:"prevPage"`
	NextPage      any  `json:"nextPage"`
}

// Collection represents a collection slug from Payload.
// It's defined as a string under slug within the Collection Config.
type Collection string

const (
	// CollectionMedia defines the Payload media collection slug.
	CollectionMedia Collection = "media"
	// CollectionUsers defines the Payload users collection slug.
	CollectionUsers Collection = "users"
)

// FindById finds a collection entity by its ID.
func (c CollectionService) FindById(ctx context.Context, collection Collection, id int, out any) error {
	path := fmt.Sprintf("/api/%s/%d", collection, id)
	return c.Client.DoAndUnmarshal(ctx, http.MethodGet, path, nil, out)
}

// FindBySlug finds a collection entity by its slug.
func (c CollectionService) FindBySlug(ctx context.Context, collection Collection, slug string, out any) error {
	path := fmt.Sprintf("/api/%s/%s", collection, slug)
	return c.Client.DoAndUnmarshal(ctx, http.MethodGet, path, nil, out)
}

// List lists all collection entities.
func (c CollectionService) List(ctx context.Context, collection Collection, out any) error {
	path := fmt.Sprintf("/api/%s", collection)
	return c.Client.DoAndUnmarshal(ctx, http.MethodGet, path, nil, out)
}

// Create creates a new collection entity.
func (c CollectionService) Create(ctx context.Context, collection Collection, body any) error {
	path := fmt.Sprintf("/api/%s", collection)
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = c.Client.Do(ctx, http.MethodPost, path, bytes.NewReader(buf))
	return err
}

// UpdateByID updates a collection entity by its ID.
func (c CollectionService) UpdateByID(ctx context.Context, collection Collection, id int, body any) error {
	path := fmt.Sprintf("/api/%s/%d", collection, id)
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = c.Client.Do(ctx, http.MethodPut, path, bytes.NewReader(buf))
	return err
}

// DeleteByID deletes a collection entity by its ID.
func (c CollectionService) DeleteByID(ctx context.Context, collection Collection, id int) error {
	path := fmt.Sprintf("/api/%s/%d", collection, id)
	_, err := c.Client.Do(ctx, http.MethodDelete, path, nil)
	return err
}
