package webkitctx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	ctx := context.Background()
	rid := "12345"

	ctx = WithRequestID(ctx, rid)

	got, ok := RequestID(ctx)
	assert.True(t, ok)
	assert.Equal(t, rid, got)
}

func TestHeadTemplate(t *testing.T) {
	ctx := context.Background()
	name := "test"
	content := "<script>alert('test')</script>"

	ctx = WithHeadSnippet(ctx, name, content)

	got, ok := HeadSnippets(ctx)
	assert.True(t, ok)
	assert.Len(t, got, 1)
	assert.Equal(t, name, got[0].Name)
	assert.Equal(t, content, string(got[0].Content))
}

func TestFooterTemplate(t *testing.T) {
	ctx := context.Background()
	name := "test"
	content := "<script>alert('test')</script>"

	ctx = WithFooterSnippet(ctx, name, content)

	got, ok := FooterSnippets(ctx)
	assert.True(t, ok)
	assert.Len(t, got, 1)
	assert.Equal(t, name, got[0].Name)
	assert.Equal(t, content, string(got[0].Content))
}
