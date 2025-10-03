package cli

import (
	"context"
	"log/slog"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

var updateCmd = &cli.Command{
	Name:   "update",
	Action: wrapCommand(update),
}

func update(_ context.Context, input commandInput) error {
	gen := cgtools.NewGenerator(input.FS)

	tpl, err := templates.LoadTemplate(".editorconfig")
	if err != nil {
		return err
	}

	err = gen.GenerateTemplate(".editorconfig", tpl, nil)

	slog.Info("Created file, nice")

	return nil
}
