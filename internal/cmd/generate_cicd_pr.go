package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

var createCICDCmd = &cli.Command{
	Name:   "cicd",
	Action: cmdtools.WrapCommand(createCICD),
}

func createCICD(_ context.Context, input cmdtools.CommandInput) error {
	gen := cgtools.NewGenerator(afero.NewBasePathFs(input.FS, "./.github"))

	// TODO: We need to apply defaults to the app spec for commands to work?
	for _, app := range input.AppDef().Apps {

		tpl := templates.MustLoadTemplate(".github/workflows/pr.yaml.tmpl")
		file := fmt.Sprintf("./workflows/%s.yaml", app.Name)

		err := gen.GenerateTemplate(file, tpl, app)
		if err != nil {
			return err
		}
	}

	return nil
}
