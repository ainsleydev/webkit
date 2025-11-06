package cicd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
)

var DeployAppCmd = &cli.Command{
	Name:   "deploy-app",
	Usage:  "Generate deploy-app workflow for all deployable apps",
	Action: cmdtools.Wrap(DeployAppWorkflow),
}

// DeployAppWorkflow creates a unified deploy workflow for all apps.
func DeployAppWorkflow(ctx context.Context, input cmdtools.CommandInput) error {
	return generateWorkflow(ctx, input, "deploy-app", func(app appdef.App) bool {
		return app.Build.Dockerfile != "" && app.ShouldRelease() &&
			app.Infra.Provider == "digitalocean"
	}, map[string]any{
		"TerraformVersion": infra.TerraformVersion,
	})
}
