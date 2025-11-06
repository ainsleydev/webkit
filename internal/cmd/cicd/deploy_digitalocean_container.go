package cicd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
)

var DeployDigitalOceanContainerCmd = &cli.Command{
	Name:   "deploy-digitalocean-container",
	Usage:  "Generate deploy workflow for DigitalOcean Container apps",
	Action: cmdtools.Wrap(DeployDigitalOceanContainerWorkflow),
}

// DeployDigitalOceanContainerWorkflow creates a deployment workflow for DigitalOcean Container apps.
func DeployDigitalOceanContainerWorkflow(ctx context.Context, input cmdtools.CommandInput) error {
	return generateWorkflow(ctx, input, "deploy-digitalocean-container", func(app appdef.App) bool {
		return app.Build.Dockerfile != "" && app.ShouldRelease() &&
			app.Infra.Provider == "digitalocean" && app.Infra.Type == "container"
	}, map[string]any{})
}
