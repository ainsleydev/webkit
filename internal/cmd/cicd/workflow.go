package cicd

import (
	"context"
	"path/filepath"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

// generateWorkflow is a helper function to generate workflows with app filtering.
func generateWorkflow(
	_ context.Context,
	input cmdtools.CommandInput,
	workflowName string,
	filter func(appdef.App) bool,
	data map[string]any,
) error {
	appDef := input.AppDef()

	// Filter apps based on the provided filter function.
	var filteredApps []appdef.App
	for _, app := range appDef.Apps {
		if filter(app) {
			filteredApps = append(filteredApps, app)
		}
	}

	// If no apps match the filter, skip generating the workflow.
	if len(filteredApps) == 0 {
		return nil
	}

	// Add filtered apps to data.
	data["Apps"] = filteredApps

	tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, workflowName+".yaml.tmpl"))
	path := filepath.Join(workflowsPath, workflowName+".yaml")

	// Track all apps as sources for this workflow.
	var trackingOptions []scaffold.Option
	for _, app := range filteredApps {
		trackingOptions = append(trackingOptions, scaffold.WithTracking(manifest.SourceApp(app.Name)))
	}

	return input.Generator().Template(path, tpl, data, trackingOptions...)
}
