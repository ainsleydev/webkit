package markup

import (
	"log/slog"

	"github.com/goccy/go-json"
)

// MarshalLDJSONScript marshals the given value to a script tag with type application/ld+json.
func MarshalLDJSONScript(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		slog.Error("Marshalling application/ld+json script: " + err.Error())
		return ""
	}
	return `<script type="application/ld+json">` + string(b) + `</script>`
}
