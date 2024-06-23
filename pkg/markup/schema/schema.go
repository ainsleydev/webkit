package schemaorg

import (
	"encoding/json"
	"log/slog"
)

// TODO:
// - Add more structured data types (Product, Article, QA)

const (
	// Context is the schema.org context definition which is
	// defined in every JSON-LD object.
	Context = "https://schema.org"
)

// ToLDJSONScript marshals the given value to a script tag with
// type application/ld+json.
func ToLDJSONScript(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		slog.Error("Marshalling application/ld+json script: " + err.Error())
		return ""
	}
	return `<script type="application/ld+json">` + string(b) + `</script>`
}

// marshal is a helper function to marshal the JSON-LD object
// with tab indentation.
func marshal(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "\t")
}
