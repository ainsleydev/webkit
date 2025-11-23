package files

import (
	"context"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/state/manifest"
	"github.com/ainsleydev/webkit/internal/templates"
)

var claudeSettingsTemplates = map[string]string{
	".claude/settings.json": ".claude/settings.json",
}

// ClaudeSettings scaffolds the Claude Code settings files.
func ClaudeSettings(_ context.Context, input cmdtools.CommandInput) error {
	app := input.AppDef()

	for file, template := range claudeSettingsTemplates {
		err := input.Generator().Template(file,
			templates.MustLoadTemplate(template),
			app,
			scaffold.WithTracking(manifest.SourceProject()),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
