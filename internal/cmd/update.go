package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/cmd/internal/operations"
)

var updateCmd = &cli.Command{
	Name:        "update",
	Usage:       "Update project dependencies from app.json",
	Description: "Rebuilds all generated files based on current app.json configuration",
	Action:      cmdtools.Wrap(update),
}

var updateOps = []cmdtools.RunCommand{
	operations.CreateCodeStyleFiles,
	operations.CreateGitSettings,
	operations.CreatePackageJson,
	operations.CreateCICD,
}

func update(ctx context.Context, input cmdtools.CommandInput) error {
	for _, op := range updateOps {
		err := op(ctx, input)
		if err != nil {
			return err
		}
	}
	return nil
}
