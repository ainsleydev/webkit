package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
)

var updateCmd = &cli.Command{
	Name:   "update",
	Action: cmdtools.WrapCommand(update),
}

func update(ctx context.Context, input cmdtools.CommandInput) error {
	return nil
}
