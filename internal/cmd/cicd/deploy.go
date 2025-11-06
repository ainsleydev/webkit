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
func DeployContainerWorkflow(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	// Filter apps to only include container apps with builds enabled.
	var containerApps []appdef.App
	for _, app := range appDef.Apps {
		// Only include apps that have a Dockerfile, should be released, and are container type.
		if app.Build.Dockerfile != "" && app.ShouldRelease() &&
			app.Infra.Provider == "digitalocean" && app.Infra.Type == "container" {
			containerApps = append(containerApps, app)
		}
	}

	// If no container apps to deploy, skip generating the workflow.
	if len(containerApps) == 0 {
		return nil
	}

	tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, "deploy-container.yaml.tmpl"))
	path := filepath.Join(workflowsPath, "deploy-container.yaml")

	data := map[string]any{
		"Apps": containerApps,
	}

	// Track all apps as sources for this workflow.
	var trackingOptions []scaffold.Option
	for _, app := range containerApps {
		trackingOptions = append(trackingOptions, scaffold.WithTracking(manifest.SourceApp(app.Name)))
	}

	return input.Generator().Template(path, tpl, data, trackingOptions...)
}

// DeployVMWorkflow creates a deploy-only workflow for VM apps.
func DeployVMWorkflow(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	// Filter apps to only include VM apps with builds enabled.
	var vmApps []appdef.App
	for _, app := range appDef.Apps {
		// Only include apps that have a Dockerfile, should be released, and are VM type.
		if app.Build.Dockerfile != "" && app.ShouldRelease() &&
			app.Infra.Provider == "digitalocean" && app.Infra.Type == "vm" {
			vmApps = append(vmApps, app)
		}
	}

	// If no VM apps to deploy, skip generating the workflow.
	if len(vmApps) == 0 {
		return nil
	}

	tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, "deploy-vm.yaml.tmpl"))
	path := filepath.Join(workflowsPath, "deploy-vm.yaml")

	data := map[string]any{
		"Apps":             vmApps,
		"TerraformVersion": infra.TerraformVersion,
	}

	// Track all apps as sources for this workflow.
	var trackingOptions []scaffold.Option
	for _, app := range vmApps {
		trackingOptions = append(trackingOptions, scaffold.WithTracking(manifest.SourceApp(app.Name)))
	}

	return input.Generator().Template(path, tpl, data, trackingOptions...)
}
