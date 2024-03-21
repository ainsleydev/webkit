package payloadcms

import (
	"net/http"
)

type (
	Client struct {
		client  *http.Client
		baseURL string
	}
)

func (c *Client) PerformRequest(method, url string, body any) ([]byte, error) {
	//kkk, err := json.Marshal(body)
	//if err != nil {
	//	return err
	//}
	//
	//req, err := http.NewRequest(method, url, )

	return nil, nil

}

type Response[T any] struct {
	Docs []T `json:"docs"`
	//Id        string    `json:"id"`
	//Title     string    `json:"title"`
	//Content   string    `json:"content"`
	//Slug      string    `json:"slug"`
	//CreatedAt time.Time `json:"createdAt"`
	//UpdatedAt time.Time `json:"updatedAt"`
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
