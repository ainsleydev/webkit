package tpl

import (
	"fmt"
	"html/template"
	"log/slog"

	"github.com/goccy/go-json"
)

// Funcs is a map of utility functions that can be used in the std
// html/template package.
var Funcs = template.FuncMap{
	"json": func(v any) template.JS {
		a, err := json.Marshal(v)
		if err != nil {
			slog.Error("Error marshalling JSON: " + err.Error())
			return "{}"
		}
		return template.JS(a)
	},
	"jsonPretty": func(v any) template.JS {
		a, err := json.MarshalIndent(v, "", "\t")
		if err != nil {
			slog.Error("Error marshalling JSON: " + err.Error())
			return "{}"
		}
		return template.JS(a)
	},
	"safeHTML": func(v any) template.HTML {
		return template.HTML(fmt.Sprint(v))
	},
}
