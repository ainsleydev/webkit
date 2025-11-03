package cicd

import (
	"context"
	"path/filepath"
	"strings"

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
	tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, "pr.yaml.tmpl"))

	// Build a map of resource names to Resource objects for quick lookup.
	resourceMap := make(map[string]appdef.Resource)
	for _, resource := range appDef.Resources {
		resourceMap[resource.Name] = resource
	}

	// Find Payload apps with Postgres dependencies.
	var appsWithDBs []AppWithDatabase
	for _, app := range appDef.Apps {
		// Only check Payload apps.
		if app.Type != appdef.AppTypePayload {
			continue
		}

		// Check if this app has a Postgres dependency.
		var dbResource *appdef.Resource
		app.Env.WalkE(func(entry appdef.EnvWalkEntry) error {
			// Only check resource-type env vars.
			if entry.Source != appdef.EnvSourceResource {
				return nil
			}

			// Parse resource reference (format: "resource_name.output_name").
			valueStr, ok := entry.Value.(string)
			if !ok {
				return nil
			}

			parts := strings.SplitN(valueStr, ".", 2)
			if len(parts) != 2 {
				return nil
			}

			resourceName := parts[0]

			// Check if this resource is a Postgres database.
			if resource, exists := resourceMap[resourceName]; exists {
				if resource.Type == appdef.ResourceTypePostgres {
					dbResource = &resource
				}
			}

			return nil
		})

		// If we found a Postgres dependency, add this app to the list.
		if dbResource != nil {
			enviro := env.Production

			appsWithDBs = append(appsWithDBs, AppWithDatabase{
				App:      app,
				Database: *dbResource,
				Secrets: map[string]string{
					"DatabaseURL": dbResource.GitHubSecretName(enviro, "connection_url"),
					"DatabaseID":  dbResource.GitHubSecretName(enviro, "id"),
				},
			})
		}
	}

	data := map[string]any{
		"Apps":          appDef.Apps,
		"AppsWithDBs":   appsWithDBs,
		"HasMigrations": len(appsWithDBs) > 0,
	}
	file := filepath.Join(workflowsPath, "pr.yaml")

	return input.Generator().Template(file, tpl, data,
		scaffold.WithTracking(manifest.SourceProject()),
	)
}
