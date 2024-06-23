package schemaorg

// Organisation represents a structured data definition for an organisation
// This can be used to let search engines about administrative details.
//
// TODO: Not all fields are included her, such as VAT ID, founding date, etc.
//
// See:
// - https://schema.org/Organization
// - https://developers.google.com/search/docs/appearance/structured-data/organization
//
//	{
//	   "@context": "https://schema.org",
//	   "@type": "Organization",
//	   "image": "https://www.example.com/example_image.jpg",
//	   "url": "https://www.example.com",
//	   "sameAs": [
//	       "https://example.net/profile/example1234",
//	       "https://example.org/example1234"
//	   ],
//	   "logo": "https://www.example.com/images/logo.png",
//	   "name": "Example Corporation",
//	   "description": "The example corporation is well-known for producing high-quality widgets",
//	   "email": "contact@example.com",
//	   "telephone": "+47-99-999-9999",
//	   "address": {
//	       "@type": "PostalAddress",
//	       "streetAddress": "Rue Improbable 99",
//	       "addressLocality": "Paris",
//	       "addressCountry": "FR",
//	       "addressRegion": "Ile-de-France",
//	       "postalCode": "75001"
//	   },
//	   "vatID": "FR12345678901",
//	   "iso6523Code": "0199:724500PMK2A2M1SQQ228"
//	}
type Organisation struct {
	ID          string   `json:"@id"`         // Full URL
	URL         string   `json:"url"`         // Full URL
	LegalName   string   `json:"legalName"`   // The legal name of the organisation
	Description string   `json:"description"` // A description of the organisation, can be the same as the tagline.
	Logo        string   `json:"logo"`        // Full URL, no SVGs
	SameAs      []string `json:"sameAs"`      // An array of full social media URLs
	Address     Address  `json:"address"`
}

// MarshalJSON implements the json.Marshaler interface to generate
// the JSON-LD for the Organisation.
func (s *Organisation) MarshalJSON() ([]byte, error) {
	type Alias Organisation // Define an alias type to avoid stack overflow
	return marshal(&struct {
		Context string `json:"@context"`
		Type    string `json:"@type"`
		*Alias
	}{
		Context: Context,
		Type:    "Organisation",
		Alias:   (*Alias)(s),
	})
}

// Address represents a structured data definition for the physical or mailing
// address of an organization according to schema.org.
//
// See: https://schema.org/PostalAddress
type Address struct {
	StreetAddress   string `json:"streetAddress"`   // I.e ainsley.dev, 71-75 Shelton Street, Covent Garden, London, WC2H 9JQ
	AddressLocality string `json:"addressLocality"` // I.e London
	AddressRegion   string `json:"addressRegion"`   // I.e Greater London
	AddressCountry  string `json:"addressCountry"`  // I.e UK
	PostalCode      string `json:"postalCode"`      // I.e WC2H 9JQ
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
