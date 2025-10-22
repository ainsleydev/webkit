package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
)

var driftCmd = &cli.Command{
	Name:        "drift",
	Usage:       "Detect manual modifications to generated files",
	Description: "Checks if generated files have been manually modified or deleted since the last webkit update",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "format",
			Usage: "Output format: text, markdown, or json",
			Value: "text",
		},
	},
	Action: cmdtools.Wrap(drift),
}

// drift detects if tracked files have been manually modified or deleted.
//
// Note: This only detects drift from the last webkit update. If templates have
// been updated in a newer version of WebKit but webkit update hasn't been run,
// this will not detect that the files are outdated.
func drift(ctx context.Context, input cmdtools.CommandInput) error {
	cmd := input.Command
	format := cmd.String("format")
	printer := input.Printer()

	// Validate format
	validFormats := map[string]bool{"text": true, "markdown": true, "json": true}
	if !validFormats[format] {
		format = "text"
	}

	// Only show info message for text format (not for structured output).
	if format == "text" {
		printer.Info("Checking for drift...")
	}

	// Run update in memory to see what should be generated
	memFS := afero.NewMemMapFs()
	memTracker := manifest.NewTracker()

	memInput := cmdtools.CommandInput{
		FS:          memFS,
		AppDefCache: input.AppDef(),
		Manifest:    memTracker,
		BaseDir:     input.BaseDir,
	}
	memInput.Printer().SetWriter(io.Discard)

	if err := update(ctx, memInput); err != nil {
		return errors.Wrap(err, "running update for in-mem manifest")
	}

	// Compare actual vs expected
	drifted, err := manifest.DetectDrift(input.FS, memFS)
	if err != nil {
		return errors.Wrap(err, "detecting drift")
	}

	// Format and output results
	output, err := formatDriftOutput(drifted, format)
	if err != nil {
		return errors.Wrap(err, "formatting output")
	}

	printer.Println(output)

	// Exit with error code if drift detected
	if len(drifted) > 0 {
		return cmdtools.ExitWithCode(1)
	}

	return nil
}

// driftFormatter is a function type for formatting drift results.
type driftFormatter func([]manifest.DriftEntry) (string, error)

// driftFormatters maps format names to their formatting functions.
var driftFormatters = map[string]driftFormatter{
	"text": func(drifted []manifest.DriftEntry) (string, error) {
		return formatDriftAsText(drifted), nil
	},
	"markdown": func(drifted []manifest.DriftEntry) (string, error) {
		return formatDriftAsMarkdown(drifted), nil
	},
	"json": formatDriftAsJSON,
}

// formatDriftOutput formats drift results based on the requested format.
func formatDriftOutput(drifted []manifest.DriftEntry, format string) (string, error) {
	formatter, exists := driftFormatters[format]
	if !exists {
		return "", fmt.Errorf("unsupported format: %s", format)
	}
	return formatter(drifted)
}

// formatDriftAsText formats drift results as human-readable text output.
func formatDriftAsText(drifted []manifest.DriftEntry) string {
	if len(drifted) == 0 {
		return "✓ No drift detected - all files are up to date"
	}

	var output strings.Builder

	// Group by type
	modifiedFiles := manifest.DriftReasonModified.FilterEntries(drifted)
	outdatedFiles := manifest.DriftReasonOutdated.FilterEntries(drifted)
	newFiles := manifest.DriftReasonNew.FilterEntries(drifted)
	deletedFiles := manifest.DriftReasonDeleted.FilterEntries(drifted)

	// Report findings
	if len(modifiedFiles) > 0 {
		output.WriteString("⚠ Manual modifications detected:\n")
		output.WriteString("  These files were manually edited:\n")
		for _, d := range modifiedFiles {
			output.WriteString(fmt.Sprintf("    • %s\n", d.Path))
		}
		output.WriteString("\n")
	}

	if len(outdatedFiles) > 0 {
		output.WriteString("Outdated files detected:\n")
		output.WriteString("app.json changed, these files need regeneration:\n")
		for _, d := range outdatedFiles {
			output.WriteString(fmt.Sprintf("    • %s\n", d.Path))
		}
		output.WriteString("\n")
	}

	if len(newFiles) > 0 {
		output.WriteString("Missing files detected:\n")
		output.WriteString("These files should exist:\n")
		for _, d := range newFiles {
			output.WriteString(fmt.Sprintf("    • %s\n", d.Path))
		}
		output.WriteString("\n")
	}

	if len(deletedFiles) > 0 {
		output.WriteString("Orphaned files detected:\n")
		output.WriteString("These files should be removed:\n")
		for _, d := range deletedFiles {
			output.WriteString(fmt.Sprintf("    • %s\n", d.Path))
		}
		output.WriteString("\n")
	}

	output.WriteString("Run 'webkit update' to sync all files")

	return output.String()
}

