package docs

import (
	"bytes"
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

const (
	// outputPath is where the generated AGENTS.md file will be written
	outputPath = "AGENTS.md"
	// baseTemplateName is the name of the base template in internal/templates/
	baseTemplateName = "AGENTS.md"
	// customContentPathTmpl is the path for custom template content
	customContentPathTmpl = "docs/AGENTS.md.tmpl"
	// customContentPath is the path for custom static content
	customContentPath = "docs/AGENTS.md"
)

// Generate creates the AGENTS.md file at the project root by combining
// the base template with optional custom content from docs/.
func Generate(_ context.Context, input cmdtools.CommandInput) error {
	baseTemplate, err := templates.LoadTemplate(baseTemplateName)
	if err != nil {
		return errors.Wrap(err, "loading base template")
	}

	customContent, err := loadCustomContent(input.FS, input.AppDef())
	if err != nil {
		return errors.Wrap(err, "loading custom content")
	}

	data := map[string]any{
		"Definition": input.AppDef(),
		"Content":    customContent,
	}

	err = input.Generator().Template(
		outputPath,
		baseTemplate,
		data,
		scaffold.WithTracking(manifest.SourceProject()),
	)
	if err != nil {
		return errors.Wrap(err, "generating AGENTS.md")
	}

	input.Printer().Success("Generated AGENTS.md")

	return nil
}

// loadCustomContent attempts to load custom content from docs/ directory.
// It tries docs/AGENTS.md.tmpl first, then docs/AGENTS.md, and returns
// an empty string if neither exists.
func loadCustomContent(fs afero.Fs, appDef any) (string, error) {
	// Try loading template file first
	if exists, _ := afero.Exists(fs, customContentPathTmpl); exists {
		content, err := afero.ReadFile(fs, customContentPathTmpl)
		if err != nil {
			return "", errors.Wrap(err, "reading custom template")
		}

		tmpl, err := templates.LoadTemplate("AGENTS.md.tmpl")
		if err != nil {
			// If LoadTemplate fails, treat it as a static file
			return string(content), nil
		}

		buf := &bytes.Buffer{}
		data := map[string]any{
			"Definition": appDef,
		}
		if err := tmpl.Execute(buf, data); err != nil {
			return "", errors.Wrap(err, "executing custom template")
		}

		return buf.String(), nil
	}

	// Fallback to static markdown file
	if exists, _ := afero.Exists(fs, customContentPath); exists {
		content, err := afero.ReadFile(fs, customContentPath)
		if err != nil {
			return "", errors.Wrap(err, "reading custom content")
		}
		return string(content), nil
	}

	return "", nil
}
