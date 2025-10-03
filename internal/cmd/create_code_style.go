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

	err := gen.GenerateTemplate(".editorconfig", templates.MustLoadTemplate(".editorconfig"), nil)
	if err != nil {
		return err
	}

	// TODO,
	// .prettierrc
	// .prettierignore
	// .eslintrc
	// .eslintignore
	// .stylelintrc

	return nil
}
