package payloadcms

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
)

type GlobalsService struct {
	Client *Client
}

type Global string

func (s GlobalsService) Get(ctx context.Context, global Global, in any) (Response, error) {
	path := fmt.Sprintf("/api/globals/%s", global)
	return s.Client.Get(ctx, path, in)
}

func (s GlobalsService) Update(ctx context.Context, global Global, in any) (Response, error) {
	path := fmt.Sprintf("/api/globals/%s", global)
	buf, err := json.Marshal(in)
	if err != nil {
		return Response{}, err
	}
	return s.Client.Do(ctx, http.MethodPost, path, bytes.NewReader(buf), nil)
}
