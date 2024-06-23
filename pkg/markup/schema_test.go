package markup

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaOrgFAQPage_MarshalJSON(t *testing.T) {
	tt := map[string]struct {
		input SchemaOrgFAQPage
		want  string
	}{
		"OK": {
			input: SchemaOrgFAQPage{
				{Question: "What is the best color?", Answer: "Blue"},
			},
			want: `{
				"@context": "https://schema.org",
				"@type": "FAQPage",
				"mainEntity": [
					{
						"@type":"Question",
						"name":"What is the best color?",
						"acceptedAnswer":{"@type":"Answer","text":"Blue"}
					}
				]
			}`,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := test.input.MarshalJSON()
			fmt.Println(string(got))
			assert.Nil(t, err)
			assert.JSONEq(t, test.want, string(got))
		})
	}
}
