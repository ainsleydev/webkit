package docs

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/fsext"
	"github.com/ainsleydev/webkit/internal/gen"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

// Agents creates the AGENTS.md file at the project root by combining
// the base template with generated guidelines from internal/gen/docs/.
func Agents(_ context.Context, input cmdtools.CommandInput) error {
	// Generate root AGENTS.md.
	if err := generateRootAgents(input); err != nil {
		return errors.Wrap(err, "generating root AGENTS.md")
	}

	// Generate app-specific AGENTS.md files for Payload and SvelteKit apps.
	if err := generateAppSpecificAgents(input); err != nil {
		return errors.Wrap(err, "generating app-specific AGENTS.md")
	}

	return nil
}

func generateRootAgents(input cmdtools.CommandInput) error {
	baseTemplate := templates.MustLoadTemplate("docs/AGENTS.md")

	data := map[string]any{
		"Definition": input.AppDef(),
		"Content":    mustLoadCustomContent(input.FS, "AGENTS.md"),
		"CodeStyle":  fsext.MustReadFromEmbed(gen.Embed, "docs/CODE_STYLE.md"),
	}

	err := input.Generator().Template(
		"AGENTS.md",
		baseTemplate,
		data,
		scaffold.WithTracking(manifest.SourceProject()),
	)
	if err != nil {
		return errors.Wrap(err, "generating AGENTS.md")
	}

	return nil
}

// generateAppSpecificAgents creates AGENTS.md files in app subdirectories
// for Payload and SvelteKit apps.
func generateAppSpecificAgents(input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	// Generate for Payload apps
	payloadApps := appDef.GetAppsByType(appdef.AppTypePayload)
	for _, app := range payloadApps {
		err := input.Generator().Template(
			filepath.Join(app.Path, "AGENTS.md"),
			templates.MustLoadTemplate("docs/AGENTS.PAYLOAD.md"),
			map[string]any{
				"Payload": fsext.MustReadFromEmbed(gen.Embed, "docs/PAYLOAD.md"),
			},
			scaffold.WithTracking(manifest.SourceProject()),
		)
		if err != nil {
			return errors.Wrap(err, "generating Payload AGENTS.md")
		}
	}

	// Generate for SvelteKit apps
	svelteKitApps := appDef.GetAppsByType(appdef.AppTypeSvelteKit)
	for _, app := range svelteKitApps {
		err := input.Generator().Template(
			filepath.Join(app.Path, "AGENTS.md"),
			templates.MustLoadTemplate("docs/AGENTS.SVELTEKIT.md"),
			map[string]any{
				"SvelteKit": fsext.MustReadFromEmbed(gen.Embed, "docs/SVELTEKIT.md"),
			},
			scaffold.WithTracking(manifest.SourceProject()),
		)
		if err != nil {
			return errors.Wrap(err, "generating SvelteKit AGENTS.md")
		}
	}

	return nil
}
