package cicd

import (
	"context"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

var ActionsCmd = &cli.Command{
	Name:   "actions",
	Usage:  "Copies action files templates",
	Action: cmdtools.Wrap(ActionTemplates),
}

var actionTemplates = map[string]string{
	"db-add-ip/action.yaml":    "db-add-ip/action.yaml",
	"db-remove-ip/action.yaml": "db-remove-ip/action.yaml",
	"pr-comment/action.yaml":   "pr-comment/action.yaml",
	"setup-infra/action.yaml":  "setup-infra/action.yaml",
	"slack-notify/action.yaml": "slack-notify/action.yaml",
}

// ActionTemplates copies action.yaml files from the templates folder
// so services can use re-usable workflow helpers in CI/CD.
func ActionTemplates(_ context.Context, input cmdtools.CommandInput) error {
	for from, to := range actionTemplates {
		err := input.Generator().CopyFromEmbed(templates.Embed,
			filepath.Join(actionsPath, from),
			filepath.Join(actionsPath, to),
			scaffold.WithTracking(manifest.SourceProject()),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
