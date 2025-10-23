package cicd

import (
	"context"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

var PRCmd = &cli.Command{
	Name:   "pr",
	Usage:  "Creates PR workflows for apps and drift detection",
	Action: cmdtools.Wrap(PR),
}

// PR generates all PR-related GitHub workflows including
// drift detection and app-specific PR workflows.
func PR(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, "pr.yaml.tmpl"))

	data := map[string]any{
		"Apps": appDef.Apps,
	}
	file := filepath.Join(workflowsPath, "pr.yaml")

	return input.Generator().Template(file, tpl, data,
		scaffold.WithTracking(manifest.SourceProject()),
	)
}
