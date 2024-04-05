package payloadcms

import (
	"context"
	"io"
	"net/http"
)

type (
	Client struct {
		client  *http.Client
		baseURL string
	}
)

func (c *Client) PerformRequest(ctx context.Context, method, url string, body any) ([]byte, error) {
	//kkk, err := json.Marshal(body)
	//if err != nil {
	//	return err
	//}
	//
	req, err := http.NewRequest(method, "https://dummyjson.com/products", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil

}

func (c *Client) doGetRequest() {

}
