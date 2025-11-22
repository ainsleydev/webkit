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

var BackupCmd = &cli.Command{
	Name:   "backup",
	Usage:  "Generate backup workflows for resources",
	Action: cmdtools.Wrap(BackupWorkflow),
}

// BackupWorkflow creates backup workflows for every resource if the
// backup config is enabled.
func BackupWorkflow(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	enviro := env.Production

	tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, "backup.yaml.tmpl"))
	path := filepath.Join(workflowsPath, "backup.yaml")

	// Build nested data map with all resource secrets grouped by resource name.
	// This allows cleaner template access: {{ index $.Data .Name "DatabaseURL" }}
	secretData := make(map[string]map[string]string)
	for _, resource := range appDef.Resources {
		resourceSecrets := make(map[string]string)

		switch resource.Type {
		case appdef.ResourceTypePostgres:
			resourceSecrets["DatabaseURL"] = resource.GitHubSecretName(enviro, "connection_url")
			resourceSecrets["DatabaseID"] = resource.GitHubSecretName(enviro, "id")

		case appdef.ResourceTypeS3:
			// NOTE: S3 backup is currently only compatible with DigitalOcean Spaces.
			// Backblaze B2 and other providers are not yet supported.
			if resource.Provider != appdef.ResourceProviderDigitalOcean {
				continue
			}

			resourceSecrets["AccessKey"] = resource.GitHubSecretName(enviro, "access_key")
			resourceSecrets["SecretKey"] = resource.GitHubSecretName(enviro, "secret_key")
			resourceSecrets["Region"] = resource.GitHubSecretName(enviro, "region")
			resourceSecrets["BucketName"] = resource.GitHubSecretName(enviro, "bucket_name")

		case appdef.ResourceTypeSQLite:
			// NOTE: SQLite backup is currently only compatible with Turso.
			// The database name is constructed as ${project_name}-${resource_name} (same as in terraform).
			// Authentication is handled via TURSO_API_TOKEN environment variable.
			if resource.Provider != appdef.ResourceProviderTurso {
				continue
			}
		}

		secretData[resource.Name] = resourceSecrets
	}

	data := map[string]any{
		"Resources": appDef.Resources,
		"Data":      secretData,
		// TODO: This may change at some point, see workflow for more details.
		"BucketName": appDef.Project.Name,
		"Env":        enviro,
	}

	// Track all resources as sources for this workflow.
	var trackingOptions []scaffold.Option
	for _, resource := range appDef.Resources {
		trackingOptions = append(trackingOptions, scaffold.WithTracking(manifest.SourceResource(resource.Name)))
	}

	return input.Generator().Template(path, tpl, data, trackingOptions...)
}
