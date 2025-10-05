package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

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
			scaffoldCmd,
			scratchCmd,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err.Error())
	}
}
