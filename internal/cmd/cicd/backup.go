package cicd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
	"github.com/ainsleydev/webkit/pkg/env"

	"github.com/spf13/afero"
)

var BackupCmd = &cli.Command{
	Name:   "backup",
	Usage:  "Generate backup workflows for resources",
	Action: cmdtools.Wrap(BackupWorkflow),
}

// BackupWorkflow creates backup workflows for every resource if the
// backup config is enabled.
func BackupWorkflow(_ context.Context, input cmdtools.CommandInput) error {
	gen := scaffold.New(afero.NewBasePathFs(input.FS, workflowsPath), input.Manifest)
	appDef := input.AppDef()
	enviro := env.Production

	if len(appDef.Resources) == 0 {
		return nil
	}

	for _, resource := range appDef.Resources {

		// Postgres DB
		if resource.Type == appdef.ResourceTypePostgres {
			tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, "backup-postgres.yaml.tmpl"))
			path := fmt.Sprintf("backup-postgres-%s.yaml", resource.Name)

			if err := gen.Template(path, tpl, map[string]any{
				"Resource":    resource,
				"DatabaseURL": resource.GitHubSecretName(enviro, "connection_url"),
				"DatabaseID":  resource.GitHubSecretName(enviro, "id"),
				"AccessKey":   resource.GitHubSecretName(enviro, "access_key"),
				"SecretKey":   resource.GitHubSecretName(enviro, "secret_key"),
			}); err != nil {
				return err
			}
		}

		// S3 Storage
		if resource.Type == appdef.ResourceTypeS3 {
			// TODO
		}
	}

	return nil
}
