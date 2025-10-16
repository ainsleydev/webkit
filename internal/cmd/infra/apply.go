package infra

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
)

var ApplyCmd = &cli.Command{
	Name:   "apply",
	Usage:  "Creates or updates infrastructure based off the apps and resources defined in app.json",
	Action: cmdtools.Wrap(Plan),
}

func Apply(_ context.Context, input cmdtools.CommandInput) error {
	_ = input.AppDef()

	return nil
}
