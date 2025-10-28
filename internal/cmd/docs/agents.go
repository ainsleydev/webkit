package docs

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	docsutil "github.com/ainsleydev/webkit/internal/docs"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

// Agents creates the AGENTS.md file at the project root by combining
// the base template with generated guidelines from internal/gen/docs/.
func Agents(_ context.Context, input cmdtools.CommandInput) error {
	// Generate root AGENTS.md
	if err := generateRootAgents(input); err != nil {
		return errors.Wrap(err, "generating root AGENTS.md")
	}

	// Generate app-specific AGENTS.md files for Payload and SvelteKit apps
	if err := generateAppSpecificAgents(input); err != nil {
		return errors.Wrap(err, "generating app-specific AGENTS.md")
	}

	return nil
}

// generateRootAgents creates the root AGENTS.md file.
func generateRootAgents(input cmdtools.CommandInput) error {
	baseTemplate := templates.MustLoadTemplate("AGENTS.md")

	// Load custom content
	customContent, err := docsutil.LoadCustomContent(input.FS)
	if err != nil {
		return errors.Wrap(err, "loading custom content")
	}

	// Load generated guidelines
	codeStyle, _ := docsutil.LoadGenFile(input.FS, "CODE_STYLE.md")
	git, _ := docsutil.LoadGenFile(input.FS, "GIT.md")

	// Try to load manifest, but don't fail if it doesn't exist
	def, _ := appdef.Read(input.FS)

	data := map[string]any{
		"Definition": def,
		"Content":    customContent,
		"CodeStyle":  codeStyle,
		"Git":        git,
	}

	// Conditionally add app-specific guidelines for root AGENTS.md
	if docsutil.HasAppType(def, appdef.AppTypePayload) {
		payload, _ := docsutil.LoadGenFile(input.FS, "PAYLOAD.md")
		data["Payload"] = payload
	}

	if docsutil.HasAppType(def, appdef.AppTypeSvelteKit) {
		svelteKit, _ := docsutil.LoadGenFile(input.FS, "SVELTEKIT.md")
		data["SvelteKit"] = svelteKit
	}

	err = input.Generator().Template(
		"AGENTS.md",
		baseTemplate,
		data,
		scaffold.WithTracking(manifest.SourceProject()),
	)

	return errors.Wrap(err, "generating AGENTS.md")
}

// generateAppSpecificAgents creates AGENTS.md files in app subdirectories
// for Payload and SvelteKit apps.
func generateAppSpecificAgents(input cmdtools.CommandInput) error {
	// Try to load manifest, but don't fail if it doesn't exist
	def, _ := appdef.Read(input.FS)
	if def == nil {
		return nil
	}

	// Generate for Payload apps
	payloadApps := docsutil.GetAppsByType(def, appdef.AppTypePayload)
	for _, app := range payloadApps {
		if err := generateAppAgentsFile(input, app, "PAYLOAD.md"); err != nil {
			return errors.Wrap(err, "generating Payload AGENTS.md")
		}
	}

	// Generate for SvelteKit apps
	svelteKitApps := docsutil.GetAppsByType(def, appdef.AppTypeSvelteKit)
	for _, app := range svelteKitApps {
		if err := generateAppAgentsFile(input, app, "SVELTEKIT.md"); err != nil {
			return errors.Wrap(err, "generating SvelteKit AGENTS.md")
		}
	}

	return nil
}

// generateAppAgentsFile creates an AGENTS.md file in the app's directory.
func generateAppAgentsFile(input cmdtools.CommandInput, app appdef.App, genFile string) error {
	content, err := docsutil.LoadGenFile(input.FS, genFile)
	if err != nil {
		return errors.Wrap(err, "loading generated file")
	}

	// If no content, skip
	if content == "" {
		return nil
	}

	// Write to app directory
	agentsPath := filepath.Join(app.Path, "AGENTS.md")
	err = input.Generator().Bytes(
		agentsPath,
		[]byte(content),
		scaffold.WithTracking(manifest.SourceProject()),
	)

	return errors.Wrap(err, "writing app AGENTS.md")
}
