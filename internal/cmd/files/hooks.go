package files

import (
	"context"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/state/manifest"
	"github.com/ainsleydev/webkit/internal/templates"
)

// Hooks scaffolds lefthook configuration for git hooks.
func Hooks(_ context.Context, input cmdtools.CommandInput) error {
	app := input.AppDef()

	return input.Generator().Template(
		"lefthook.yaml",
		templates.MustLoadTemplate("lefthook.yaml"),
		app,
		scaffold.WithTracking(manifest.SourceProject()),
	)
}
