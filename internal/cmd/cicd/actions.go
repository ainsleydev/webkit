package cicd

import (
	"context"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

var ActionsCmd = &cli.Command{
	Name:   "actions",
	Usage:  "Copies action files templates",
	Action: cmdtools.Wrap(ActionTemplates),
}

var actionTemplates = map[string]string{
	"db-add-ip/action.yaml": "db-add-ip.yaml",
}

func ActionTemplates(_ context.Context, input cmdtools.CommandInput) error {
	base := filepath.Join(".github", "actions")
	gen := scaffold.New(afero.NewBasePathFs(input.FS, base))
	app := input.AppDef()

	for file, template := range actionTemplates {
		tpl := templates.MustLoadTemplate(filepath.Join(".github", "actions", template))
		err := gen.Template(file, tpl, app)
		if err != nil {
			return err
		}
	}

	return nil
}
