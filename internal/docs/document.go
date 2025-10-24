package docs

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

// GenerateDocumentOptions configures document generation.
type GenerateDocumentOptions struct {
	// FS is the filesystem to read custom content from and write output to.
	FS afero.Fs
	// Generator is used to write the final document file.
	Generator scaffold.Generator
	// DocumentType specifies which document to generate (e.g., AGENTS.md).
	DocumentType DocumentType
	// CustomContentPath is the path to check for custom content (e.g., "docs", "ai/docs").
	CustomContentPath string
	// Data is optional template data. Can be nil for webkit repo, or contain
	// app.json definition for service/app repos.
	Data map[string]any
	// TrackingSource is the manifest source identifier for tracking.
	TrackingSource string
}

// GenerateDocument creates a document file by combining the base template
// with optional custom content from the specified directory.
//
// For service/app repos with app.json:
//
//	GenerateDocument(GenerateDocumentOptions{
//	    DocumentType: DocumentTypeAgents,
//	    CustomContentPath: "docs",
//	    Data: map[string]any{"Definition": appDef},
//	    ...
//	})
//
// For webkit repo without app.json:
//
//	GenerateDocument(GenerateDocumentOptions{
//	    DocumentType: DocumentTypeAgents,
//	    CustomContentPath: "ai/docs",
//	    Data: nil,  // No manifest needed
//	    ...
//	})
func GenerateDocument(opts GenerateDocumentOptions) error {
	baseTemplate := templates.MustLoadTemplate(opts.DocumentType.TemplateName())

	customContent, err := loadCustomContent(
		opts.FS,
		opts.CustomContentPath,
		opts.DocumentType,
		opts.Data,
	)
	if err != nil {
		return errors.Wrap(err, "loading custom content")
	}

	// Merge user data with Content field.
	templateData := make(map[string]any)
	if opts.Data != nil {
		for k, v := range opts.Data {
			templateData[k] = v
		}
	}
	templateData["Content"] = customContent

	err = opts.Generator.Template(
		opts.DocumentType.String(),
		baseTemplate,
		templateData,
		scaffold.WithTracking(opts.TrackingSource),
	)
	if err != nil {
		return errors.Wrapf(err, "generating %s", opts.DocumentType)
	}

	return nil
}

// loadCustomContent attempts to load custom content from the specified directory.
// It tries {path}/{document}.tmpl first, then {path}/{document}, and returns
// an empty string if neither exists.
func loadCustomContent(
	fs afero.Fs,
	customPath string,
	docType DocumentType,
	templateData map[string]any,
) (string, error) {
	templatePath := fmt.Sprintf("%s/%s", customPath, docType.CustomTemplateName())
	markdownPath := fmt.Sprintf("%s/%s", customPath, docType.String())

	// Try loading the template file first.
	if exists, _ := afero.Exists(fs, templatePath); exists {
		tmpl, err := templates.LoadTemplateFromFS(fs, templatePath)
		if err != nil {
			return "", errors.Wrap(err, "loading custom template")
		}

		buf := &bytes.Buffer{}
		if err = tmpl.Execute(buf, templateData); err != nil {
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
