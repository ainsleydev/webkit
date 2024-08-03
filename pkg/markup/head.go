package markup

import (
	"context"
	"github.com/ainsleydev/webkit/pkg/markup/internal/templates"
	"io"
	"time"

	webkitctx "github.com/ainsleydev/webkit/pkg/context"
	schemaorg "github.com/ainsleydev/webkit/pkg/markup/schema"
)

// HeadProps defines the properties that should be included in the
// head of the document.
type HeadProps struct {
	// Required meta properties
	Title         string
	Description   string
	Image         string
	DatePublished time.Time
	DateModified  time.Time
	Locale        string // ISO 639-1 defaults to "en_GB"

	// Other
	URL    string // The full URL of the page.
	IsPage bool   // If true, the page is a page, not a post.

	// Additional meta properties
	Private   bool   // If true, the page should not be indexed by search engines.
	Canonical string // The canonical URL of the page.

	// Schema, Meta & Opengraph
	OpenGraph    *OpenGraph
	Twitter      *TwitterCard
	Organisation *schemaorg.Organisation
	Navigation   *schemaorg.BreadcrumbList

	// Other (Code Injection)
	Other string

	// To define additional meta tags and any other HTML, see webkitctx.WithHeadSnippet
	Snippets []webkitctx.MarkupSnippet
}

// Render renders the head of the document to the provided writer.
func (h HeadProps) Render(ctx context.Context, w io.Writer) error {
	s, ok := webkitctx.HeadSnippets(ctx)
	if ok {
		h.Snippets = s
	}
	return templates.Render(ctx, w, "head.html", h)
}
