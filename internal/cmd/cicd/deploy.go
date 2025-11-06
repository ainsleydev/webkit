package cicd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
)

var DeployContainerCmd = &cli.Command{
	Name:   "deploy-container",
	Usage:  "Generate deploy-container workflow for container apps",
	Action: cmdtools.Wrap(DeployContainerWorkflow),
}

var DeployVMCmd = &cli.Command{
	Name:   "deploy-vm",
	Usage:  "Generate deploy-vm workflow for VM apps",
	Action: cmdtools.Wrap(DeployVMWorkflow),
}

// DeployContainerWorkflow creates a deploy-only workflow for container apps.
func DeployContainerWorkflow(ctx context.Context, input cmdtools.CommandInput) error {
	return generateWorkflow(ctx, input, "deploy-container", func(app appdef.App) bool {
		return app.Build.Dockerfile != "" && app.ShouldRelease() &&
			app.Infra.Provider == "digitalocean" && app.Infra.Type == "container"
	}, map[string]any{})
}

// DeployVMWorkflow creates a deploy-only workflow for VM apps.
func DeployVMWorkflow(ctx context.Context, input cmdtools.CommandInput) error {
	return generateWorkflow(ctx, input, "deploy-vm", func(app appdef.App) bool {
		return app.Build.Dockerfile != "" && app.ShouldRelease() &&
			app.Infra.Provider == "digitalocean" && app.Infra.Type == "vm"
	}, map[string]any{
		"TerraformVersion": infra.TerraformVersion,
	})
}
