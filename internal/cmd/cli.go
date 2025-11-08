package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/cicd"
	"github.com/ainsleydev/webkit/internal/cmd/docs"
	"github.com/ainsleydev/webkit/internal/cmd/env"
	"github.com/ainsleydev/webkit/internal/cmd/infra"
	"github.com/ainsleydev/webkit/internal/cmd/payload"
	"github.com/ainsleydev/webkit/internal/cmd/secrets"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/pkg/log"
)

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
			validateCmd,
			schemaCmd,
			scaffoldCmd,
			secrets.Command,
			env.Command,
			infra.Command,
			cicd.Command,
			docs.Command,
			payload.Command,
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
