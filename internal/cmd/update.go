package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/cicd"
	"github.com/ainsleydev/webkit/internal/cmd/files"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
)

var updateCmd = &cli.Command{
	Name:        "update",
	Usage:       "Update project dependencies from app.json",
	Description: "Rebuilds all generated files based on current app.json configuration",
	Action:      cmdtools.Wrap(update),
}

var updateOps = []cmdtools.RunCommand{
	files.CreateCodeStyleFiles,
	files.CreateGitSettings,
	files.CreatePackageJson,
	cicd.CreatePRWorkflow,
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
