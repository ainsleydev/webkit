package schemaorg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress_MarshalJSON(t *testing.T) {
	address := Address{
		StreetAddress: "71-75 Shelton Street, Covent Garden",
		Locality:      "London",
		Region:        "Greater London",
		Country:       "UK",
		PostalCode:    "WC2H 9JQ",
	}

	want := `{
		"@type": "PostalAddress",
		"streetAddress": "71-75 Shelton Street, Covent Garden",
		"addressLocality": "London",
		"addressRegion": "Greater London",
		"addressCountry": "UK",
		"postalCode": "WC2H 9JQ"
	}`

	got, err := address.MarshalJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, want, string(got))
}
