package cicd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/templates"
)

// AppPRWorkflow bootstraps all the GitHub workflows for a
// WebKit application.
func AppPRWorkflow(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	if len(appDef.Apps) == 0 {
		return nil
	}

	for _, app := range appDef.Apps {
		tpl := templates.MustLoadTemplate(filepath.Join(workflowsPath, "pr.yaml.tmpl"))
		file := filepath.Join(workflowsPath, fmt.Sprintf("pr-%s.yaml", app.Name))

		if err := input.Generator().Template(file, tpl, &app); err != nil {
			return err
		}
	}

	return nil
}
