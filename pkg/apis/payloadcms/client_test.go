package payloadcms_test

import (
	"context"
	"testing"

	"github.com/ainsleydev/webkit/pkg/apis/payloadcms"
)

func TestClient_PerformRequest(t *testing.T) {
	c := payloadcms.CollectServiceOp{}

	resp, err := c.Find(context.TODO(), "collection")

	resp
}
