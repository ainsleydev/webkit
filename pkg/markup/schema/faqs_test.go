package schemaorg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFAQPage_MarshalJSON(t *testing.T) {
	in := FAQPage{
		{Question: "What is the best color?", Answer: "Blue"},
		{Question: "What is the best animal?", Answer: "<p>Cat</p>"},
	}

	want := `{
		"@context": "https://schema.org",
		"@type": "FAQPage",
		"mainEntity": [
			{
				"@type":"Question",
				"name":"What is the best color?",
				"acceptedAnswer":{"@type":"Answer","text":"Blue"}
			},
			{
				"@type":"Question",
				"name":"What is the best animal?",
				"acceptedAnswer":{"@type":"Answer","text":"<p>Cat</p>"}
			}
		]
	}`

	got, err := in.MarshalJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, want, string(got))
}
