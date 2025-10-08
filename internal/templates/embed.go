package templates

import (
	"embed"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

//go:embed *
var embeddedTemplates embed.FS

// LoadTemplate returns a parsed template from embedded FS
func LoadTemplate(name string) (*template.Template, error) {
	content, err := embeddedTemplates.ReadFile(name)
	if err != nil {
		return nil, err
	}

	funcs := sprig.FuncMap()
	funcs["secret"] = secret

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
