package markup

import (
	"context"
	"embed"
	_ "embed"
	"html/template"
	"io"
	"time"
)

//go:embed *.html
var templatesFS embed.FS

// headTemplate is the template for the head of the HTML document.
// It requires a HeadProps struct to be passed in when executing the template.
var headTemplate = template.Must(template.ParseFS(templatesFS,
	"head.html",
	"opengraph.html",
	"twitter.html",
))

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
	OpenGraph  *OpenGraph
	Twitter    *TwitterCard
	Org        *SchemaOrgOrganisation
	Navigation *SchemaOrgItemList

	// Other (Code Injection)
	Other string
}

// Render renders the head of the document to the provided writer.
func (h HeadProps) Render(_ context.Context, w io.Writer) error {
	return headTemplate.ExecuteTemplate(w, "head.html", h)
}
