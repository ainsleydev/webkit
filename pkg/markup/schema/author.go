package schemaorg

// Person - TODO
type Person struct {
	Name      string `json:"name"`
	URL       string `json:"url,omitempty"`
	Email     string `json:"email,omitempty"`
	Image     string `json:"image,omitempty"`
	JobTitle  string `json:"jobTitle,omitempty"`
	Telephone string `json:"telephone,omitempty"`
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
