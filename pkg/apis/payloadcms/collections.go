package payloadcms

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/google/go-querystring/query"
)

// CollectionService handles communication with the collect related methods of
// the Shopify API.
type CollectionService struct {
	Client *Client
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

type (
	// ListParams represents additional query parameters for the find endpoint.
	ListParams struct {
		Sort  string         `json:"sort" url:"sort"`   // Sort the returned documents by a specific field.
		Where map[string]any `json:"where" url:"where"` // Constrain returned documents with a where query.
		Limit int            `json:"limit" url:"limit"` // Limit the returned documents to a certain number.
		Page  int            `json:"page" url:"page"`   // Get a specific page of documents.
	}
	// ListResponse represents a list of entities that is sent back
	// from the Payload CMS.
	ListResponse[T any] struct {
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
	// CreateResponse represents a response from the Payload CMS
	// when a new entity is created.
	CreateResponse[T any] struct {
		Doc     T      `json:"doc"`
		Message string `json:"message"`
		Errors  []any  `json:"errors"`
	}
	// UpdateResponse represents a response from the Payload CMS
	// when an entity is updated.
	UpdateResponse[T any] struct {
		Doc     T      `json:"doc"`
		Message string `json:"message"`
		Errors  []any  `json:"error"`
	}
)

// FindById finds a collection entity by its ID.
func (s CollectionService) FindById(ctx context.Context, collection Collection, id int, out any) (Response, error) {
	path := fmt.Sprintf("/api/%s/%d", collection, id)
	return s.Client.Do(ctx, http.MethodGet, path, nil, out)
}

// FindBySlug finds a collection entity by its slug.
func (s CollectionService) FindBySlug(ctx context.Context, collection Collection, slug string, out any) (Response, error) {
	path := fmt.Sprintf("/api/%s/%s", collection, slug)
	return s.Client.Do(ctx, http.MethodGet, path, nil, out)
}

// List lists all collection entities.
func (s CollectionService) List(ctx context.Context, collection Collection, params ListParams, out any) (Response, error) {
	v, err := query.Values(params)
	if err != nil {
		return Response{}, err
	}
	path := fmt.Sprintf("/api/%s?%s", collection, v)
	return s.Client.Do(ctx, http.MethodGet, path, nil, out)
}

// Create creates a new collection entity.
func (s CollectionService) Create(ctx context.Context, collection Collection, in any) (Response, error) {
	path := fmt.Sprintf("/api/%s", collection)
	buf, err := json.Marshal(in)
	if err != nil {
		return Response{}, err
	}
	return s.Client.Do(ctx, http.MethodPost, path, bytes.NewReader(buf), nil)
}

// UpdateByID updates a collection entity by its ID.
func (s CollectionService) UpdateByID(ctx context.Context, collection Collection, id int, in any) (Response, error) {
	path := fmt.Sprintf("/api/%s/%d", collection, id)
	buf, err := json.Marshal(in)
	if err != nil {
		return Response{}, err
	}
	return s.Client.Do(ctx, http.MethodPut, path, bytes.NewReader(buf), nil)
}

// DeleteByID deletes a collection entity by its ID.
func (s CollectionService) DeleteByID(ctx context.Context, collection Collection, id int) (Response, error) {
	path := fmt.Sprintf("/api/%s/%d", collection, id)
	return s.Client.Do(ctx, http.MethodDelete, path, nil, nil)
}
