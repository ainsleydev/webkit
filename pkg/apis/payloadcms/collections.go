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
	client *Client
}

// https://vladimir.varank.in/notes/2022/05/a-real-life-use-case-for-generics-in-go-api-for-client-side-pagination/

func Find[T any](ctx context.Context, client *Client, collection string) (*Response[T], error) {
	//path := `/api/` + collection

	resource := new(Response[T])
	buf, err := client.PerformRequest(ctx, http.MethodGet, "", nil)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(buf, resource); err != nil {
		return nil, err
	}

	return resource, nil
}

// List collects
//func (s *CollectServiceOp) List(ctx context.Context, options interface{}) ([]Collect, error) {
//	path := fmt.Sprintf("%s.json", collectsBasePath)
//	resource := new(CollectsResource)
//	err := s.client.Get(ctx, path, resource, options)
//	return resource.Collects, err
//}
