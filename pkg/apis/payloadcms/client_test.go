package payloadcms_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ainsleydev/webkit/pkg/apis/payloadcms"
)

type Products struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func TestClient_PerformRequest(t *testing.T) {
	c := payloadcms.CollectServiceOp{Client: &payloadcms.Client{}}

	r := payloadcms.ListResponse[Products]{}
	err := c.Find(context.TODO(), "collection", &r)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(r)
}
