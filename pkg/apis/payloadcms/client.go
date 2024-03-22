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

type Response[T any] struct {
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
