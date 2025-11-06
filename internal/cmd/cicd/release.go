package cicd

import (
	"context"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

var ReleaseCmd = &cli.Command{
	Name:   "release",
	Usage:  "Generate release workflow for apps",
	Action: cmdtools.Wrap(ReleaseWorkflow),
}

// ReleaseWorkflow creates a release workflow for all apps with builds enabled.
func ReleaseWorkflow(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	// Filter apps to only include those with builds enabled.
	var appsToRelease []appdef.App
	for _, app := range appDef.Apps {
		// Only include apps that have a Dockerfile and should be released.
		if app.Build.Dockerfile != "" && app.ShouldRelease() {
			appsToRelease = append(appsToRelease, app)
		}
	}

	// If no apps to release, skip generating the workflow.
	if len(appsToRelease) == 0 {
		return nil
	}

	tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, "release.yaml.tmpl"))
	path := filepath.Join(workflowsPath, "release.yaml")

	data := map[string]any{
		"Apps":             appsToRelease,
		"TerraformVersion": infra.TerraformVersion,
	}

	// Track all apps as sources for this workflow.
	var trackingOptions []scaffold.Option
	for _, app := range appsToRelease {
		trackingOptions = append(trackingOptions, scaffold.WithTracking(manifest.SourceApp(app.Name)))
	}

	return input.Generator().Template(path, tpl, data, trackingOptions...)
}
