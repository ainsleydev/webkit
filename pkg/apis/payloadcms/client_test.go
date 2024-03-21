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
	find, err := payloadcms.Find[Products](context.TODO(), nil, "collection")
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(find)
}
