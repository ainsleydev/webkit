package templates

import (
	"embed"
	"text/template"
)

//go:embed *
var embeddedTemplates embed.FS

// LoadTemplate returns a parsed template from embedded FS
func LoadTemplate(name string) (*template.Template, error) {
	content, err := embeddedTemplates.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return template.New(name).Parse(string(content))
}
