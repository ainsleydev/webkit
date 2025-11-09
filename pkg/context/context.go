package webkitctx

import (
	"context"
)

// ContextKey is a typed string for context keys to avoid collisions.
type ContextKey string

const (
	// ContextKeyRequestID is the key used to define a unique
	// request identifier.
	ContextKeyRequestID ContextKey = "request_id"

	// ContextKeyURL is the key used to retrieve the full URL in
	// the context.
	ContextKeyURL = "url"

	// ContextKeyHeadSnippets is the key used to define the head
	// templates for a request.
	ContextKeyHeadSnippets ContextKey = "head_snippets"

	// ContextKeyFooterSnippets is the key used to define the head
	// templates for a request.
	ContextKeyFooterSnippets ContextKey = "footer_snippets"
)

// RequestID extracts a unique request identifier from a context
func RequestID(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(ContextKeyRequestID).(string)
	return rid, ok
}

// WithRequestID returns a new context with the given request identifier.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKeyRequestID, requestID)
}

// URL extracts the full URL from a context.
func URL(ctx context.Context) (string, bool) {
	url, ok := ctx.Value(ContextKeyURL).(string)
	return url, ok
}

// WithURL returns a new context with the given full URL.
func WithURL(ctx context.Context, url string) context.Context {
	return context.WithValue(ctx, ContextKeyURL, url) //nolint
}

// MarkupSnippet represents HTML content to be injected into a document.
type MarkupSnippet struct {
	Name    string
	Content string
}

// WithHeadSnippet sets an HTML snippet to be injected into
// the <head> of the document.
func WithHeadSnippet(ctx context.Context, name, content string) context.Context {
	c, ok := ctx.Value(ContextKeyHeadSnippets).([]MarkupSnippet)
	if !ok {
		c = make([]MarkupSnippet, 0)
	}
	c = append(c, MarkupSnippet{
		Name:    name,
		Content: content,
	})
	return context.WithValue(ctx, ContextKeyHeadSnippets, c)
}

// HeadSnippets returns the head content templates from the context.
func HeadSnippets(ctx context.Context) ([]MarkupSnippet, bool) {
	c, ok := ctx.Value(ContextKeyHeadSnippets).([]MarkupSnippet)
	return c, ok
}

// WithFooterSnippet sets an HTML snippet to be injected into
// the footer of the document.
func WithFooterSnippet(ctx context.Context, name string, content string) context.Context {
	c, ok := ctx.Value(ContextKeyFooterSnippets).([]MarkupSnippet)
	if !ok {
		c = make([]MarkupSnippet, 0)
	}
	c = append(c, MarkupSnippet{
		Name:    name,
		Content: content,
	})
	return context.WithValue(ctx, ContextKeyFooterSnippets, c)
}

// FooterSnippets returns the footer content templates from the context.
func FooterSnippets(ctx context.Context) ([]MarkupSnippet, bool) {
	c, ok := ctx.Value(ContextKeyFooterSnippets).([]MarkupSnippet)
	return c, ok
}
