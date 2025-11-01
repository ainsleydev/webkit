package files

import (
	"context"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

// Hooks scaffolds lefthook configuration for git hooks.
// This includes pre-commit and pre-push hooks for:
// - Checking SOPS secrets are encrypted
// - Running format commands
func Hooks(_ context.Context, input cmdtools.CommandInput) error {
	app := input.AppDef()

	return input.Generator().Template(
		"lefthook.yml",
		templates.MustLoadTemplate("lefthook.yml"),
		app,
		scaffold.WithTracking(manifest.SourceProject()),
	)
}
