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
	path := filepath.Join(workflowsPath, "backup.yaml")

	// Build data map with all resource secrets
	secretData := make(map[string]string)
	for _, resource := range appDef.Resources {
		switch resource.Type {
		case appdef.ResourceTypePostgres:
			secretData[fmt.Sprintf("%s_DatabaseURL", resource.Name)] = resource.GitHubSecretName(enviro, "connection_url")
			secretData[fmt.Sprintf("%s_DatabaseID", resource.Name)] = resource.GitHubSecretName(enviro, "id")

		case appdef.ResourceTypeS3:
			secretData[fmt.Sprintf("%s_AccessKey", resource.Name)] = resource.GitHubSecretName(enviro, "access_key")
			secretData[fmt.Sprintf("%s_SecretKey", resource.Name)] = resource.GitHubSecretName(enviro, "secret_key")
			secretData[fmt.Sprintf("%s_Region", resource.Name)] = resource.GitHubSecretName(enviro, "region")
			secretData[fmt.Sprintf("%s_BucketName", resource.Name)] = resource.GitHubSecretName(enviro, "bucket_name")
		}
	}

	data := map[string]any{
		"Resources":  appDef.Resources,
		"Data":       secretData,
		"BucketName": appDef.Project.Name, // TODO: This may change at some point, see workflow for more details.
	}

	return input.Generator().Template(path, tpl, data)
}
