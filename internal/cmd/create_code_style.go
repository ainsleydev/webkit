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
	gen := cgtools.NewGenerator(input.FS)

	files := []string{
		".editorconfig",
		".prettierrc",
		".prettierignore",
	}

	for _, file := range files {
		// Presume that the name of the generated file is the same as the
		// one that resides in templates, but it can be changed later.
		err := gen.GenerateTemplate(file, templates.MustLoadTemplate(file), nil)
		if err != nil {
			return err
		}
	}

	// TODO,
	// .eslintrc
	// .eslintignore
	// .stylelintrc

	return nil
}
