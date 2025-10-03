package cmd

import (
	"context"
	"log/slog"

	"github.com/urfave/cli/v3"

	cmdtools "github.com/ainsleydev/webkit/internal/cmd/internal"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

var updateCmd = &cli.Command{
	Name:   "update",
	Action: cmdtools.WrapCommand(update),
}

func update(_ context.Context, input cmdtools.CommandInput) error {
	gen := cgtools.NewGenerator(input.FS)

	err := gen.GenerateTemplate(".editorconfig", templates.MustLoadTemplate(".editorconfig"), nil)
	if err != nil {
		return err
	}

	err = gen.GenerateTemplate(".gitignore", templates.MustLoadTemplate(".gitignore"), nil)
	if err != nil {
		return err
	}

	slog.Info("Created file, nice")

	return nil
}
