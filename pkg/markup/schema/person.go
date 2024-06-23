package schemaorg

// Person type (alive, dead, undead, or fictional). Which can be used for
// associating someone with a company, article or any other entity.
//
// See:
// - https://schema.org/Person
type Person struct {
	// The full name of the person.
	Name string `json:"name"`

	// A URL to a page that unambiguously identifies the person.
	// E.g. the URL of the person's Wikipedia page, Wikidata entry, or official website.
	URL string `json:"url,omitempty"`

	// The persons' email address
	Email string `json:"email,omitempty"`

	// The person's image as a URL string.
	Image string `json:"image,omitempty"`

	// The job title of the person (for example, Financial Manager).
	JobTitle string `json:"jobTitle,omitempty"`

	// The telephone number.
	Telephone string `json:"telephone,omitempty"`

	// URL of a reference Web page that unambiguously indicates the item's identity.
	// E.g. the URL of the item's Wikipedia page, Wikidata entry, or official website.
	SameAs []string `json:"sameAs,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface to generate
// the JSON-LD for a Person.
func (s *Person) MarshalJSON() ([]byte, error) {
	type Alias Person
	return marshal(&struct {
		*Alias
		Type string `json:"@type"`
	}{
		Alias: (*Alias)(s),
		Type:  "Person",
	})
}
