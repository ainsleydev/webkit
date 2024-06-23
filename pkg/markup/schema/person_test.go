package schemaorg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPerson_MarshalJSON(t *testing.T) {
	person := Person{
		Name:      "John Doe",
		URL:       "https://example.com/profile/johndoe",
		Email:     "johndoe@example.com",
		Image:     "https://example.com/photos/johndoe.jpg",
		JobTitle:  "Software Engineer",
		Telephone: "+1234567890",
		SameAs: []string{
			"https://example.com/social/johndoe",
			"https://linkedin.com/in/johndoe",
		},
	}

	want := `{
		"@type": "Person",
		"name": "John Doe",
		"url": "https://example.com/profile/johndoe",
		"email": "johndoe@example.com",
		"image": "https://example.com/photos/johndoe.jpg",
		"jobTitle": "Software Engineer",
		"telephone": "+1234567890",
		"sameAs": [
			"https://example.com/social/johndoe",
			"https://linkedin.com/in/johndoe"
		]
	}`

	got, err := person.MarshalJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, want, string(got))
}
