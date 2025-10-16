package cmd

import (
	"context"
	"os"

	"github.com/charmbracelet/lipgloss/tree"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/printer"
)

var printCmd = &cli.Command{
	Name:        "printer",
	Description: "Utility to function to display the output of the printer",
	Hidden:      true,
	Action: cmdtools.Wrap(func(ctx context.Context, input cmdtools.CommandInput) error {
		c := printer.New(os.Stdout)

		c.Success("Deployment completed successfully.")
		c.Info("Fetching configuration...")
		c.Warn("Config file deprecated, using fallback.")
		c.Error("Failed to connect to database.")

		c.Table(
			[]string{"Service", "Status"},
			[][]string{
				{"Database", "Running"},
				{"API", "Stopped"},
			},
		)

		c.List("Config loaded", "Deployment started", "Done ✅")

		c.Tree("project",
			"cmd",
			"internal",
			tree.New().Root("pkg").Child("config", "printer", "infra"),
		)

		return nil
	}),
}
