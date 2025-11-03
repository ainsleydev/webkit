package cicd

import (
	"context"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
	"github.com/ainsleydev/webkit/pkg/env"
)

var PRCmd = &cli.Command{
	Name:   "pr",
	Usage:  "Creates PR workflows for apps and drift detection",
	Action: cmdtools.Wrap(PR),
}

// AppWithDatabase represents a Payload app that has a Postgres database dependency.
type AppWithDatabase struct {
	App      appdef.App
	Database appdef.Resource
	Secrets  map[string]string
}

// PR generates all PR-related GitHub workflows including
// drift detection and app-specific PR workflows.
func PR(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	return input.Generator().Template(
		filepath.Join(workflowsPath, "pr.yaml"),
		templates.MustLoadTemplate(filepath.Join(workflowsPath, "pr.yaml.tmpl")),
		map[string]any{
			"Apps":                appDef.Apps,
			"PayloadPostgresApps": findPayloadAppsWithPostgres(appDef.Apps, appDef.Resources),
		},
		scaffold.WithTracking(manifest.SourceProject()),
	)
}

// findPayloadAppsWithPostgres finds Payload apps that depend on a Postgres resource.
func findPayloadAppsWithPostgres(apps []appdef.App, resources []appdef.Resource) []AppWithDatabase {
	// Build a map of resource names to Resource objects for quick lookup.
	resourceMap := make(map[string]appdef.Resource)
	for _, resource := range resources {
		resourceMap[resource.Name] = resource
	}

	var result []AppWithDatabase
	for _, app := range apps {
		// Only check Payload apps.
		if app.Type != appdef.AppTypePayload {
			continue
		}

		// Check if Payload has a dependency of Postgres
		var dbResource *appdef.Resource
		app.Env.Walk(func(entry appdef.EnvWalkEntry) {
			if entry.Source != appdef.EnvSourceResource {
				return
			}

			resourceName, _, ok := appdef.ParseResourceReference(entry.Value)
			if !ok {
				return
			}

			if resource, exists := resourceMap[resourceName]; exists {
				if resource.Type == appdef.ResourceTypePostgres {
					dbResource = &resource
				}
			}
		})

		// If we found a Postgres dependency, add this app to the list.
		if dbResource != nil {
			enviro := env.Production

			result = append(result, AppWithDatabase{
				App:      app,
				Database: *dbResource,
				Secrets: map[string]string{
					"DatabaseURL": dbResource.GitHubSecretName(enviro, "connection_url"),
					"DatabaseID":  dbResource.GitHubSecretName(enviro, "id"),
				},
			})
		}
	}

	return result
}
