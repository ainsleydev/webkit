package files

import (
	"context"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

var codeStyleTemplates = map[string]string{
	".editorconfig":    ".editorconfig",
	".prettierrc":      ".prettierrc",
	".prettierignore":  ".prettierignore",
	"eslint.config.js": "eslint.config.js.tmpl",
	// TODO: .stylelintrc
}

// CodeStyle scaffolds' developer and formatting files for
// the project, mainly dotfiles.
//
// IDEA: Might be good in the AppDef if we could specify what files
// we want to generate or exclude from this.
func CodeStyle(_ context.Context, input cmdtools.CommandInput) error {
	app := input.AppDef()

	for file, template := range codeStyleTemplates {
		err := input.Generator().Template(file,
			templates.MustLoadTemplate(template),
			app,
			scaffold.WithTracking(manifest.SourceProject()),
		)
		if err != nil {
			return err
		}
	}

	// Only bootstrap the Go files if necessary.
	if app.ContainsGo() {
		return input.
			Generator().
			CopyFromEmbed(templates.Embed, ".golangci.yaml", ".golangci.yaml")
	}

	return nil
}
