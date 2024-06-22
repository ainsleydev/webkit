package markup

import (
	_ "embed"
	"time"
)

//go:embed head.html
var headTpl string

// HeadTemplate is the template for the head of the HTML document.
// It requires a HeadProps struct to be passed in when executing the template.
//var HeadTemplate = template.Must(template.New("").ParseFiles(
//	"./head.html",
//	"./opengraph.html",
//))

// HeadProps defines the properties that should be included in the
// head of the document.
type HeadProps struct {
	// Required meta properties
	Title         string
	Description   string
	Image         string // TODO ???
	DatePublished time.Time
	DateModified  time.Time
	Locale        string // ISO 639-1 defaults to "en_GB"

	// Other
	URL    string // The full URL of the page.
	IsPage bool   // If true, the page is a page, not a post.

	// Resources
	Scripts []string
	Styles  []string
	Fonts   []struct { // Remove inline
		Link string
		Type string
	}
	Hash int64 // Unix now

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
