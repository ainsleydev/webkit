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
type CollectServiceOp[T any] struct {
	Client *Client
}

// https://vladimir.varank.in/notes/2022/05/a-real-life-use-case-for-generics-in-go-api-for-client-side-pagination/

func (c *CollectServiceOp[T]) Find(ctx context.Context, collection string) (*Response[T], error) {
	//path := `/api/` + collection
	buf, err := c.Client.PerformRequest(ctx, http.MethodGet, "", nil)
	if err != nil {
		return nil, err
	}
	r := &Response[T]{}
	if err = json.Unmarshal(buf, r); err != nil {
		return nil, err
	}
	return r, nil
}

func Find[T *Response[T]](c *CollectServiceOp, ctx context.Context, collection string) (*T, error) {
	buf, err := c.Client.PerformRequest(ctx, http.MethodGet, "", nil)
	if err != nil {
		return nil, err
	}
	out := new(T)
	if err = json.Unmarshal(buf, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *CollectServiceOpp[T]) FindTEMP(ctx context.Context, collection string, response *Response[T]) error {
	buf, err := c.Client.PerformRequest(ctx, http.MethodGet, "", nil)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(buf, response); err != nil {
		return err
	}
	return nil
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
