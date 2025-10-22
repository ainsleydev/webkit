package cicd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
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

	if len(appDef.Resources) == 0 {
		return nil
	}

	tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, "backup.yaml.tmpl"))

	for _, resource := range appDef.Resources {
		path := filepath.Join(workflowsPath, fmt.Sprintf("backup-%s.yaml", resource.Name))

		data := map[string]any{
			"Resource": resource,
		}

		// Add resource-specific data based on type
		switch resource.Type {
		case appdef.ResourceTypePostgres:
			data["DatabaseURL"] = resource.GitHubSecretName(enviro, "connection_url")
			data["DatabaseID"] = resource.GitHubSecretName(enviro, "id")
			data["BucketName"] = appDef.Project.Name // TODO: This may change at some point, see workflow for more details.

		case appdef.ResourceTypeS3:
			data["AccessKey"] = resource.GitHubSecretName(enviro, "access_key")
			data["SecretKey"] = resource.GitHubSecretName(enviro, "secret_key")
			data["Region"] = resource.GitHubSecretName(enviro, "region")
			data["BucketName"] = resource.GitHubSecretName(enviro, "bucket_name")
		}

		if err := input.Generator().Template(path, tpl, data); err != nil {
			return err
		}
	}

	return nil
}
