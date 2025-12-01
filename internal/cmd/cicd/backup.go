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

var BackupCmd = &cli.Command{
	Name:   "backup",
	Usage:  "Generate backup workflows for resources",
	Action: cmdtools.Wrap(BackupWorkflow),
}

// filterBackupableResources returns resources that support automated backups
// based on the backup.enabled flag and provider compatibility.
func filterBackupableResources(resources []appdef.Resource) []appdef.Resource {
	var result []appdef.Resource
	for _, resource := range resources {
		if !resource.IsBackupEnabled() {
			continue
		}

		// Check provider compatibility.
		switch resource.Type {
		case appdef.ResourceTypeS3:
			// S3 backup only compatible with DigitalOcean Spaces.
			if resource.Provider != appdef.ResourceProviderDigitalOcean {
				continue
			}
		case appdef.ResourceTypeSQLite:
			// SQLite backup only compatible with Turso.
			if resource.Provider != appdef.ResourceProviderTurso {
				continue
			}
		case appdef.ResourceTypePostgres:
			// Postgres supports all providers.
		default:
			// Skip unknown resource types.
			continue
		}

		result = append(result, resource)
	}
	return result
}

// BackupWorkflow creates backup workflows for every resource if the
// backup config is enabled.
func BackupWorkflow(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	enviro := env.Production

	tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, "backup.yaml.tmpl"))
	path := filepath.Join(workflowsPath, "backup.yaml")

	// Filter resources that support backup.
	backupableResources := filterBackupableResources(appDef.Resources)

	// Build nested data map with all resource secrets grouped by resource name.
	// This allows cleaner template access: {{ index $.Data .Name "DatabaseURL" }}
	secretData := make(map[string]map[string]string)
	for _, resource := range backupableResources {
		resourceSecrets := make(map[string]string)

		switch resource.Type {
		case appdef.ResourceTypePostgres:
			resourceSecrets["DatabaseURL"] = resource.GitHubSecretName(enviro, "connection_url")
			resourceSecrets["DatabaseID"] = resource.GitHubSecretName(enviro, "id")

		case appdef.ResourceTypeS3:
			resourceSecrets["AccessKey"] = resource.GitHubSecretName(enviro, "access_key")
			resourceSecrets["SecretKey"] = resource.GitHubSecretName(enviro, "secret_key")
			resourceSecrets["Region"] = resource.GitHubSecretName(enviro, "region")
			resourceSecrets["BucketName"] = resource.GitHubSecretName(enviro, "bucket_name")

		case appdef.ResourceTypeSQLite:
			// SQLite backup uses TURSO_API_TOKEN environment variable.
		}

		secretData[resource.Name] = resourceSecrets
	}

	data := map[string]any{
		"Resources":         backupableResources,
		"Data":              secretData,
		"MonitoringEnabled": appDef.Monitoring.IsEnabled(),
		// TODO: This may change at some point, see workflow for more details.
		"BucketName": appDef.Project.Name,
		"Env":        enviro,
	}

	// Track backupable resources as sources for this workflow.
	var trackingOptions []scaffold.Option
	for _, resource := range backupableResources {
		trackingOptions = append(trackingOptions, scaffold.WithTracking(manifest.SourceResource(resource.Name)))
	}

	return input.Generator().Template(path, tpl, data, trackingOptions...)
}
