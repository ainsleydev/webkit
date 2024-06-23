package webkitctx

import (
	"context"
)

// ContextKey defines the type used to set a context key.
type ContextKey string

const (
	// ContextKeyRequestID is the key used to define a unique
	// request identifier.
	ContextKeyRequestID ContextKey = "request_id"

	// ContextKeyHeadTemplates is the key used to define the head
	// templates for a request.
	ContextKeyHeadTemplates ContextKey = "head_templates"

	// ContextKeyFooterTemplates is the key used to define the head
	// templates for a request.
	ContextKeyFooterTemplates ContextKey = "footer_templates"
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

// MarkupContent defines the structure of a piece of HTML content that
// is to be injected into the head or body of a document.
type MarkupContent struct {
	Name    string
	Content string
}

// WithHTMLHeadTemplate sets an HTML snippet to be injected into
// the <head> of the document.
func WithHTMLHeadTemplate(ctx context.Context, name, content string) context.Context {
	c, ok := ctx.Value(ContextKeyHeadTemplates).([]MarkupContent)
	if !ok {
		c = make([]MarkupContent, 0)
	}
	c = append(c, MarkupContent{
		Name:    name,
		Content: content,
	})
	return context.WithValue(ctx, ContextKeyHeadTemplates, c)
}

// HeadContent returns the head content templates from the context.
func HeadContent(ctx context.Context) ([]MarkupContent, bool) {
	c, ok := ctx.Value(ContextKeyHeadTemplates).([]MarkupContent)
	return c, ok
}

// WithHTMLFooterTemplate sets an HTML snippet to be injected into
// the footer of the document.
func WithHTMLFooterTemplate(ctx context.Context, name string, content string) context.Context {
	c, ok := ctx.Value(ContextKeyFooterTemplates).([]MarkupContent)
	if !ok {
		c = make([]MarkupContent, 0)
	}
	c = append(c, MarkupContent{
		Name:    name,
		Content: content,
	})
	return context.WithValue(ctx, ContextKeyFooterTemplates, c)
}

// FooterContent returns the footer content templates from the context.
func FooterContent(ctx context.Context) ([]MarkupContent, bool) {
	c, ok := ctx.Value(ContextKeyFooterTemplates).([]MarkupContent)
	return c, ok
}
