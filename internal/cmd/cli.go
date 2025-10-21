// Package cmd provides the main CLI commands and subcommands for the WebKit CLI tool.
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/cicd"
	"github.com/ainsleydev/webkit/internal/cmd/env"
	"github.com/ainsleydev/webkit/internal/cmd/infra"
	"github.com/ainsleydev/webkit/internal/cmd/secrets"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/pkg/log"
)

// Run initializes and executes the WebKit CLI application with all registered commands.
func Run() {
	cmd := &cli.Command{
		Name:  "webkit",
		Usage: "make an explosive entrance",
		Before: func(ctx context.Context, _ *cli.Command) (context.Context, error) {
			log.Bootstrap("Webkit")
			return ctx, nil
		},
		Commands: []*cli.Command{
			updateCmd,
			scaffoldCmd,
			secrets.Command,
			env.Command,
			infra.Command,
			cicd.Command,
			driftCmd,
			printCmd,
			scratchCmd,
			versionCmd,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		var silentErr *cmdtools.ExitError
		if errors.As(err, &silentErr) {
			os.Exit(silentErr.Code)
			return
		}
		fmt.Println(err.Error()) //nolint
	}
}
