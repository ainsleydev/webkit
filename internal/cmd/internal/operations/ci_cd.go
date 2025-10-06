package operations

import (
	"context"
	"fmt"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

// CreateCICD bootstraps all of the GitHub workflows for a
// WebKit application.
func CreateCICD(_ context.Context, input cmdtools.CommandInput) error {
	gen := scaffold.New(afero.NewBasePathFs(input.FS, "./.github"))
	app := input.AppDef()

	for _, app := range app.Apps {
		tpl := templates.MustLoadTemplate(".github/workflows/pr.yaml.tmpl")
		file := fmt.Sprintf("./workflows/app-pr-%s.yaml", app.Name)

		if err := gen.Template(file, tpl, &app); err != nil {
			return err
		}
	}

	for _, resource := range app.Resources {
		backupEnabled := resource.Backup.Enabled

		if resource.Type == appdef.ResourceTypePostgres && backupEnabled {
			tpl := templates.MustLoadTemplate(".github/workflows/backup-postgres.yaml.tmpl")
			file := fmt.Sprintf("./workflows/resource-backup-%s.yaml", resource.Name)

			if err := gen.Template(file, tpl, &resource); err != nil {
				return err
			}
		}
	}

	return nil
}
