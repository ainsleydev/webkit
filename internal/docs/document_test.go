package docs

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/printer"
	"github.com/ainsleydev/webkit/internal/scaffold"
)

func TestGenerateDocument(t *testing.T) {
	t.Parallel()

	t.Run("Generates with custom markdown content", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		tracker := manifest.NewTracker()
		console := printer.NewTest()
		gen := scaffold.New(fs, tracker, console)

		customContent := "# Custom Content\n\nThis is custom."
		err := afero.WriteFile(fs, "ai/docs/AGENTS.md", []byte(customContent), 0o644)
		require.NoError(t, err)

		opts := GenerateDocumentOptions{
			FS:                fs,
			Generator:         gen,
			DocumentType:      DocumentTypeAgents,
			CustomContentPath: "ai/docs",
			Data:              nil,
			TrackingSource:    manifest.SourceProject(),
		}

		err = GenerateDocument(opts)
		require.NoError(t, err)

		exists, err := afero.Exists(fs, "AGENTS.md")
		require.NoError(t, err)
		assert.True(t, exists)

		content, err := afero.ReadFile(fs, "AGENTS.md")
		require.NoError(t, err)
		assert.Contains(t, string(content), "# Agent Guidelines")
		assert.Contains(t, string(content), customContent)
	})

	t.Run("Generates with custom template content and data", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		tracker := manifest.NewTracker()
		console := printer.NewTest()
		gen := scaffold.New(fs, tracker, console)

		customTemplate := "# Custom Template\n\nData: {{ .Name }}"
		err := afero.WriteFile(fs, "ai/docs/AGENTS.md.tmpl", []byte(customTemplate), 0o644)
		require.NoError(t, err)

		opts := GenerateDocumentOptions{
			FS:                fs,
			Generator:         gen,
			DocumentType:      DocumentTypeAgents,
			CustomContentPath: "ai/docs",
			Data: map[string]any{
				"Name": "WebKit",
			},
			TrackingSource: manifest.SourceProject(),
		}

		err = GenerateDocument(opts)
		require.NoError(t, err)

		content, err := afero.ReadFile(fs, "AGENTS.md")
		require.NoError(t, err)
		assert.Contains(t, string(content), "# Agent Guidelines")
		assert.Contains(t, string(content), "Data: WebKit")
	})

	t.Run("Generates without custom content", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		tracker := manifest.NewTracker()
		console := printer.NewTest()
		gen := scaffold.New(fs, tracker, console)

		opts := GenerateDocumentOptions{
			FS:                fs,
			Generator:         gen,
			DocumentType:      DocumentTypeAgents,
			CustomContentPath: "ai/docs",
			Data:              nil,
			TrackingSource:    manifest.SourceProject(),
		}

		err := GenerateDocument(opts)
		require.NoError(t, err)

		content, err := afero.ReadFile(fs, "AGENTS.md")
		require.NoError(t, err)
		assert.Contains(t, string(content), "# Agent Guidelines")
	})

	t.Run("Template takes precedence over markdown", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		tracker := manifest.NewTracker()
		console := printer.NewTest()
		gen := scaffold.New(fs, tracker, console)

		err := afero.WriteFile(fs, "ai/docs/AGENTS.md", []byte("Markdown content"), 0o644)
		require.NoError(t, err)
		err = afero.WriteFile(fs, "ai/docs/AGENTS.md.tmpl", []byte("Template content"), 0o644)
		require.NoError(t, err)

		opts := GenerateDocumentOptions{
			FS:                fs,
			Generator:         gen,
			DocumentType:      DocumentTypeAgents,
			CustomContentPath: "ai/docs",
			Data:              nil,
			TrackingSource:    manifest.SourceProject(),
		}

		err = GenerateDocument(opts)
		require.NoError(t, err)

		content, err := afero.ReadFile(fs, "AGENTS.md")
		require.NoError(t, err)
		assert.Contains(t, string(content), "Template content")
		assert.NotContains(t, string(content), "Markdown content")
	})

	t.Run("Works with service/app repo with Definition data", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		tracker := manifest.NewTracker()
		console := printer.NewTest()
		gen := scaffold.New(fs, tracker, console)

		customTemplate := "# App: {{ .Definition.Name }}"
		err := afero.WriteFile(fs, "docs/AGENTS.md.tmpl", []byte(customTemplate), 0o644)
		require.NoError(t, err)

		type AppDef struct{ Name string }
		opts := GenerateDocumentOptions{
			FS:                fs,
			Generator:         gen,
			DocumentType:      DocumentTypeAgents,
			CustomContentPath: "docs",
			Data: map[string]any{
				"Definition": AppDef{Name: "MyApp"},
			},
			TrackingSource: manifest.SourceProject(),
		}

		err = GenerateDocument(opts)
		require.NoError(t, err)

		content, err := afero.ReadFile(fs, "AGENTS.md")
		require.NoError(t, err)
		assert.Contains(t, string(content), "# App: MyApp")
	})
}

func TestLoadCustomContent(t *testing.T) {
	t.Parallel()

	t.Run("Loads template file with data", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		customTemplate := "# Template\n\nHello {{ .Name }}"
		err := afero.WriteFile(fs, "custom/AGENTS.md.tmpl", []byte(customTemplate), 0o644)
		require.NoError(t, err)

		data := map[string]any{"Name": "Test"}
		content, err := loadCustomContent(fs, "custom", DocumentTypeAgents, data)
		require.NoError(t, err)
		assert.Contains(t, content, "# Template")
		assert.Contains(t, content, "Hello Test")
	})

	t.Run("Loads markdown file", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		customMarkdown := "# Markdown\n\nStatic content"
		err := afero.WriteFile(fs, "custom/AGENTS.md", []byte(customMarkdown), 0o644)
		require.NoError(t, err)

		content, err := loadCustomContent(fs, "custom", DocumentTypeAgents, nil)
		require.NoError(t, err)
		assert.Equal(t, customMarkdown, content)
	})

	t.Run("Returns empty string when no files exist", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		content, err := loadCustomContent(fs, "custom", DocumentTypeAgents, nil)
		require.NoError(t, err)
		assert.Empty(t, content)
	})

	t.Run("Template file error", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		invalidTemplate := "{{ .Invalid"
		err := afero.WriteFile(fs, "custom/AGENTS.md.tmpl", []byte(invalidTemplate), 0o644)
		require.NoError(t, err)

		content, err := loadCustomContent(fs, "custom", DocumentTypeAgents, nil)
		assert.Error(t, err)
		assert.Empty(t, content)
	})
}

func TestDocumentType(t *testing.T) {
	t.Parallel()

	t.Run("String returns correct value", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "AGENTS.md", DocumentTypeAgents.String())
	})

	t.Run("TemplateName returns correct value", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "AGENTS.md", DocumentTypeAgents.TemplateName())
	})

	t.Run("CustomTemplateName returns correct value", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "AGENTS.md.tmpl", DocumentTypeAgents.CustomTemplateName())
	})
}
