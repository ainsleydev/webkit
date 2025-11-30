package templates

import (
	"embed"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/ainsleydev/webkit/internal/fsext"
)

//go:embed *
var Embed embed.FS

// LoadTemplate returns a parsed template from the embedded FS.
func LoadTemplate(name string) (*template.Template, error) {
	content, err := fsext.ReadFromEmbed(Embed, name)
	if err != nil {
		return nil, err
	}
	return template.New(name).Funcs(templateFuncs()).Parse(content)
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
	funcs["ghExpr"] = githubExpression
	funcs["ghVar"] = githubVariable
	funcs["ghSecret"] = githubSecret
	funcs["ghInput"] = githubInput
	funcs["ghEnv"] = githubEnv
	funcs["prettyConfigKey"] = prettyConfigKey
	return funcs
}
