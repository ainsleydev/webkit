package cicd

import (
	"context"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/state/manifest"
	"github.com/ainsleydev/webkit/internal/templates"
	"github.com/ainsleydev/webkit/pkg/env"
)

var VMMaintenanceCmd = &cli.Command{
	Name:   "vm-maintenance",
	Usage:  "Generate server maintenance workflow for Digital Ocean VMs",
	Action: cmdtools.Wrap(VMMaintenanceWorkflow),
}

// VMMaintenanceWorkflow creates a weekly server maintenance workflow
// for all apps that use Digital Ocean VMs.
func VMMaintenanceWorkflow(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	enviro := env.Production

	if len(appDef.Apps) == 0 {
		return nil
	}

	// Check if there are any Digital Ocean or Hetzner VM apps
	hasVMApps := false
	var vmApps []appdef.App
	for _, app := range appDef.Apps {
		isDigitalOceanVM := app.Infra.Provider == appdef.ResourceProviderDigitalOcean && app.Infra.Type == "vm"
		isHetznerVM := app.Infra.Provider == appdef.ResourceProviderHetzner && app.Infra.Type == "vm"
		if isDigitalOceanVM || isHetznerVM {
			hasVMApps = true
			vmApps = append(vmApps, app)
		}
	}

	// Only generate the workflow if there are VM apps
	if !hasVMApps {
		return nil
	}

	tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, "server-maintenance.yaml.tmpl"))
	path := filepath.Join(workflowsPath, "server-maintenance.yaml")

	data := map[string]any{
		"Apps": appDef.Apps,
		"Env":  enviro,
	}

	// Track all VM apps as sources for this workflow
	var trackingOptions []scaffold.Option
	for _, app := range vmApps {
		trackingOptions = append(trackingOptions, scaffold.WithTracking(manifest.SourceApp(app.Name)))
	}

	return input.Generator().Template(path, tpl, data, trackingOptions...)
}
