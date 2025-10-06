package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

var createCodeStyleFilesCmd = &cli.Command{
	Name:   "code-style",
	Action: cmdtools.WrapCommand(createCodeStyleFiles),
}

func createCodeStyleFiles(_ context.Context, input cmdtools.CommandInput) error {
	gen := scaffold.New(input.FS)
	app := input.AppDef()

	files := map[string]string{
		".editorconfig":     ".editorconfig",
		".prettierrc":       ".prettierrc",
		".prettierignore":   ".prettierignore",
		".eslint.config.js": "eslint.config.js.tmpl",
	}

	// TODO: .stylelintrc

	for file, template := range files {
		tpl := templates.MustLoadTemplate(template)
		err := gen.Template(file, tpl, app)
		if err != nil {
			return err
		}
	}

	return nil
}
