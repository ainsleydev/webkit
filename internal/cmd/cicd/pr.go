package cicd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/templates"
)

var PRCmd = &cli.Command{
	Name:   "pr",
	Usage:  "Creates PR workflows for apps and drift detection",
	Action: cmdtools.Wrap(PRWorkflow),
}

// PRWorkflow generates all PR-related GitHub workflows including
// drift detection and app-specific PR workflows.
func PRWorkflow(_ context.Context, input cmdtools.CommandInput) error {
	tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, "pr.yaml.tmpl"))

	// Generate drift detection workflow (template detects this by empty .Name)
	driftPath := filepath.Join(workflowsPath, "drift.yaml")
	if err := input.Generator().Template(driftPath, tpl, nil); err != nil {
		return err
	}

	// Generate app-specific PR workflows
	appDef := input.AppDef()
	for _, app := range appDef.Apps {
		file := filepath.Join(workflowsPath, fmt.Sprintf("pr-%s.yaml", app.Name))
		if err := input.Generator().Template(file, tpl, &app); err != nil {
			return err
		}
	}

	return nil
}
