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
	"db-add-ip/action.yaml":        "db-add-ip/action.yaml",
	"setup-webkit-cli/action.yaml": "setup-webkit-cli/action.yaml",
}

// ActionTemplates copies action.yaml files from the templates folder
// so services can use re-usable workflow helpers in CI/CD.
func ActionTemplates(_ context.Context, input cmdtools.CommandInput) error {
	gen := scaffold.New(afero.NewBasePathFs(input.FS, actionsPath), input.Manifest)

	for from, to := range actionTemplates {
		err := gen.CopyFromEmbed(templates.Embed, filepath.Join(actionsPath, from), to)
		if err != nil {
			return err
		}
	}

	return nil
}
