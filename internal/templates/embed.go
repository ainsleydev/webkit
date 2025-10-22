package templates

import (
	"embed"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/spf13/afero"
)

//go:embed *
var Embed embed.FS

// LoadTemplate returns a parsed template from the embedded FS.
func LoadTemplate(name string) (*template.Template, error) {
	content, err := Embed.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return template.New(name).Funcs(templateFuncs()).Parse(string(content))
}

// LoadTemplateFromFS returns a parsed template from an afero filesystem.
func LoadTemplateFromFS(fs afero.Fs, name string) (*template.Template, error) {
	content, err := afero.ReadFile(fs, name)
	if err != nil {
		return nil, err
	}
	return template.New(name).Funcs(templateFuncs()).Parse(string(content))
}

// MustLoadTemplate calls LoadTemplate but panics if parsing fails.
func MustLoadTemplate(name string) *template.Template {
	t, err := LoadTemplate(name)
	if err != nil {
		panic(err)
	}
	return t
}

func templateFuncs() template.FuncMap {
	funcs := sprig.FuncMap()
	funcs["ghVar"] = githubVariable
	funcs["ghSecret"] = githubSecret
	funcs["ghInput"] = githubInput
	funcs["ghEnv"] = githubEnv
	return funcs
}
