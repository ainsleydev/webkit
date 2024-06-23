package schemaorg

// BreadcrumbList is an ItemList consisting of a chain of linked Web pages, typically
// described using at least their URL and their name, and usually ending
// with the current page.
//
// See:
// - https://schema.org/BreadcrumbList
// - https://developers.google.com/search/docs/data-types/breadcrumb
//
//	{
//	   "@context": "https://schema.org",
//	   "@type": "BreadcrumbList",
//	   "itemListElement": [
//	       {
//	           "@type": "ListItem",
//	           "position": 1,
//	           "name": "Books",
//	           "item": "https://example.com/books"
//	       },
//	       {
//	           "@type": "ListItem",
//	           "position": 2,
//	           "name": "Science Fiction",
//	           "item": "https://example.com/books/sciencefiction"
//	       }
//	   ]
//	}
type BreadcrumbList []BreadcrumbItem

// BreadcrumbItem contains details about an individual item in the list.
// Google only supports item, name and position.
//
// See: https://schema.org/ListItem.
type BreadcrumbItem struct {
	// The position of the breadcrumb in the breadcrumb trail.
	// Position 1 signifies the beginning of the trail.
	Position int `json:"position"`

	// The title of the breadcrumb displayed for the user.
	Name string `json:"name"`

	// The URL to the webpage that represents the breadcrumb.
	Item string `json:"item,omitempty"`
}

// Alias the types for JSON-LD.
type (
	breadcrumbs struct {
		Context string                  `json:"@context"`
		Type    string                  `json:"@type"`
		Items   []breadcrumbsItemListEl `json:"itemListElement"`
	}
	breadcrumbsItemListEl struct {
		Type     string `json:"@type"`
		Position int    `json:"position"`
		Name     string `json:"name"`
		Item     string `json:"item,omitempty"`
	}
)

// MarshalJSON implements the json.Marshaler interface to generate
// the JSON-LD for the BreadcrumbList.
func (s BreadcrumbList) MarshalJSON() ([]byte, error) {
	b := breadcrumbs{
		Context: Context,
		Type:    "BreadcrumbList",
		Items:   make([]breadcrumbsItemListEl, len(s)),
	}
	for i, bi := range s {
		b.Items[i] = breadcrumbsItemListEl{
			Type:     "ListItem",
			Position: bi.Position,
			Name:     bi.Name,
			Item:     bi.Item,
		}
	}
	return marshal(b)
}
