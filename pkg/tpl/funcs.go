package tpl

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log/slog"

	"github.com/goccy/go-json"
)

type componentRenderer interface {
	Render(context.Context, io.Writer) error
}

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
	"safeAttr": func(v any) template.HTMLAttr {
		return template.HTMLAttr(fmt.Sprint(v))
	},
	"renderComponent": func(renderer componentRenderer) string {
		buf := bytes.Buffer{}
		err := renderer.Render(context.Background(), &buf)
		if err != nil {
			slog.Error("Rendering component: " + err.Error())
			return ""
		}
		return buf.String()
	},
}
