package cicd

import (
	"context"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/templates"
)

var DriftCmd = &cli.Command{
	Name:   "drift",
	Usage:  "Creates the drift detection workflow",
	Action: cmdtools.Wrap(DriftDetection),
}

// DriftData is the template data for drift detection workflow.
type DriftData struct {
	IsDrift bool
}

// DriftDetection generates the drift detection workflow using the PR template.
func DriftDetection(_ context.Context, input cmdtools.CommandInput) error {
	tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, "pr.yaml.tmpl"))
	path := filepath.Join(workflowsPath, "drift.yaml")
	data := &DriftData{IsDrift: true}
	return input.Generator().Template(path, tpl, data)
}
