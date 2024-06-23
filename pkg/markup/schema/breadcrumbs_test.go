package schemaorg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBreadcrumbList_MarshalJSON(t *testing.T) {
	in := BreadcrumbList{
		{Position: 1, Name: "Home", Item: "https://example.com"},
		{Position: 2, Name: "Catalog", Item: "https://example.com/catalog"},
	}

	want := `{
		"@context": "https://schema.org",
		"@type": "BreadcrumbList",
		"itemListElement": [
			{
				"@type": "ListItem",
				"position": 1,
				"name": "Home",
				"item": "https://example.com"
			},
			{
				"@type": "ListItem",
				"position": 2,
				"name": "Catalog",
				"item": "https://example.com/catalog"
			}
		]
	}`

	got, err := in.MarshalJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, want, string(got))
}
