package docs

import (
	"bytes"
	"context"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/docs"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

const (
	// agentsPath is the path to the static agents markdown file.
	agentsPath = "docs/AGENTS.md"

	// agentsPathTpl is the path to the agents template file.
	agentsPathTpl = "docs/AGENTS.md.tmpl"
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
	baseTemplate := templates.MustLoadTemplate("AGENTS.md")

	// Try to load manifest, but don't fail if it doesn't exist
	def, _ := appdef.Read(input.FS)

	// Load custom content (supports both static and template files)
	customContent, err := loadCustomContent(input.FS, def)
	if err != nil {
		return errors.Wrap(err, "loading custom content")
	}

	// Load generated guidelines
	codeStyle := docsutil.MustLoadGenFile(input.FS, docsutil.CodeStyleTemplate)

	data := map[string]any{
		"Definition": def,
		"Content":    customContent,
		"CodeStyle":  codeStyle,
	}

	// Conditionally add app-specific guidelines for root AGENTS.md
	if def != nil && def.HasAppType(appdef.AppTypePayload) {
		payload := docsutil.MustLoadGenFile(input.FS, docsutil.PayloadTemplate)
		data["Payload"] = payload
	}

	if def != nil && def.HasAppType(appdef.AppTypeSvelteKit) {
		svelteKit := docsutil.MustLoadGenFile(input.FS, docsutil.SvelteKitTemplate)
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
	payloadApps := def.GetAppsByType(appdef.AppTypePayload)
	for _, app := range payloadApps {
		if err := generateAppAgentsFile(input, app, docsutil.PayloadTemplate); err != nil {
			return errors.Wrap(err, "generating Payload AGENTS.md")
		}
	}

	// Generate for SvelteKit apps
	svelteKitApps := def.GetAppsByType(appdef.AppTypeSvelteKit)
	for _, app := range svelteKitApps {
		if err := generateAppAgentsFile(input, app, docsutil.SvelteKitTemplate); err != nil {
			return errors.Wrap(err, "generating SvelteKit AGENTS.md")
		}
	}

	return nil
}

// generateAppAgentsFile creates an AGENTS.md file in the app's directory.
func generateAppAgentsFile(input cmdtools.CommandInput, app appdef.App, genFile docsutil.Template) error {
	content := docsutil.MustLoadGenFile(input.FS, genFile)

	// Determine which template to use based on the app type
	var templateName string
	var dataKey string
	switch genFile {
	case docsutil.PayloadTemplate:
		templateName = "AGENTS.PAYLOAD.md"
		dataKey = "Payload"
	case docsutil.SvelteKitTemplate:
		templateName = "AGENTS.SVELTEKIT.md"
		dataKey = "SvelteKit"
	default:
		return errors.New("unknown app type for template")
	}

	// Load the app-specific template
	tmpl := templates.MustLoadTemplate(templateName)

	// Prepare data for template
	data := map[string]any{
		dataKey: content,
	}

	// Write to app directory using template
	agentsPath := filepath.Join(app.Path, "AGENTS.md")
	err := input.Generator().Template(
		agentsPath,
		tmpl,
		data,
		scaffold.WithTracking(manifest.SourceProject()),
	)

	return errors.Wrap(err, "writing app AGENTS.md")
}

// loadCustomContent loads custom content from docs/AGENTS.md or docs/AGENTS.md.tmpl.
// If a template file exists, it will be rendered with the provided definition.
// Template files take precedence over static markdown files.
func loadCustomContent(fs afero.Fs, def *appdef.Definition) (string, error) {
	// Check for template file first (takes precedence)
	tplExists, err := afero.Exists(fs, agentsPathTpl)
	if err != nil {
		return "", errors.Wrap(err, "checking for template file")
	}

	if tplExists {
		// Read and render template
		tplContent, err := afero.ReadFile(fs, agentsPathTpl)
		if err != nil {
			return "", errors.Wrap(err, "reading template file")
		}

		// Parse and execute template
		tmpl, err := template.New("agents").Funcs(sprig.FuncMap()).Parse(string(tplContent))
		if err != nil {
			return "", errors.Wrap(err, "parsing template")
		}

		var buf bytes.Buffer
		data := map[string]any{
			"Definition": def,
		}
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", errors.Wrap(err, "executing template")
		}

		return buf.String(), nil
	}

	// Check for static markdown file
	staticExists, err := afero.Exists(fs, agentsPath)
	if err != nil {
		return "", errors.Wrap(err, "checking for static file")
	}

	if staticExists {
		content, err := afero.ReadFile(fs, agentsPath)
		if err != nil {
			return "", errors.Wrap(err, "reading static file")
		}
		return string(content), nil
	}

	// No custom content found
	return "", nil
}
