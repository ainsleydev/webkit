package schemaorg

// Address represents a structured data definition for the physical or mailing
// address of an organization according to schema.org.
//
// See:
// - https://schema.org/PostalAddress
type Address struct {
	// 	The street address.
	//	For example, ainsley.dev, 71-75 Shelton Street, Covent Garden, London, WC2H 9JQ
	StreetAddress string `json:"streetAddress,omitempty"`

	// The locality in which the street address is, and which is in the region.
	// For example, London
	Locality string `json:"addressLocality,omitempty"`

	// The region in which the locality is, and which is in the country.
	// For example, Greater London or another appropriate first-level Administrative division.
	Region string `json:"addressRegion,omitempty"`

	// The country. For example, UK.
	// You can also provide the two-letter ISO 3166-1 alpha-2 country code.
	Country string `json:"addressCountry,omitempty"`

	// The postal code. For example, WC2H 9JQ
	PostalCode string `json:"postalCode,omitempty"`

	// The post office box number for PO box addresses.
	PostOfficeBoxNumber string `json:"postOfficeBoxNumber,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface to generate
// the JSON-LD for the Address.
func (s *Address) MarshalJSON() ([]byte, error) {
	type Alias Address // Define an alias type to avoid stack overflow
	return marshal(&struct {
		Type string `json:"@type"`
		*Alias
	}{
		Type:  "PostalAddress",
		Alias: (*Alias)(s),
	})
}
