package schemaorg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrganisation_MarshalJSON(t *testing.T) {
	org := Organisation{
		ID:          "https://www.example.com/organization",
		URL:         "https://www.example.com",
		LegalName:   "Example Corporation",
		Description: "The example corporation is well-known for producing high-quality widgets",
		Logo:        "https://www.example.com/images/logo.png",
		SameAs: []string{
			"https://example.net/profile/example1234",
			"https://example.org/example1234",
		},
		Address: Address{
			StreetAddress:   "Rue Improbable 99",
			AddressLocality: "Paris",
			AddressRegion:   "Ile-de-France",
			AddressCountry:  "FR",
			PostalCode:      "75001",
		},
	}

	want := `{
		"@context": "https://schema.org",
		"@type": "Organisation",
		"@id": "https://www.example.com/organization",
		"url": "https://www.example.com",
		"legalName": "Example Corporation",
		"description": "The example corporation is well-known for producing high-quality widgets",
		"logo": "https://www.example.com/images/logo.png",
		"sameAs": [
			"https://example.net/profile/example1234",
			"https://example.org/example1234"
		],
		"address": {
			"@type": "PostalAddress",
			"streetAddress": "Rue Improbable 99",
			"addressLocality": "Paris",
			"addressRegion": "Ile-de-France",
			"addressCountry": "FR",
			"postalCode": "75001"
		}
	}`

	got, err := org.MarshalJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, want, string(got))
}
