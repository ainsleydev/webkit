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

func TestGenerateAgents(t *testing.T) {
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

		opts := GenerateAgentsOptions{
			FS:                fs,
			Generator:         gen,
			CustomContentPath: "ai/docs",
			OutputPath:        "AGENTS.md",
			TemplateData:      nil,
			TrackingSource:    manifest.SourceProject(),
		}

		err = GenerateAgents(opts)
		require.NoError(t, err)

		exists, err := afero.Exists(fs, "AGENTS.md")
		require.NoError(t, err)
		assert.True(t, exists)

		content, err := afero.ReadFile(fs, "AGENTS.md")
		require.NoError(t, err)
		assert.Contains(t, string(content), "# Agent Guidelines")
		assert.Contains(t, string(content), customContent)
	})

	t.Run("Generates with custom template content", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		tracker := manifest.NewTracker()
		console := printer.NewTest()
		gen := scaffold.New(fs, tracker, console)

		customTemplate := "# Custom Template\n\nData: {{ .Definition.Name }}"
		err := afero.WriteFile(fs, "ai/docs/AGENTS.md.tmpl", []byte(customTemplate), 0o644)
		require.NoError(t, err)

		templateData := struct{ Name string }{Name: "WebKit"}

		opts := GenerateAgentsOptions{
			FS:                fs,
			Generator:         gen,
			CustomContentPath: "ai/docs",
			OutputPath:        "AGENTS.md",
			TemplateData:      templateData,
			TrackingSource:    manifest.SourceProject(),
		}

		err = GenerateAgents(opts)
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

		opts := GenerateAgentsOptions{
			FS:                fs,
			Generator:         gen,
			CustomContentPath: "ai/docs",
			OutputPath:        "AGENTS.md",
			TemplateData:      nil,
			TrackingSource:    manifest.SourceProject(),
		}

		err := GenerateAgents(opts)
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

		opts := GenerateAgentsOptions{
			FS:                fs,
			Generator:         gen,
			CustomContentPath: "ai/docs",
			OutputPath:        "AGENTS.md",
			TemplateData:      nil,
			TrackingSource:    manifest.SourceProject(),
		}

		err = GenerateAgents(opts)
		require.NoError(t, err)

		content, err := afero.ReadFile(fs, "AGENTS.md")
		require.NoError(t, err)
		assert.Contains(t, string(content), "Template content")
		assert.NotContains(t, string(content), "Markdown content")
	})
}

func TestLoadCustomContent(t *testing.T) {
	t.Parallel()

	t.Run("Loads template file", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		customTemplate := "# Template\n\nHello {{ .Definition.Name }}"
		err := afero.WriteFile(fs, "custom/AGENTS.md.tmpl", []byte(customTemplate), 0o644)
		require.NoError(t, err)

		templateData := struct{ Name string }{Name: "Test"}
		content, err := loadCustomContent(fs, "custom", templateData)
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

		content, err := loadCustomContent(fs, "custom", nil)
		require.NoError(t, err)
		assert.Equal(t, customMarkdown, content)
	})

	t.Run("Returns empty string when no files exist", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		content, err := loadCustomContent(fs, "custom", nil)
		require.NoError(t, err)
		assert.Empty(t, content)
	})

	t.Run("Template file error", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		invalidTemplate := "{{ .Invalid"
		err := afero.WriteFile(fs, "custom/AGENTS.md.tmpl", []byte(invalidTemplate), 0o644)
		require.NoError(t, err)

		content, err := loadCustomContent(fs, "custom", nil)
		assert.Error(t, err)
		assert.Empty(t, content)
	})
}
