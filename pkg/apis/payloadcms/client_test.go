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
	c := payloadcms.CollectServiceOp[Products]{Client: &payloadcms.Client{}}

	got, err := c.Find(context.TODO(), "collection")

	fmt.Println(got.Docs[0].ID)

	r := payloadcms.Response[Products]{}
	err := c.Find(context.TODO(), "collection", &r)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(find)
}
