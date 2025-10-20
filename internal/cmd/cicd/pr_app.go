package cicd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

// AppPRWorkflow bootstraps all the GitHub workflows for a
// WebKit application.
func AppPRWorkflow(_ context.Context, input cmdtools.CommandInput) error {
	gen := scaffold.New(afero.NewBasePathFs(input.FS, "./.github"), input.Manifest)
	appDef := input.AppDef()

	if len(appDef.Apps) == 0 {
		return nil
	}

	for _, app := range appDef.Apps {
		tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, "pr.yaml.tmpl"))
		file := fmt.Sprintf("./workflows/pr-%s.yaml", app.Name)

		if err := gen.Template(file, tpl, &app); err != nil {
			return err
		}
	}

	return nil
}
