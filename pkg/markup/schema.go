package markup

import "encoding/json"

// SchemaOrgOrganisation represents a structured data definition for an organisation
// according to schema.org. This can be used to provide details about the
// organisation and improve search engine understanding.
//
// See: https://schema.org/Organization
type SchemaOrgOrganisation struct {
	Context     string                       `json:"@context"`    // Always "https://schema.org"
	Type        string                       `json:"@type"`       // Always "Organization"
	ID          string                       `json:"@id"`         // Full URL
	URL         string                       `json:"url"`         // Full URL
	LegalName   string                       `json:"legalName"`   // The legal name of the organisation
	Description string                       `json:"description"` // A description of the organisation, can be the same as the tagline.
	Logo        string                       `json:"logo"`        // Full URL, no SVGs
	SameAs      []string                     `json:"sameAs"`      // An array of full social media URLs
	Address     SchemaOrgOrganisationAddress `json:"address"`
}

// MarshalJSON is a custom JSON marshaller for the SchemaOrgOrganisation struct.
// It sets the context and type to the correct values before marshalling.
func (o *SchemaOrgOrganisation) MarshalJSON() ([]byte, error) {
	type Alias SchemaOrgOrganisation // Alias to prevent stack overflow
	alias := (*Alias)(o)
	alias.Context = "https://schema.org"
	alias.Type = "Organization"
	alias.Address.Type = "PostalAddress"
	return json.Marshal(alias)
}

// SchemaOrgOrganisationAddress represents a structured data definition for the
// physical or mailing address of an organization according to schema.org.
//
// See: https://schema.org/PostalAddress
type SchemaOrgOrganisationAddress struct {
	Type            string `json:"@type"`           // Always "PostalAddress"
	StreetAddress   string `json:"streetAddress"`   // I.e ainsley.dev, 71-75 Shelton Street, Covent Garden, London, WC2H 9JQ
	AddressLocality string `json:"addressLocality"` // I.e London
	AddressRegion   string `json:"addressRegion"`   // I.e Greater London
	AddressCountry  string `json:"addressCountry"`  // I.e UK
	PostalCode      string `json:"postalCode"`      // I.e WC2H 9JQ
}

// SchemaOrgNavItemList defines a structured data representation for a navigational list
// of items on a webpage. This helps search engines understand the website's
// structure and potentially improve search ranking.
//
// See: https://schema.org/WebPage
type SchemaOrgNavItemList struct {
	Context         string                     `json:"@context"`        // Always "https://schema.org"
	Type            string                     `json:"@type"`           // Always "ItemList"
	ItemListElement []SchemaOrgItemListElement `json:"itemListElement"` // The list of items
}

// MarshalJSON is a custom JSON marshaller for the SchemaOrgItemList struct.
// It sets the context and type to the correct values before marshalling.
func (o *SchemaOrgNavItemList) MarshalJSON() ([]byte, error) {
	type Alias SchemaOrgNavItemList // Alias to prevent stack overflow
	alias := (*Alias)(o)
	alias.Context = "https://schema.org"
	alias.Type = "ItemList"
	return json.Marshal(alias)
}

// SchemaOrgItemListElement represents a single item within a navigational
// list on a webpage.
//
// See: https://schema.org/ItemList
type SchemaOrgItemListElement struct {
	Type        string `json:"@type"`       // Always "ListItem"
	Position    int    `json:"position"`    // I.e 1, 2, 3
	Name        string `json:"name"`        // I.e "Home"
	Description string `json:"description"` // I.e "The homepage of the website" usually the same as the description tag.
	URL         string `json:"url"`         // Full URL
}

// SchemaOrgFAQPage defines an FAQPage is a WebPage presenting one or more
// "Frequently asked questions" (see also QAPage).
//
// Example as defined on: https://developers.google.com/search/docs/data-types/faqpage
//
//	{
//	   "@context": "https://schema.org",
//	   "@type": "FAQPage",
//	   "mainEntity": [
//	       {
//	           "@type": "Question",
//	           "name": "How to find an apprenticeship?",
//	           "acceptedAnswer": {
//	               "@type": "Answer",
//	               "text": "<p>We provide an official service to search through available apprenticeships. To get started, create an account here, specify the desired region, and your preferences. You will be able to search through all officially registered open apprenticeships.</p>"
//	           }
//	       },
//	       {
//	           "@type": "Question",
//	           "name": "Whom to contact?",
//	           "acceptedAnswer": {
//	               "@type": "Answer",
//	               "text": "You can contact the apprenticeship office through our official phone hotline above, or with the web-form below. We generally respond to written requests within 7-10 days."
//	           }
//	       }
//	   ]
//	}
//
// See: https://schema.org/FAQPage
type SchemaOrgFAQPage []SchemaOrgQuestionAnswer

// SchemaOrgQuestionAnswer represents a single question and answer within an FAQPage.
type SchemaOrgQuestionAnswer struct {
	// The Question type defines a single answered question within the FAQ.
	// Every Question instance must be contained within the mainEntity
	// property array of the schema.org/FAQPage.
	Question string `json:"-"`

	// The full answer to the question. The answer may contain HTML content
	// such as links and lists.
	//
	// Google Search displays the following HTML tags; all other tags are ignored:
	// <h1> through <h6>, <br>, <ol>, <ul>, <li>, <a>, <p>, <div>, <b>, <strong>, <i>, and <em>.
	Answer string `json:"-"`
}

// MarshalJSON is a custom JSON marshaller for the SchemaOrgFAQPage struct.
func (s SchemaOrgFAQPage) MarshalJSON() ([]byte, error) {
	type mainEntity struct {
		Type           string `json:"@type"`
		Name           string `json:"name"`
		AcceptedAnswer struct {
			Type string `json:"@type"`
			Text string `json:"text"`
		} `json:"acceptedAnswer"`
	}

	var entities []mainEntity
	for _, qa := range s {
		entities = append(entities, mainEntity{
			Type: "Question",
			Name: qa.Question,
			AcceptedAnswer: struct {
				Type string `json:"@type"`
				Text string `json:"text"`
			}{
				Type: "Answer",
				Text: qa.Answer,
			},
		})
	}

	marshal, err := json.MarshalIndent(entities, "", "\t")
	if err != nil {
		return nil, err
	}

	b := `{
	"@context": "https://schema.org",
	"@type": "FAQPage",
	"mainEntity":` + string(marshal) + `
}`

	return []byte(b), nil
}

// SchemaOrgBreadcrumbList is an ItemList consisting of a chain of linked Web pages,
// typically described using at least their URL and their name, and usually ending
// with the current page.
//
// Example as defined on: https://developers.google.com/search/docs/data-types/breadcrumb
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
//
// See: https://schema.org/BreadcrumbList
type SchemaOrgBreadcrumbList struct {
	Context string                    `json:"@context"`
	Type    string                    `json:"@type"`
	Items   []SchemaOrgBreadcrumbItem `json:"itemListElement"`
}

// SchemaOrgBreadcrumbItem contains details about an individual item in the list.
// Google only supports item, name and position.
//
// See: https://schema.org/ListItem. wing:
type SchemaOrgBreadcrumbItem struct {
	// Always "ListItem"
	Type string `json:"@type"`

	// The position of the breadcrumb in the breadcrumb trail.
	// Position 1 signifies the beginning of the trail.
	Position int `json:"position"`

	// The title of the breadcrumb displayed for the user.
	Name string `json:"name"`

	// The URL to the webpage that represents the breadcrumb.
	Item string `json:"item,omitempty"`
}
