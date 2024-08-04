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
	// Full URL
	ID string `json:"@id"`
	// Full URL
	URL string `json:"url"`
	// The legal name of the organisation
	LegalName string `json:"legalName"`
	// A description of the organisation, can be the same as the tagline.
	Description string `json:"description"`
	// Full URL, no SVGs
	Logo string `json:"logo"`
	// An array of full social media URLs
	SameAs  []string `json:"sameAs"`
	Address Address  `json:"address"`
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
