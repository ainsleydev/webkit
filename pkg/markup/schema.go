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

// SchemaOrgItemList defines a structured data representation for a navigational list
// of items on a webpage. This helps search engines understand the website's
// structure and potentially improve search ranking.
//
// See: https://schema.org/WebPage
type SchemaOrgItemList struct {
	Context         string                     `json:"@context"`        // Always "https://schema.org"
	Type            string                     `json:"@type"`           // Always "ItemList"
	ItemListElement []SchemaOrgItemListElement `json:"itemListElement"` // The list of items
}

// MarshalJSON is a custom JSON marshaller for the SchemaOrgItemList struct.
// It sets the context and type to the correct values before marshalling.
func (o *SchemaOrgItemList) MarshalJSON() ([]byte, error) {
	type Alias SchemaOrgItemList // Alias to prevent stack overflow
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
