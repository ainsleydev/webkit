package docs

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

// GenerateAgentsOptions configures the AGENTS.md generation.
type GenerateAgentsOptions struct {
	// FS is the filesystem to read custom content from and write output to.
	FS afero.Fs
	// Generator is used to write the final AGENTS.md file.
	Generator scaffold.Generator
	// CustomContentPath is the path to check for custom content (e.g., "docs", "ai/docs").
	CustomContentPath string
	// OutputPath is where the generated AGENTS.md will be written.
	OutputPath string
	// TemplateData is any additional data to pass to the template.
	TemplateData any
	// TrackingSource is the manifest source identifier for tracking.
	TrackingSource string
}

const (
	// agentsTemplateName is the filename for custom template content.
	agentsTemplateName = "AGENTS.md.tmpl"
	// agentsMarkdownName is the filename for custom static content.
	agentsMarkdownName = "AGENTS.md"
)

// GenerateAgents creates the AGENTS.md file by combining the base template
// with optional custom content from the specified directory.
func GenerateAgents(opts GenerateAgentsOptions) error {
	baseTemplate := templates.MustLoadTemplate("AGENTS.md")

	customContent, err := loadCustomContent(opts.FS, opts.CustomContentPath, opts.TemplateData)
	if err != nil {
		return errors.Wrap(err, "loading custom content")
	}

	data := map[string]any{
		"Definition": opts.TemplateData,
		"Content":    customContent,
	}

	err = opts.Generator.Template(
		opts.OutputPath,
		baseTemplate,
		data,
		scaffold.WithTracking(opts.TrackingSource),
	)
	if err != nil {
		return errors.Wrap(err, "generating AGENTS.md")
	}

	return nil
}

// loadCustomContent attempts to load custom content from the specified directory.
// It tries {path}/AGENTS.md.tmpl first, then {path}/AGENTS.md, and returns
// an empty string if neither exists.
func loadCustomContent(fs afero.Fs, customPath string, templateData any) (string, error) {
	templatePath := customPath + "/" + agentsTemplateName
	markdownPath := customPath + "/" + agentsMarkdownName

	// Try loading the template file first.
	if exists, _ := afero.Exists(fs, templatePath); exists {
		tmpl, err := templates.LoadTemplateFromFS(fs, templatePath)
		if err != nil {
			return "", errors.Wrap(err, "loading custom agents template")
		}

		buf := &bytes.Buffer{}
		data := map[string]any{
			"Definition": templateData,
		}
		if err = tmpl.Execute(buf, data); err != nil {
			return "", errors.Wrap(err, "executing custom template")
		}

		return buf.String(), nil
	}

	// Fallback to a static markdown file.
	if exists, _ := afero.Exists(fs, markdownPath); exists {
		content, err := afero.ReadFile(fs, markdownPath)
		if err != nil {
			return "", errors.Wrap(err, "reading custom content")
		}
		return string(content), nil
	}

	return "", nil
}
