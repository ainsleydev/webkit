package payloadcms

import (
	"context"
	"net/http"

	"github.com/goccy/go-json"
)

const collectsBasePath = "collects"

// CollectService is an interface for interfacing with the collect endpoints
// of the Shopify API.
// See: https://help.shopify.com/api/reference/products/collect
//type CollectService interface {
//	List(context.Context, interface{}) ([]Collect, error)
//	Count(context.Context, interface{}) (int, error)
//	Get(context.Context, uint64, interface{}) (*Collect, error)
//	Create(context.Context, Collect) (*Collect, error)
//	Delete(context.Context, uint64) error
//}

// CollectServiceOp handles communication with the collect related methods of
// the Shopify API.
type CollectServiceOp struct {
	Client *Client
}

type ListResponse[T any] struct {
	Docs []T `json:"products"`
	//Id        string    `json:"id"`
	//Title     string    `json:"title"`
	//Content   string    `json:"content"`
	//Slug      string    `json:"slug"`
	//CreatedAt time.Time `json:"createdAt"`
	//UpdatedAt time.Time `json:"updatedAt"`

	Total int `json:"total"`
	//TotalDocs     int  `json:"totalDocs"`
	//Limit         int  `json:"limit"`
	//TotalPages    int  `json:"totalPages"`
	//Page          int  `json:"page"`
	//PagingCounter int  `json:"pagingCounter"`
	//HasPrevPage   bool `json:"hasPrevPage"`
	//HasNextPage   bool `json:"hasNextPage"`
	//PrevPage      any  `json:"prevPage"`
	//NextPage      any  `json:"nextPage"`
}

// https://vladimir.varank.in/notes/2022/05/a-real-life-use-case-for-generics-in-go-api-for-client-side-pagination/

func (c *CollectServiceOp) Find(ctx context.Context, collection string, out ListResponse[any]) (ListResponse[T any], error) {
	//path := `/api/` + collection
	buf, err := c.Client.PerformRequest(ctx, http.MethodGet, "", nil)
	if err != nil {
		return err
	}
	r := &Response[T]{}
	if err = json.Unmarshal(buf, r); err != nil {
		return nil, err
	}
	return r, nil
}

func unmarshal[T any](data []byte) (*T, error) {
	out := new(T)
	if err := json.Unmarshal(data, out); err != nil {
		return nil, err
	}
	return out, nil
}

// List collects
//func (s *CollectServiceOp) List(ctx context.Context, options interface{}) ([]Collect, error) {
//	path := fmt.Sprintf("%s.json", collectsBasePath)
//	resource := new(CollectsResource)
//	err := s.client.Get(ctx, path, resource, options)
//	return resource.Collects, err
//}
