// Package templates provides utilities for loading and parsing embedded Go templates
// with custom functions for GitHub Actions workflow generation.
package templates

import (
	"embed"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

//go:embed *
var Embed embed.FS

// LoadTemplate returns a parsed template from the embedded filesystem.
func LoadTemplate(name string) (*template.Template, error) {
	content, err := Embed.ReadFile(name)
	if err != nil {
		return nil, err
	}

	funcs := sprig.FuncMap()
	funcs["ghVar"] = githubVariable
	funcs["ghSecret"] = githubSecret
	funcs["ghInput"] = githubInput
	funcs["ghEnv"] = githubEnv

	return template.New(name).Funcs(funcs).Parse(string(content))
}

// MustLoadTemplate calls LoadTemplate but panics if the
// template could not be parsed.
func MustLoadTemplate(name string) *template.Template {
	t, err := LoadTemplate(name)
	if err != nil {
		panic(err)
	}
	return t
}