// formatDriftAsMarkdown formats drift results as GitHub-friendly markdown.
func formatDriftAsMarkdown(drifted []manifest.DriftEntry) string {
	if len(drifted) == 0 {
		return "## WebKit Drift Detection\n\n✅ **No drift detected** - all files are up to date"
	}

	var output strings.Builder
	output.WriteString("## WebKit Drift Detection\n\n")

	// Group by type
	modifiedFiles := manifest.DriftReasonModified.FilterEntries(drifted)
	outdatedFiles := manifest.DriftReasonOutdated.FilterEntries(drifted)
	newFiles := manifest.DriftReasonNew.FilterEntries(drifted)
	deletedFiles := manifest.DriftReasonDeleted.FilterEntries(drifted)

	// Report findings with collapsible sections for better readability
	if len(modifiedFiles) > 0 {
		output.WriteString(fmt.Sprintf("### ⚠ Manual modifications detected (%d file", len(modifiedFiles)))
		if len(modifiedFiles) != 1 {
			output.WriteString("s")
		}
		output.WriteString(")\n\n")
		output.WriteString("These files were manually edited:\n")
		for _, d := range modifiedFiles {
			output.WriteString(fmt.Sprintf("- `%s`\n", d.Path))
		}
		output.WriteString("\n")
	}

	if len(outdatedFiles) > 0 {
		output.WriteString(fmt.Sprintf("### 📝 Outdated files detected (%d file", len(outdatedFiles)))
		if len(outdatedFiles) != 1 {
			output.WriteString("s")
		}
		output.WriteString(")\n\n")
		output.WriteString("app.json changed, these files need regeneration:\n")
		for _, d := range outdatedFiles {
			output.WriteString(fmt.Sprintf("- `%s`\n", d.Path))
		}
		output.WriteString("\n")
	}

	if len(newFiles) > 0 {
		output.WriteString(fmt.Sprintf("### 📄 Missing files detected (%d file", len(newFiles)))
		if len(newFiles) != 1 {
			output.WriteString("s")
		}
		output.WriteString(")\n\n")
		output.WriteString("These files should exist:\n")
		for _, d := range newFiles {
			output.WriteString(fmt.Sprintf("- `%s`\n", d.Path))
		}
		output.WriteString("\n")
	}

	if len(deletedFiles) > 0 {
		output.WriteString(fmt.Sprintf("### 🗑️ Orphaned files detected (%d file", len(deletedFiles)))
		if len(deletedFiles) != 1 {
			output.WriteString("s")
		}
		output.WriteString(")\n\n")
		output.WriteString("These files should be removed:\n")
		for _, d := range deletedFiles {
			output.WriteString(fmt.Sprintf("- `%s`\n", d.Path))
		}
		output.WriteString("\n")
	}

	output.WriteString("---\n")
	output.WriteString("**Action Required:** Run `webkit update` to sync all files\n")

	return output.String()
}

// formatDriftAsJSON formats drift results as JSON.
func formatDriftAsJSON(drifted []manifest.DriftEntry) (string, error) {
	type driftOutput struct {
		DriftDetected bool                  `json:"drift_detected"`
		TotalFiles    int                   `json:"total_files"`
		Files         []manifest.DriftEntry `json:"files"`
		Summary       map[string]int        `json:"summary"`
	}

	// Build summary counts
	summary := map[string]int{
		"modified": len(manifest.DriftReasonModified.FilterEntries(drifted)),
		"outdated": len(manifest.DriftReasonOutdated.FilterEntries(drifted)),
		"new":      len(manifest.DriftReasonNew.FilterEntries(drifted)),
		"deleted":  len(manifest.DriftReasonDeleted.FilterEntries(drifted)),
	}

	output := driftOutput{
		DriftDetected: len(drifted) > 0,
		TotalFiles:    len(drifted),
		Files:         drifted,
		Summary:       summary,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", errors.Wrap(err, "marshalling drift to JSON")
	}

	return string(data), nil
}
