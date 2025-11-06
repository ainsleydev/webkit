package cicd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
)

var DeployDigitalOceanVMCmd = &cli.Command{
	Name:   "deploy-digitalocean-vm",
	Usage:  "Generate deploy workflow for DigitalOcean VM apps",
	Action: cmdtools.Wrap(DeployDigitalOceanVMWorkflow),
}

// DeployDigitalOceanVMWorkflow creates a deployment workflow for DigitalOcean VM apps.
func DeployDigitalOceanVMWorkflow(ctx context.Context, input cmdtools.CommandInput) error {
	return generateWorkflow(ctx, input, "deploy-digitalocean-vm", func(app appdef.App) bool {
		return app.Build.Dockerfile != "" && app.ShouldRelease() &&
			app.Infra.Provider == "digitalocean" && app.Infra.Type == "vm"
	}, map[string]any{
		"TerraformVersion": infra.TerraformVersion,
	})
}
