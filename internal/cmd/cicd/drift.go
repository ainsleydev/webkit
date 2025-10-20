package cicd

import (
	"context"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

var DriftCmd = &cli.Command{
	Name:   "drift",
	Usage:  "Creates the drift detection workflow",
	Action: cmdtools.Wrap(DriftDetection),
}

// DriftDetection simply copies the drift detection workflow.
func DriftDetection(_ context.Context, input cmdtools.CommandInput) error {
	gen := scaffold.New(input.FS, input.Manifest)
	path := filepath.Join(workflowsPath, "drift.yaml")

	return gen.CopyFromEmbed(templates.Embed, path, path)
}
