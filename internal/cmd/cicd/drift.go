package cicd

import (
	"context"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/templates"
)

var DriftCmd = &cli.Command{
	Name:   "drift",
	Usage:  "Creates the drift detection workflow",
	Action: cmdtools.Wrap(DriftDetection),
}

// DriftDetection simply copies the drift detection workflow.
func DriftDetection(_ context.Context, input cmdtools.CommandInput) error {
	path := filepath.Join(workflowsPath, "drift.yaml")
	return input.Generator().CopyFromEmbed(templates.Embed, path, path)
}
