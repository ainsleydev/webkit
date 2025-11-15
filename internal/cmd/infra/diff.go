package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/goccy/go-json"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/printer"
)

// DiffCmd defines the diff command for detecting infrastructure changes.
var DiffCmd = &cli.Command{
	Name:  "diff",
	Usage: "Detect infrastructure changes and determine if Terraform apply is needed",
	Description: `Analyses infrastructure changes between the current app.json and a previous version
to determine if Terraform apply is needed. This helps optimise CI/CD workflows by skipping
Terraform operations when only cosmetic drift would be shown.

Exit codes:
  0 - No Terraform needed (skip)
  1 - Terraform needed (run)
  2 - Error occurred`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Output format (text, json, github)",
			Value:   "text",
		},
		&cli.StringFlag{
			Name:    "base",
			Aliases: []string{"b"},
			Usage:   "Git ref to compare against",
			Value:   "HEAD~1",
		},
		&cli.BoolFlag{
			Name:    "silent",
			Aliases: []string{"s"},
			Usage:   "Suppress informational output",
		},
	},
	Action: cmdtools.Wrap(Diff),
}

// Diff executes the diff command to analyse infrastructure changes.
func Diff(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()
	baseRef := input.Cmd().String("base")
	format := input.Cmd().String("format")
	silent := input.Cmd().Bool("silent")

	if !silent && format == "text" {
		printer.Info(fmt.Sprintf("Comparing app.json with %s", baseRef))
	}

	// Load current app.json.
	current := input.AppDef()

	// Load previous app.json from git.
	previous, err := appdef.LoadFromGit(ctx, baseRef)
	if err != nil {
		return err
	}

	// Compare definitions.
	analysis := appdef.Compare(current, previous)

	// Output based on format.
	switch format {
	case "json":
		return outputJSON(analysis)
	case "github":
		return outputGitHub(analysis)
	default:
		return outputText(analysis, printer, silent)
	}
}

// outputText outputs the analysis in human-readable text format.
func outputText(analysis appdef.ChangeAnalysis, printer *printer.Console, silent bool) error {
	if !silent {
		printer.Print(fmt.Sprintf("Decision: %s", analysis.Reason))

		if len(analysis.ChangedApps) > 0 {
			printer.Print("\nChanged apps:")
			for _, app := range analysis.ChangedApps {
				status := "unchanged"
				if app.InfraChanged {
					status = "infrastructure changed"
				} else if app.EnvChanged {
					status = "env changed"
				}
				printer.Print(fmt.Sprintf("  - %s: %s", app.Name, status))
			}
		}
		printer.Print("")
	}

	if analysis.Skip {
		if !silent {
			printer.Success("Terraform apply can be skipped")
		}
		return nil
	}

	if !silent {
		printer.Warning("Terraform apply is needed")
	}
	return cli.Exit("", 1)
}

// outputJSON outputs the analysis in JSON format.
func outputJSON(analysis appdef.ChangeAnalysis) error {
	data, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	if analysis.Skip {
		return nil
	}
	return cli.Exit("", 1)
}

// outputGitHub outputs the analysis in GitHub Actions format.
func outputGitHub(analysis appdef.ChangeAnalysis) error {
	skip := "false"
	if analysis.Skip {
		skip = "true"
	}

	// Output GitHub Actions output variables.
	fmt.Fprintf(os.Stdout, "skip_terraform=%s\n", skip)
	fmt.Fprintf(os.Stdout, "reason=%s\n", analysis.Reason)

	// Output as GitHub Actions notice.
	fmt.Fprintf(os.Stdout, "::notice::%s\n", analysis.Reason)

	if analysis.Skip {
		return nil
	}
	return cli.Exit("", 1)
}
