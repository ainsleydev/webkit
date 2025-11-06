package cicd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
)

var ReleaseCmd = &cli.Command{
	Name:   "release",
	Usage:  "Generate release workflow for apps",
	Action: cmdtools.Wrap(ReleaseWorkflow),
}

// ReleaseWorkflow creates a release workflow for all apps with builds enabled.
func ReleaseWorkflow(ctx context.Context, input cmdtools.CommandInput) error {
	return generateWorkflow(ctx, input, "release", func(app appdef.App) bool {
		return app.Build.Dockerfile != "" && app.ShouldRelease()
	}, map[string]any{
		"TerraformVersion": infra.TerraformVersion,
	})
}
