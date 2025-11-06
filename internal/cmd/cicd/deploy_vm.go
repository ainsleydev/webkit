package cicd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
)

var DeployVMCmd = &cli.Command{
	Name:   "deploy-vm",
	Usage:  "Generate deploy workflow for VM apps",
	Action: cmdtools.Wrap(DeployVMWorkflow),
}

// DeployVMWorkflow creates a deployment workflow for VM apps.
func DeployVMWorkflow(ctx context.Context, input cmdtools.CommandInput) error {
	return generateWorkflow(ctx, input, "deploy-vm", func(app appdef.App) bool {
		return app.Build.Dockerfile != "" && app.ShouldRelease() &&
			app.Infra.Provider == "digitalocean" && app.Infra.Type == "vm"
	}, map[string]any{
		"TerraformVersion": infra.TerraformVersion,
	})
}
