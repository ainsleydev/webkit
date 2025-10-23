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
	// agentsPathTpl is the path for custom template content
	agentsPathTpl = "docs/AGENTS.md.tmpl"

	// agentsPath is the path for custom static content
	agentsPath = "docs/AGENTS.md"
)

// Agents creates the AGENTS.md file at the project root by combining
// the base template with optional custom content from docs/.
func Agents(_ context.Context, input cmdtools.CommandInput) error {
	baseTemplate := templates.MustLoadTemplate("AGENTS.md")

	customContent, err := loadCustomContent(input.FS, input.AppDef())
	if err != nil {
		return errors.Wrap(err, "loading custom content")
	}

	data := map[string]any{
		"Definition": input.AppDef(),
		"Content":    customContent,
	}

	err = input.Generator().Template(
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

// loadCustomContent attempts to load custom content from docs/ directory.
// It tries docs/AGENTS.md.tmpl first, then docs/AGENTS.md, and returns
// an empty string if neither exists.
func loadCustomContent(fs afero.Fs, appDef any) (string, error) {
	// Try loading the template file first,
	if exists, _ := afero.Exists(fs, agentsPathTpl); exists {
		tmpl, err := templates.LoadTemplateFromFS(fs, agentsPathTpl)
		if err != nil {
			return "", errors.Wrap(err, "loading custom agents content")
		}

		buf := &bytes.Buffer{}
		data := map[string]any{
			"Definition": appDef,
		}
		if err = tmpl.Execute(buf, data); err != nil {
			return "", errors.Wrap(err, "executing custom template")
		}

		return buf.String(), nil
	}

	// Fallback to a static markdown file,
	if exists, _ := afero.Exists(fs, agentsPath); exists {
		content, err := afero.ReadFile(fs, agentsPath)
		if err != nil {
			return "", errors.Wrap(err, "reading custom content")
		}
		return string(content), nil
	}

	return "", nil
}
